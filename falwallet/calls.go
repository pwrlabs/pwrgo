package falwallet

import (
	"os"
	"encoding/hex"
	"path/filepath"
	"encoding/binary"

	"github.com/pwrlabs/pwrgo/encode"
	"github.com/pwrlabs/pwrgo/rpc"
)

func (w *Falcon512Wallet) GetAddress() string {
	return "0x" + hex.EncodeToString(w.Address)
}

func (w *Falcon512Wallet) GetPublicKey() []byte {
	return w.PublicKey
}

func (w *Falcon512Wallet) GetPrivateKey() []byte {
	return w.PrivateKey
}

func (w *Falcon512Wallet) GetBalance() int {
	balance := rpc.GetBalanceOfAddress(w.GetAddress())
	return balance
}

func (w *Falcon512Wallet) GetNonce() int {
	nonce := rpc.GetNonceOfAddress(w.GetAddress())
	return nonce
}

// Sign signs a message with the wallet's private key
func (w *Falcon512Wallet) SignTx(buffer []byte) ([]byte, error) {
	signature, err := encode.Sign(buffer, w.PrivateKey, encode.SigCompressed)
	if err != nil {
		return nil, err
	}

	signatureLenBytes := encode.DecToBytes(len(signature), 2)

	txn_bytes := make([]byte, len(buffer))
	copy(txn_bytes, buffer)

	txn_bytes = append(txn_bytes, signatureLenBytes...)
	txn_bytes = append(txn_bytes, signature...)

	return txn_bytes, nil
}

// VerifySign verifies a signature against the wallet's public key
func (w *Falcon512Wallet) VerifySign(message, signature []byte) bool {
	err := encode.Verify(signature, message, w.PublicKey, encode.SigCompressed)
	return err == nil
}

// StoreWallet stores the wallet to a file
func (w *Falcon512Wallet) StoreWallet(filePath string) error {
	var buffer []byte

	// Add public key length and data
	pubKeyLengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(pubKeyLengthBytes, uint32(len(w.PublicKey)))
	buffer = append(buffer, pubKeyLengthBytes...)
	buffer = append(buffer, w.PublicKey...)

	// Add private key length and data
	privKeyLengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(privKeyLengthBytes, uint32(len(w.PrivateKey)))
	buffer = append(buffer, privKeyLengthBytes...)
	buffer = append(buffer, w.PrivateKey...)

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(filePath, buffer, 0600)
}
