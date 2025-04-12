package falwallet

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"

	"github.com/pwrlabs/pwrgo/encode"
	"github.com/pwrlabs/pwrgo/rpc"
	"github.com/ebfe/keccak"
)

// New creates a new Falcon512Wallet with generated keys
func New(rpcEndpoint ...string) (*Falcon512Wallet, error) {
	keyPair, err := encode.GenerateKeyPair(9) // 9 for Falcon-512
	if err != nil {
		return nil, err
	}

	wallet, _ := FromKeys(
		keyPair.PublicKey, keyPair.PrivateKey, rpcEndpoint...,
	)

	return wallet, nil
}

// FromKeys creates a wallet from existing keys
func FromKeys(publicKey, privateKey []byte, rpcEndpoint ...string) (*Falcon512Wallet, error) {
	// Get the hash of the public key
	hash := hash224(publicKey)
	address := hash[:20]

	endpoint := "https://pwrrpc.pwrlabs.io"
	if len(rpcEndpoint) > 0 {
		endpoint = rpcEndpoint[0]
	}

	return &Falcon512Wallet{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
		Address:    address,
		rpc:        rpc.SetRpcNodeUrl(endpoint),
	}, nil
}

// LoadWallet loads a wallet from a file
func LoadWallet(filePath string, rpcEndpoint ...string) (*Falcon512Wallet, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if len(data) < 8 { // At minimum we need two 4-byte length fields
		return nil, fmt.Errorf("file too small: %d bytes", len(data))
	}

	// Read public key length
	pubLength := binary.BigEndian.Uint32(data[0:4])
	if pubLength == 0 || pubLength > 2048 {
		return nil, fmt.Errorf("invalid public key length: %d", pubLength)
	}

	if 4+pubLength > uint32(len(data)) {
		return nil, fmt.Errorf("file too small for public key of length %d", pubLength)
	}

	publicKeyBytes := data[4 : 4+pubLength]

	// Read private key length
	if 4+pubLength+4 > uint32(len(data)) {
		return nil, errors.New("file too small for secret key length")
	}

	secLength := binary.BigEndian.Uint32(data[4+pubLength : 8+pubLength])
	if secLength == 0 || secLength > 4096 {
		return nil, fmt.Errorf("invalid secret key length: %d", secLength)
	}

	if 8+pubLength+secLength > uint32(len(data)) {
		return nil, fmt.Errorf("file too small for secret key of length %d", secLength)
	}

	privateKeyBytes := data[8+pubLength : 8+pubLength+secLength]

	return FromKeys(publicKeyBytes, privateKeyBytes, rpcEndpoint...)
}

func hash224(input []byte) []byte {
	hasher := keccak.New224()
	hasher.Write(input)
	return hasher.Sum(nil)
}
