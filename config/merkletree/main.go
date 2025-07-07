package merkletree

import (
    "bytes"
    "encoding/binary"
    "encoding/hex"
    "encoding/json"
    "errors"
    "fmt"
    "os"
    "path/filepath"
    "sync"
    "sync/atomic"
    "time"

    "go.etcd.io/bbolt"
    "golang.org/x/crypto/sha3"
)

// Constants
const (
    HASH_LENGTH              = 32
    METADATA_BUCKET          = "metadata"
    NODES_BUCKET            = "nodes"
    KEY_DATA_BUCKET         = "keydata"
    KEY_ROOT_HASH           = "rootHash"
    KEY_NUM_LEAVES          = "numLeaves"
    KEY_DEPTH               = "depth"
    KEY_HANGING_NODE_PREFIX = "hangingNode"
)

// Global state
var (
    openTrees  = make(map[string]*MerkleTree)
    treesMutex sync.RWMutex
)

// Node represents a single node in the Merkle Tree
type Node struct {
    Hash   []byte `json:"hash"`
    Left   []byte `json:"left,omitempty"`
    Right  []byte `json:"right,omitempty"`
    Parent []byte `json:"parent,omitempty"`
    
    // The old hash of the node before it was updated
    NodeHashToRemoveFromDB []byte `json:"nodeHashToRemoveFromDB,omitempty"`
    
    mu sync.RWMutex `json:"-"`
}

// NewLeafNode creates a new leaf node with a known hash
func NewLeafNode(hash []byte) (*Node, error) {
    if hash == nil {
        return nil, errors.New("node hash cannot be nil")
    }
    
    node := &Node{
        Hash: make([]byte, len(hash)),
    }
    copy(node.Hash, hash)
    
    return node, nil
}

// NewNodeWithFields creates a node with all fields
func NewNodeWithFields(hash, left, right, parent []byte) (*Node, error) {
    if hash == nil {
        return nil, errors.New("node hash cannot be nil")
    }
    
    node := &Node{
        Hash:   make([]byte, len(hash)),
        Left:   copyBytes(left),
        Right:  copyBytes(right),
        Parent: copyBytes(parent),
    }
    copy(node.Hash, hash)
    
    return node, nil
}

// NewInternalNode creates a node with left and right hashes, auto-calculating node hash
func NewInternalNode(left, right []byte) (*Node, error) {
    if left == nil && right == nil {
        return nil, errors.New("at least one of left or right hash must be non-nil")
    }
    
    node := &Node{
        Left:  copyBytes(left),
        Right: copyBytes(right),
    }
    
    hash := node.calculateHash()
    if hash == nil {
        return nil, errors.New("failed to calculate node hash")
    }
    
    node.Hash = hash
    return node, nil
}

// Helper function to copy byte slices
func copyBytes(src []byte) []byte {
    if src == nil {
        return nil
    }
    dst := make([]byte, len(src))
    copy(dst, src)
    return dst
}

// calculateHash calculates the hash of this node based on left and right child hashes
func (n *Node) calculateHash() []byte {
    n.mu.RLock()
    defer n.mu.RUnlock()
    
    if n.Left == nil && n.Right == nil {
        return nil
    }
    
    var leftHash, rightHash []byte
    if n.Left != nil {
        leftHash = n.Left
    } else {
        leftHash = n.Right
    }
    
    if n.Right != nil {
        rightHash = n.Right
    } else {
        rightHash = n.Left
    }
    
    return pwrHash256(leftHash, rightHash)
}

// encode serializes the node to JSON
func (n *Node) encode() ([]byte, error) {
    n.mu.RLock()
    defer n.mu.RUnlock()
    
    return json.Marshal(n)
}

// setParentNodeHash sets this node's parent
func (n *Node) setParentNodeHash(parentHash []byte) {
    n.mu.Lock()
    defer n.mu.Unlock()
    n.Parent = copyBytes(parentHash)
}

