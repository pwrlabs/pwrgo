package wallet

import (
	// "encoding/hex"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/ebfe/keccak"
	"github.com/pwrlabs/pwrgo/encode"
	"github.com/pwrlabs/pwrgo/rpc"
)

// NewRandom creates a new PWRWallet with a generated mnemonic phrase
func NewRandom(wordCount int, rpcEndpoint ...*rpc.RPC) (*PWRWallet, error) {
	// Calculate entropy bytes based on word count
	var entropyBytes int
	switch wordCount {
	case 12:
		entropyBytes = 16 // 128 bits
	case 15:
		entropyBytes = 20 // 160 bits
	case 18:
		entropyBytes = 24 // 192 bits
	case 21:
		entropyBytes = 28 // 224 bits
	case 24:
		entropyBytes = 32 // 256 bits
	default:
		return nil, fmt.Errorf("invalid word count: %d", wordCount)
	}
	_ = entropyBytes

	publicKey, secretKey, seedPhrase := generateRandomKeypair(wordCount)

	publicKeyBytes, _ := hex.DecodeString(publicKey)
	secretKeyBytes, _ := hex.DecodeString(secretKey)

	wallet, err := FromKeys([]byte(seedPhrase), publicKeyBytes, secretKeyBytes, rpcEndpoint...)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	return wallet, nil
}

func New(seedPhrase string, rpcEndpoint ...*rpc.RPC) (*PWRWallet, error) {
	publicKey, secretKey := generateKeypair(seedPhrase)

	publicKeyBytes, _ := hex.DecodeString(publicKey)
	secretKeyBytes, _ := hex.DecodeString(secretKey)
	
	wallet, _ := FromKeys(
		[]byte(seedPhrase), publicKeyBytes, secretKeyBytes, rpcEndpoint...,
	)

	return wallet, nil
}

// FromKeys creates a wallet from existing keys
func FromKeys(seedPhrase []byte, publicKey, privateKey []byte, rpcEndpoint ...*rpc.RPC) (*PWRWallet, error) {
	// Get the hash of the public key
	hash := hash224(publicKey[1:])
	address := hash[:20]

	endpoint := "https://pwrrpc.pwrlabs.io"
	if len(rpcEndpoint) > 0 {
		endpoint = rpcEndpoint[0].GetRpcNodeUrl()
	}

	return &PWRWallet{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
		Address:    address,
		seedPhrase: seedPhrase,
		rpc:        rpc.SetRpcNodeUrl(endpoint),
	}, nil
}

// LoadWallet loads a wallet from a file
func LoadWallet(path string, password string, rpcEndpoint ...*rpc.RPC) (*PWRWallet, error) {
	encryptedData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	seedPhraseBytes, err := encode.Decrypt(encryptedData, password)
	if err != nil {
		return nil, err
	}

	seedPhrase := string(seedPhraseBytes)
	wallet, err := New(seedPhrase, rpcEndpoint...)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func hash224(input []byte) []byte {
	hasher := keccak.New224()
	hasher.Write(input)
	return hasher.Sum(nil)
}
