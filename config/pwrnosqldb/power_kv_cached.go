package pwrnosqldb

import (
	"bytes"
	"context"
	"encoding/hex"
	"log"
	"sync"
	"time"
)

// ByteArrayWrapper wraps []byte to use as map key
type ByteArrayWrapper struct {
	data []byte
	hash string
}

// NewByteArrayWrapper creates a new ByteArrayWrapper
func NewByteArrayWrapper(data []byte) *ByteArrayWrapper {
	return &ByteArrayWrapper{
		data: data,
		hash: hex.EncodeToString(data),
	}
}

// Equals checks if two ByteArrayWrapper are equal
func (b *ByteArrayWrapper) Equals(other *ByteArrayWrapper) bool {
	return bytes.Equal(b.data, other.data)
}

// String returns the hex representation
func (b *ByteArrayWrapper) String() string {
	return b.hash
}

// PowerKvCached represents the cached PowerKv client
type PowerKvCached struct {
	db           *PowerKv
	cache        sync.Map // map[string][]byte
	isShutdown   bool
	shutdownMu   sync.RWMutex
	activeWrites sync.WaitGroup
	ctx          context.Context
	cancel       context.CancelFunc
}

// NewPowerKvCached creates a new PowerKvCached instance
func NewPowerKvCached(projectID, secret string) (*PowerKvCached, error) {
	db, err := NewPowerKv(projectID, secret)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &PowerKvCached{
		db:     db,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

// toBytes converts various types to bytes (same as PowerKv)
func (p *PowerKvCached) toBytes(data interface{}) ([]byte, error) {
	return p.db.toBytes(data)
}

// Put stores data with the given key (non-blocking)
func (p *PowerKvCached) Put(key, value interface{}) error {
	p.shutdownMu.RLock()
	if p.isShutdown {
		p.shutdownMu.RUnlock()
		return NewPowerKvError("InvalidInput", "PowerKvCached has been shut down")
	}
	p.shutdownMu.RUnlock()

	keyBytes, err := p.toBytes(key)
	if err != nil {
		return err
	}

	valueBytes, err := p.toBytes(value)
	if err != nil {
		return err
	}

	return p.PutBytes(keyBytes, valueBytes)
}

// PutBytes stores byte data with byte key (non-blocking)
func (p *PowerKvCached) PutBytes(key, value []byte) error {
	p.shutdownMu.RLock()
	if p.isShutdown {
		p.shutdownMu.RUnlock()
		return NewPowerKvError("InvalidInput", "PowerKvCached has been shut down")
	}
	p.shutdownMu.RUnlock()

	keyWrapper := NewByteArrayWrapper(key)
	keyStr := keyWrapper.String()

	// Check if we have an old value
	oldValueInterface, hadOldValue := p.cache.Load(keyStr)
	var oldValue []byte
	if hadOldValue {
		oldValue = oldValueInterface.([]byte)
	}

	// Update cache immediately
	p.cache.Store(keyStr, value)

	// If oldValue is same as new value, no need to update db
	// If oldValue doesn't exist, it means this key is being inserted for the first time, so we need to update db
	if !hadOldValue || !bytes.Equal(oldValue, value) {
		// Start background write (non-blocking)
		p.activeWrites.Add(1)
		go p.backgroundWrite(key, value, keyStr)
	}

	return nil
}

// backgroundWrite performs background write with retry logic
func (p *PowerKvCached) backgroundWrite(keyBytes, valueBytes []byte, keyStr string) {
	defer p.activeWrites.Done()

	for {
		// Check if shutdown was requested
		select {
		case <-p.ctx.Done():
			return
		default:
		}

		// Check if cache value has changed
		currentCachedInterface, exists := p.cache.Load(keyStr)
		if !exists {
			log.Printf("Cache updated for key, stopping background write: %s", string(keyBytes))
			return
		}

		currentCached := currentCachedInterface.([]byte)
		if !bytes.Equal(currentCached, valueBytes) {
			log.Printf("Cache updated for key, stopping background write: %s", string(keyBytes))
			return
		}

		// Try to write to database
		success, err := p.db.PutBytes(keyBytes, valueBytes)
		if err == nil && success {
			log.Printf("Successfully updated key on PWR Chain: %s", string(keyBytes))
			return
		}

		if err != nil {
			log.Printf("Error updating key on PWR Chain: %s - %v", string(keyBytes), err)
		} else {
			log.Printf("Failed to update key on PWR Chain, retrying: %s", string(keyBytes))

			// Check if another goroutine has already updated the value
			remoteValue, remoteErr := p.db.GetValueBytes(keyBytes)
			if remoteErr == nil && bytes.Equal(remoteValue, valueBytes) {
				log.Printf("Value already updated by another process: %s", string(keyBytes))
				return
			}
		}

		select {
		case <-p.ctx.Done():
			return
		case <-time.After(10 * time.Millisecond):
			// Continue to retry
		}
	}
}

// GetValue retrieves data for the given key
func (p *PowerKvCached) GetValue(key interface{}) ([]byte, error) {
	keyBytes, err := p.toBytes(key)
	if err != nil {
		return nil, err
	}

	return p.GetValueBytes(keyBytes)
}

// GetValueBytes retrieves data for the given byte key
func (p *PowerKvCached) GetValueBytes(key []byte) ([]byte, error) {
	keyWrapper := NewByteArrayWrapper(key)
	keyStr := keyWrapper.String()

	// Check cache first
	if cachedInterface, exists := p.cache.Load(keyStr); exists {
		return cachedInterface.([]byte), nil
	}

	// If not in cache, fetch from remote
	value, err := p.db.GetValueBytes(key)
	if err != nil {
		log.Printf("Error retrieving value: %v", err)
		return nil, nil
	}

	if value != nil {
		// Cache the retrieved value
		p.cache.Store(keyStr, value)
	}

	return value, nil
}

// GetStringValue retrieves data as string
func (p *PowerKvCached) GetStringValue(key interface{}) (string, error) {
	value, err := p.GetValue(key)
	if err != nil || value == nil {
		return "", err
	}
	return string(value), nil
}

// GetIntValue retrieves data as int
func (p *PowerKvCached) GetIntValue(key interface{}) (int, error) {
	stringValue, err := p.GetStringValue(key)
	if err != nil || stringValue == "" {
		return 0, err
	}
	return p.db.GetIntValue(key)
}

// GetLongValue retrieves data as int64
func (p *PowerKvCached) GetLongValue(key interface{}) (int64, error) {
	stringValue, err := p.GetStringValue(key)
	if err != nil || stringValue == "" {
		return 0, err
	}
	return p.db.GetLongValue(key)
}

// GetDoubleValue retrieves data as float64
func (p *PowerKvCached) GetDoubleValue(key interface{}) (float64, error) {
	stringValue, err := p.GetStringValue(key)
	if err != nil || stringValue == "" {
		return 0, err
	}
	return p.db.GetDoubleValue(key)
}

// Shutdown gracefully shuts down the cached client
func (p *PowerKvCached) Shutdown() error {
	log.Println("Shutting down PowerKvCached...")

	p.shutdownMu.Lock()
	p.isShutdown = true
	p.shutdownMu.Unlock()

	// Cancel context to signal all goroutines to stop
	p.cancel()

	// Wait for all active writes to complete
	done := make(chan struct{})
	go func() {
		p.activeWrites.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("All background writes completed")
	case <-time.After(60 * time.Second):
		log.Println("Forced shutdown after 60 seconds timeout")
	}

	return nil
}