// updateLeaf updates a leaf (left or right) if it matches the old hash
func (n *Node) updateLeaf(oldLeafHash, newLeafHash []byte) error {
    n.mu.Lock()
    defer n.mu.Unlock()
    
    if n.Left != nil && bytes.Equal(n.Left, oldLeafHash) {
        n.Left = copyBytes(newLeafHash)
        return nil
    } else if n.Right != nil && bytes.Equal(n.Right, oldLeafHash) {
        n.Right = copyBytes(newLeafHash)
        return nil
    }
    
    return errors.New("old hash not found among this node's children")
}

// Hash functions - using Keccak256 to match Java implementation
func pwrHash256(data ...[]byte) []byte {
    hasher := sha3.NewLegacyKeccak256()
    for _, d := range data {
        hasher.Write(d)
    }
    return hasher.Sum(nil)
}

func CalculateLeafHash(key, data []byte) []byte {
    return pwrHash256(key, data)
}

// MerkleTree represents a Merkle Tree backed by BoltDB storage
type MerkleTree struct {
    treeName string
    path     string
    
    db *bbolt.DB
    
    // Caches
    nodesCache   map[string]*Node // Using hex string as key
    hangingNodes map[int][]byte
    keyDataCache map[string][]byte // Using hex string as key
    
    numLeaves int32
    depth     int32
    rootHash  []byte
    
    closed            int32 // atomic
    hasUnsavedChanges int32 // atomic
    
    mu sync.RWMutex
}

// NewMerkleTree creates a new MerkleTree instance
func NewMerkleTree(treeName string) (*MerkleTree, error) {
    treesMutex.Lock()
    defer treesMutex.Unlock()
    
    if _, exists := openTrees[treeName]; exists {
        return nil, errors.New("there is already open instance of this tree")
    }
    
    tree := &MerkleTree{
        treeName:     treeName,
        path:         filepath.Join("merkleTree", treeName+".db"),
        nodesCache:   make(map[string]*Node),
        hangingNodes: make(map[int][]byte),
        keyDataCache: make(map[string][]byte),
    }
    
    // Ensure directory exists
    dir := filepath.Dir(tree.path)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return nil, fmt.Errorf("failed to create directory: %s, error: %v", dir, err)
    }
    
    // Initialize database
    if err := tree.initializeDB(); err != nil {
        return nil, err
    }
    
    // Load metadata
    if err := tree.loadMetaData(); err != nil {
        tree.Close()
        return nil, err
    }
    
    // Register instance
    openTrees[treeName] = tree
    
    return tree, nil
}

func (mt *MerkleTree) initializeDB() error {
    var err error
    mt.db, err = bbolt.Open(mt.path, 0600, &bbolt.Options{
        Timeout: 1 * time.Second,
    })
    if err != nil {
        return fmt.Errorf("failed to open database: %v", err)
    }
    
    // Create buckets if they don't exist
    return mt.db.Update(func(tx *bbolt.Tx) error {
        if _, err := tx.CreateBucketIfNotExists([]byte(METADATA_BUCKET)); err != nil {
            return err
        }
        if _, err := tx.CreateBucketIfNotExists([]byte(NODES_BUCKET)); err != nil {
            return err
        }
        if _, err := tx.CreateBucketIfNotExists([]byte(KEY_DATA_BUCKET)); err != nil {
            return err
        }
        return nil
    })
}

func (mt *MerkleTree) loadMetaData() error {
    return mt.db.View(func(tx *bbolt.Tx) error {
        b := tx.Bucket([]byte(METADATA_BUCKET))
        if b == nil {
            return nil
        }
        
        // Load root hash
        if v := b.Get([]byte(KEY_ROOT_HASH)); v != nil {
            mt.rootHash = make([]byte, len(v))
            copy(mt.rootHash, v)
        }
        
        // Load number of leaves
        if v := b.Get([]byte(KEY_NUM_LEAVES)); len(v) >= 4 {
            atomic.StoreInt32(&mt.numLeaves, int32(binary.BigEndian.Uint32(v)))
        }
        
        // Load depth
        if v := b.Get([]byte(KEY_DEPTH)); len(v) >= 4 {
            atomic.StoreInt32(&mt.depth, int32(binary.BigEndian.Uint32(v)))
        }
        
        // Load hanging nodes
        depth := atomic.LoadInt32(&mt.depth)
        for i := int32(0); i <= depth; i++ {
            key := fmt.Sprintf("%s%d", KEY_HANGING_NODE_PREFIX, i)
            if v := b.Get([]byte(key)); v != nil {
                node, err := mt.getNodeByHash(v)
                if err == nil && node != nil {
                    mt.hangingNodes[int(i)] = copyBytes(node.Hash)
                }
            }
        }
        
        return nil
    })
}

