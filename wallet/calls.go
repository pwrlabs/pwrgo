package wallet

import (
	"encoding/hex"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/keep-pwr-strong/falcon-go/falcon"
	"github.com/pwrlabs/pwrgo/config/aes256"
	"github.com/pwrlabs/pwrgo/config/transactions"
	"github.com/pwrlabs/pwrgo/rpc"
)

func (w *PWRWallet) GetAddress() string {
	return "0x" + hex.EncodeToString(w.Address)
}

func (w *PWRWallet) GetPublicKey() []byte {
	return w.PublicKey
}

func (w *PWRWallet) GetPrivateKey() []byte {
	return w.PrivateKey
}

func (w *PWRWallet) GetBalance() int {
	balance := w.rpc.GetBalanceOfAddress(w.GetAddress())
	return balance
}

func (w *PWRWallet) GetNonce() int {
	nonce := w.rpc.GetNonceOfAddress(w.GetAddress())
	return nonce
}

// Sign signs a message with the wallet's private key
func (w *PWRWallet) SignTx(transaction []byte) ([]byte, error) {
	txnHash := crypto.Keccak256Hash(transaction)

	signature, err := falcon.Sign(txnHash.Bytes(), w.PrivateKey, falcon.SigCompressed)
	if err != nil {
		return nil, err
	}

	signatureLenBytes := transactions.DecToBytes(len(signature), 2)

	txn_bytes := make([]byte, len(transaction))
	copy(txn_bytes, transaction)

	txn_bytes = append(txn_bytes, signature...)
	txn_bytes = append(txn_bytes, signatureLenBytes...)

	return txn_bytes, nil
}

// VerifySign verifies a signature against the wallet's public key
func (w *PWRWallet) VerifySign(message, signature []byte) bool {
	err := falcon.Verify(signature, message, w.PublicKey, falcon.SigCompressed)
	return err == nil
}

// StoreWallet stores the wallet to a file
func (w *PWRWallet) StoreWallet(filePath string, password string) error {
	var seedPhrase []byte = w.seedPhrase
	encryptedData, err := aes256.Encrypt(seedPhrase, password)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, encryptedData, 0600)
}

func (w *PWRWallet) GetRpc() *rpc.RPC {
	return w.rpc
}

func (w *PWRWallet) GetSeedPhrase() string {
	return string(w.seedPhrase)
}