func (mt *MerkleTree) errorIfClosed() error {
    if atomic.LoadInt32(&mt.closed) != 0 {
        return errors.New("merkle tree is closed")
    }
    return nil
}

// GetRootHash returns the current root hash of the Merkle tree
func (mt *MerkleTree) GetRootHash() ([]byte, error) {
    if err := mt.errorIfClosed(); err != nil {
        return nil, err
    }
    
    mt.mu.RLock()
    defer mt.mu.RUnlock()
    
    if mt.rootHash == nil {
        return nil, nil
    }
    
    result := make([]byte, len(mt.rootHash))
    copy(result, mt.rootHash)
    return result, nil
}

// GetNumLeaves returns the number of leaves in the tree
func (mt *MerkleTree) GetNumLeaves() int {
    return int(atomic.LoadInt32(&mt.numLeaves))
}

// GetDepth returns the depth of the tree
func (mt *MerkleTree) GetDepth() int {
    return int(atomic.LoadInt32(&mt.depth))
}

// GetData retrieves data for a key from the Merkle Tree
func (mt *MerkleTree) GetData(key []byte) ([]byte, error) {
    if err := mt.errorIfClosed(); err != nil {
        return nil, err
    }
    
    keyHex := hex.EncodeToString(key)
    
    mt.mu.RLock()
    if data, exists := mt.keyDataCache[keyHex]; exists {
        mt.mu.RUnlock()
        result := make([]byte, len(data))
        copy(result, data)
        return result, nil
    }
    mt.mu.RUnlock()
    
    var result []byte
    err := mt.db.View(func(tx *bbolt.Tx) error {
        b := tx.Bucket([]byte(KEY_DATA_BUCKET))
        if b == nil {
            return nil
        }
        
        if v := b.Get(key); v != nil {
            result = make([]byte, len(v))
            copy(result, v)
        }
        
        return nil
    })
    
    return result, err
}

// AddOrUpdateData adds or updates data for a key in the Merkle Tree
func (mt *MerkleTree) AddOrUpdateData(key, data []byte) error {
    if err := mt.errorIfClosed(); err != nil {
        return err
    }
    
    if key == nil {
        return errors.New("key cannot be nil")
    }
    if data == nil {
        return errors.New("data cannot be nil")
    }
    
    mt.mu.Lock()
    defer mt.mu.Unlock()
    
    // Check if key already exists
    existingData, err := mt.getData(key)
    if err != nil {
        return err
    }
    
    var oldLeafHash []byte
    if existingData != nil {
        oldLeafHash = CalculateLeafHash(key, existingData)
    }
    
    // Calculate hash from key and data
    newLeafHash := CalculateLeafHash(key, data)
    
    if oldLeafHash != nil && bytes.Equal(oldLeafHash, newLeafHash) {
        return nil // No change needed
    }
    
    // Store key-data mapping
    keyHex := hex.EncodeToString(key)
    dataCopy := make([]byte, len(data))
    copy(dataCopy, data)
    mt.keyDataCache[keyHex] = dataCopy
    atomic.StoreInt32(&mt.hasUnsavedChanges, 1)
    
    if oldLeafHash == nil {
        // Key doesn't exist, add new leaf
        leafNode, err := NewLeafNode(newLeafHash)
        if err != nil {
            return err
        }
        return mt.addLeaf(leafNode)
    } else {
        // Key exists, update leaf
        return mt.updateLeaf(oldLeafHash, newLeafHash)
    }
}

// Internal getData method that doesn't acquire locks
func (mt *MerkleTree) getData(key []byte) ([]byte, error) {
    keyHex := hex.EncodeToString(key)
    
    if data, exists := mt.keyDataCache[keyHex]; exists {
        result := make([]byte, len(data))
        copy(result, data)
        return result, nil
    }
    
    var result []byte
    err := mt.db.View(func(tx *bbolt.Tx) error {
        b := tx.Bucket([]byte(KEY_DATA_BUCKET))
        if b == nil {
            return nil
        }
        
        if v := b.Get(key); v != nil {
            result = make([]byte, len(v))
            copy(result, v)
        }
        
        return nil
    })
    
    return result, err
}

func (mt *MerkleTree) getNodeByHash(hash []byte) (*Node, error) {
    if hash == nil {
        return nil, nil
    }
    
    hashHex := hex.EncodeToString(hash)
    
    // Check cache first
    if node, exists := mt.nodesCache[hashHex]; exists {
        return node, nil
    }
    
    // Load from database
    var node *Node
    err := mt.db.View(func(tx *bbolt.Tx) error {
        b := tx.Bucket([]byte(NODES_BUCKET))
        if b == nil {
            return nil
        }
        
        if v := b.Get(hash); v != nil {
            var n Node
            if err := json.Unmarshal(v, &n); err != nil {
                return err
            }
            node = &n
        }
        
        return nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if node != nil {
        // Cache the node
        mt.nodesCache[hashHex] = node
    }
    
    return node, nil
}

func (mt *MerkleTree) addLeaf(leafNode *Node) error {
    if leafNode == nil {
        return errors.New("leaf node cannot be nil")
    }
    if leafNode.Hash == nil {
        return errors.New("leaf node hash cannot be nil")
    }
    
    // Cache the leaf node
    hashHex := hex.EncodeToString(leafNode.Hash)
    mt.nodesCache[hashHex] = leafNode
    
    numLeaves := atomic.LoadInt32(&mt.numLeaves)
    
    if numLeaves == 0 {
        mt.hangingNodes[0] = copyBytes(leafNode.Hash)
        mt.rootHash = copyBytes(leafNode.Hash)
    } else {
        hangingLeafHash := mt.hangingNodes[0]
        hangingLeaf, err := mt.getNodeByHash(hangingLeafHash)
        if err != nil {
            return err
        }
        
        if hangingLeaf == nil {
            mt.hangingNodes[0] = copyBytes(leafNode.Hash)
            parentNode, err := NewInternalNode(leafNode.Hash, nil)
            if err != nil {
                return err
            }
            leafNode.setParentNodeHash(parentNode.Hash)
            return mt.addNode(1, parentNode)
        } else {
            if hangingLeaf.Parent == nil { // hanging leaf is the root
                parentNode, err := NewInternalNode(hangingLeaf.Hash, leafNode.Hash)
                if err != nil {
                    return err
                }
                hangingLeaf.setParentNodeHash(parentNode.Hash)
                leafNode.setParentNodeHash(parentNode.Hash)
                delete(mt.hangingNodes, 0)
                return mt.addNode(1, parentNode)
            } else {
                parentNodeOfHangingLeaf, err := mt.getNodeByHash(hangingLeaf.Parent)
                if err != nil {
                    return err
                }
                if parentNodeOfHangingLeaf == nil {
                    return errors.New("parent node of hanging leaf not found")
                }
                if err := mt.addLeafToNode(parentNodeOfHangingLeaf, leafNode.Hash); err != nil {
                    return err
                }
                delete(mt.hangingNodes, 0)
            }
        }
    }
    
    atomic.AddInt32(&mt.numLeaves, 1)
    return nil
}

func (mt *MerkleTree) addNode(level int, node *Node) error {
    currentDepth := atomic.LoadInt32(&mt.depth)
    if int32(level) > currentDepth {
        atomic.StoreInt32(&mt.depth, int32(level))
    }
    
    // Cache the node
    hashHex := hex.EncodeToString(node.Hash)
    mt.nodesCache[hashHex] = node
    
    hangingNode, err := mt.getNodeByHash(mt.hangingNodes[level])
    if err != nil {
        return err
    }
    
    if hangingNode == nil {
        mt.hangingNodes[level] = copyBytes(node.Hash)
        
        if int32(level) >= atomic.LoadInt32(&mt.depth) {
            mt.rootHash = copyBytes(node.Hash)
        } else {
            parentNode, err := NewInternalNode(node.Hash, nil)
            if err != nil {
                return err
            }
            node.setParentNodeHash(parentNode.Hash)
            return mt.addNode(level+1, parentNode)
        }
    } else if hangingNode.Parent == nil {
        parent, err := NewInternalNode(hangingNode.Hash, node.Hash)
        if err != nil {
            return err
        }
        hangingNode.setParentNodeHash(parent.Hash)
        node.setParentNodeHash(parent.Hash)
        delete(mt.hangingNodes, level)
        return mt.addNode(level+1, parent)
    } else {
        parentNodeOfHangingNode, err := mt.getNodeByHash(hangingNode.Parent)
        if err != nil {
            return err
        }
        if parentNodeOfHangingNode != nil {
            if err := mt.addLeafToNode(parentNodeOfHangingNode, node.Hash); err != nil {
                return err
            }
            delete(mt.hangingNodes, level)
        } else {
            parent, err := NewInternalNode(hangingNode.Hash, node.Hash)
            if err != nil {
                return err
            }
            hangingNode.setParentNodeHash(parent.Hash)
            node.setParentNodeHash(parent.Hash)
            delete(mt.hangingNodes, level)
            return mt.addNode(level+1, parent)
        }
    }
    
    return nil
}

func (mt *MerkleTree) addLeafToNode(node *Node, leafHash []byte) error {
    leafNode, err := mt.getNodeByHash(leafHash)
    if err != nil {
        return err
    }
    if leafNode == nil {
        return fmt.Errorf("leaf node not found: %s", hex.EncodeToString(leafHash))
    }
    
    node.mu.Lock()
    if node.Left == nil {
        node.Left = copyBytes(leafHash)
    } else if node.Right == nil {
        node.Right = copyBytes(leafHash)
    } else {
        node.mu.Unlock()
        return errors.New("node already has both left and right children")
    }
    node.mu.Unlock()
    
    newHash := node.calculateHash()
    if newHash == nil {
        return errors.New("failed to calculate new hash after adding leaf")
    }
    
    return mt.updateNodeHash(node, newHash)
}

func (mt *MerkleTree) updateLeaf(oldLeafHash, newLeafHash []byte) error {
    if oldLeafHash == nil {
        return errors.New("old leaf hash cannot be nil")
    }
    if newLeafHash == nil {
        return errors.New("new leaf hash cannot be nil")
    }
    if bytes.Equal(oldLeafHash, newLeafHash) {
        return errors.New("old and new leaf hashes cannot be the same")
    }
    
    leaf, err := mt.getNodeByHash(oldLeafHash)
    if err != nil {
        return err
    }
    
    if leaf == nil {
        return fmt.Errorf("leaf not found: %s", hex.EncodeToString(oldLeafHash))
    }
    
    return mt.updateNodeHash(leaf, newLeafHash)
}

func (mt *MerkleTree) updateNodeHash(node *Node, newHash []byte) error {
    node.mu.Lock()
    
    // Store old hash for deletion if not already set
    if node.NodeHashToRemoveFromDB == nil {
        node.NodeHashToRemoveFromDB = make([]byte, len(node.Hash))
        copy(node.NodeHashToRemoveFromDB, node.Hash)
    }
    
    oldHash := make([]byte, len(node.Hash))
    copy(oldHash, node.Hash)
    node.Hash = make([]byte, len(newHash))
    copy(node.Hash, newHash)
    
    node.mu.Unlock()
    
    // Update hanging nodes
    for level, hash := range mt.hangingNodes {
        if bytes.Equal(hash, oldHash) {
            mt.hangingNodes[level] = copyBytes(newHash)
            break
        }
    }
    
    // Update cache
    oldHashHex := hex.EncodeToString(oldHash)
    newHashHex := hex.EncodeToString(newHash)
    delete(mt.nodesCache, oldHashHex)
    mt.nodesCache[newHashHex] = node
    
    // Handle different node types
    isLeaf := node.Left == nil && node.Right == nil
    isRoot := node.Parent == nil
    
    if isRoot {
        mt.rootHash = copyBytes(newHash)
        
        // Update children's parent references
        if node.Left != nil {
            leftNode, err := mt.getNodeByHash(node.Left)
            if err == nil && leftNode != nil {
                leftNode.setParentNodeHash(newHash)
            }
        }
        if node.Right != nil {
            rightNode, err := mt.getNodeByHash(node.Right)
            if err == nil && rightNode != nil {
                rightNode.setParentNodeHash(newHash)
            }
        }
    }
    
    if isLeaf && !isRoot {
        parentNode, err := mt.getNodeByHash(node.Parent)
        if err != nil {
            return err
        }
        if parentNode != nil {
            if err := parentNode.updateLeaf(oldHash, newHash); err != nil {
                return err
            }
            newParentHash := parentNode.calculateHash()
            return mt.updateNodeHash(parentNode, newParentHash)
        }
    } else if !isLeaf && !isRoot {
        // Update children's parent references
        if node.Left != nil {
            leftNode, err := mt.getNodeByHash(node.Left)
            if err == nil && leftNode != nil {
                leftNode.setParentNodeHash(newHash)
            }
        }
        if node.Right != nil {
            rightNode, err := mt.getNodeByHash(node.Right)
            if err == nil && rightNode != nil {
                rightNode.setParentNodeHash(newHash)
            }
        }
        
        // Update parent
        parentNode, err := mt.getNodeByHash(node.Parent)
        if err != nil {
            return err
        }
        if parentNode != nil {
            if err := parentNode.updateLeaf(oldHash, newHash); err != nil {
                return err
            }
            newParentHash := parentNode.calculateHash()
            return mt.updateNodeHash(parentNode, newParentHash)
        }
    }
    
    return nil
}

// ContainsKey checks if a key exists in the tree
func (mt *MerkleTree) ContainsKey(key []byte) (bool, error) {
    if err := mt.errorIfClosed(); err != nil {
        return false, err
    }
    
    if key == nil {
        return false, errors.New("key cannot be nil")
    }
    
    data, err := mt.GetData(key)
    if err != nil {
        return false, err
    }
    
    return data != nil, nil
}

// FlushToDisk flushes all in-memory changes to BoltDB
func (mt *MerkleTree) FlushToDisk() error {
    if atomic.LoadInt32(&mt.hasUnsavedChanges) == 0 {
        return nil
    }
    
    if err := mt.errorIfClosed(); err != nil {
        return err
    }
    
    mt.mu.Lock()
    defer mt.mu.Unlock()
    
    return mt.db.Update(func(tx *bbolt.Tx) error {
        // Clear existing metadata
        metaBucket := tx.Bucket([]byte(METADATA_BUCKET))
        if metaBucket != nil {
            // Delete all existing metadata
            c := metaBucket.Cursor()
            for k, _ := c.First(); k != nil; k, _ = c.Next() {
                if err := c.Delete(); err != nil {
                    return err
                }
            }
        }
        
        // Write metadata
        if mt.rootHash != nil {
            if err := metaBucket.Put([]byte(KEY_ROOT_HASH), mt.rootHash); err != nil {
                return err
            }
        }
        
        numLeavesBytes := make([]byte, 4)
        binary.BigEndian.PutUint32(numLeavesBytes, uint32(atomic.LoadInt32(&mt.numLeaves)))
        if err := metaBucket.Put([]byte(KEY_NUM_LEAVES), numLeavesBytes); err != nil {
            return err
        }
        
        depthBytes := make([]byte, 4)
        binary.BigEndian.PutUint32(depthBytes, uint32(atomic.LoadInt32(&mt.depth)))
        if err := metaBucket.Put([]byte(KEY_DEPTH), depthBytes); err != nil {
            return err
        }
        
        // Write hanging nodes
        for level, nodeHash := range mt.hangingNodes {
            key := fmt.Sprintf("%s%d", KEY_HANGING_NODE_PREFIX, level)
            if err := metaBucket.Put([]byte(key), nodeHash); err != nil {
                return err
            }
        }
        
        // Write nodes
        nodesBucket := tx.Bucket([]byte(NODES_BUCKET))
        for _, node := range mt.nodesCache {
            encoded, err := node.encode()
            if err != nil {
                return err
            }
            if err := nodesBucket.Put(node.Hash, encoded); err != nil {
                return err
            }
            
            // Delete old node if needed
            if node.NodeHashToRemoveFromDB != nil {
                if err := nodesBucket.Delete(node.NodeHashToRemoveFromDB); err != nil {
                    return err
                }
            }
        }
        
        // Write key data
        keyDataBucket := tx.Bucket([]byte(KEY_DATA_BUCKET))
        for keyHex, data := range mt.keyDataCache {
            key, err := hex.DecodeString(keyHex)
            if err != nil {
                continue
            }
            if err := keyDataBucket.Put(key, data); err != nil {
                return err
            }
        }
        
        return nil
    })
    
    // Clear caches after successful write
    mt.nodesCache = make(map[string]*Node)
    mt.keyDataCache = make(map[string][]byte)
    atomic.StoreInt32(&mt.hasUnsavedChanges, 0)

    return nil
}

// RevertUnsavedChanges reverts all unsaved changes
func (mt *MerkleTree) RevertUnsavedChanges() error {
    if atomic.LoadInt32(&mt.hasUnsavedChanges) == 0 {
        return nil
    }
    
    if err := mt.errorIfClosed(); err != nil {
        return err
    }
    
    mt.mu.Lock()
    defer mt.mu.Unlock()
    
    mt.nodesCache = make(map[string]*Node)
    mt.hangingNodes = make(map[int][]byte)
    mt.keyDataCache = make(map[string][]byte)
    
    if err := mt.loadMetaData(); err != nil {
        return err
    }
    
    atomic.StoreInt32(&mt.hasUnsavedChanges, 0)
    return nil
}

// Clear efficiently clears the entire MerkleTree
func (mt *MerkleTree) Clear() error {
    if err := mt.errorIfClosed(); err != nil {
        return err
    }
    
    mt.mu.Lock()
    defer mt.mu.Unlock()
    
    // Clear all buckets
    err := mt.db.Update(func(tx *bbolt.Tx) error {
        // Delete and recreate buckets
        if err := tx.DeleteBucket([]byte(METADATA_BUCKET)); err != nil && err != bbolt.ErrBucketNotFound {
            return err
        }
        if err := tx.DeleteBucket([]byte(NODES_BUCKET)); err != nil && err != bbolt.ErrBucketNotFound {
            return err
        }
        if err := tx.DeleteBucket([]byte(KEY_DATA_BUCKET)); err != nil && err != bbolt.ErrBucketNotFound {
            return err
        }
        
        // Recreate buckets
        if _, err := tx.CreateBucket([]byte(METADATA_BUCKET)); err != nil {
            return err
        }
        if _, err := tx.CreateBucket([]byte(NODES_BUCKET)); err != nil {
            return err
        }
        if _, err := tx.CreateBucket([]byte(KEY_DATA_BUCKET)); err != nil {
            return err
        }
        
        return nil
    })
    
    if err != nil {
        return err
    }
    
    // Reset in-memory state
    mt.nodesCache = make(map[string]*Node)
    mt.keyDataCache = make(map[string][]byte)
    mt.hangingNodes = make(map[int][]byte)
    mt.rootHash = nil
    atomic.StoreInt32(&mt.numLeaves, 0)
    atomic.StoreInt32(&mt.depth, 0)
    atomic.StoreInt32(&mt.hasUnsavedChanges, 0)
    
    return nil
}

// Close closes the database
func (mt *MerkleTree) Close() error {
    mt.mu.Lock()
    defer mt.mu.Unlock()
    
    if atomic.LoadInt32(&mt.closed) != 0 {
        return nil
    }
    
    // Flush any pending changes
    if atomic.LoadInt32(&mt.hasUnsavedChanges) != 0 {
        if err := mt.flushToDiskLocked(); err != nil {
            // Log error but continue with cleanup
            fmt.Printf("Error flushing to disk during close: %v\n", err)
        }
    }
    
    // Close database
    if mt.db != nil {
        if err := mt.db.Close(); err != nil {
            fmt.Printf("Error closing database: %v\n", err)
        }
    }
    
    // Remove from open trees
    treesMutex.Lock()
    delete(openTrees, mt.treeName)
    treesMutex.Unlock()
    
    atomic.StoreInt32(&mt.closed, 1)
    return nil
}

// flushToDiskLocked is a helper that assumes the mutex is already held
func (mt *MerkleTree) flushToDiskLocked() error {
    if atomic.LoadInt32(&mt.hasUnsavedChanges) == 0 {
        return nil
    }
    
    err := mt.db.Update(func(tx *bbolt.Tx) error {
        // Clear existing metadata
        metaBucket := tx.Bucket([]byte(METADATA_BUCKET))
        if metaBucket != nil {
            // Delete all existing metadata
            c := metaBucket.Cursor()
            for k, _ := c.First(); k != nil; k, _ = c.Next() {
                if err := c.Delete(); err != nil {
                    return err
                }
            }
        }
        
        // Write metadata
        if mt.rootHash != nil {
            if err := metaBucket.Put([]byte(KEY_ROOT_HASH), mt.rootHash); err != nil {
                return err
            }
        }
        
        numLeavesBytes := make([]byte, 4)
        binary.BigEndian.PutUint32(numLeavesBytes, uint32(atomic.LoadInt32(&mt.numLeaves)))
        if err := metaBucket.Put([]byte(KEY_NUM_LEAVES), numLeavesBytes); err != nil {
            return err
        }
        
        depthBytes := make([]byte, 4)
        binary.BigEndian.PutUint32(depthBytes, uint32(atomic.LoadInt32(&mt.depth)))
        if err := metaBucket.Put([]byte(KEY_DEPTH), depthBytes); err != nil {
            return err
        }
        
        // Write hanging nodes
        for level, nodeHash := range mt.hangingNodes {
            key := fmt.Sprintf("%s%d", KEY_HANGING_NODE_PREFIX, level)
            if err := metaBucket.Put([]byte(key), nodeHash); err != nil {
                return err
            }
        }
        
        // Write nodes
        nodesBucket := tx.Bucket([]byte(NODES_BUCKET))
        for _, node := range mt.nodesCache {
            encoded, err := node.encode()
            if err != nil {
                return err
            }
            if err := nodesBucket.Put(node.Hash, encoded); err != nil {
                return err
            }
            
            // Delete old node if needed
            if node.NodeHashToRemoveFromDB != nil {
                if err := nodesBucket.Delete(node.NodeHashToRemoveFromDB); err != nil {
                    return err
                }
            }
        }
        
        // Write key data
        keyDataBucket := tx.Bucket([]byte(KEY_DATA_BUCKET))
        for keyHex, data := range mt.keyDataCache {
            key, err := hex.DecodeString(keyHex)
            if err != nil {
                continue
            }
            if err := keyDataBucket.Put(key, data); err != nil {
                return err
            }
        }
        
        return nil
    })
    
    if err != nil {
        return err
    }
    
    // Clear caches after successful write
    mt.nodesCache = make(map[string]*Node)
    mt.keyDataCache = make(map[string][]byte)
    atomic.StoreInt32(&mt.hasUnsavedChanges, 0)
    
    return nil
}

// GetTreeName returns the tree name
func (mt *MerkleTree) GetTreeName() string {
    return mt.treeName
}

// GetPath returns the tree path
func (mt *MerkleTree) GetPath() string {
    return mt.path
}
