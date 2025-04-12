package wallet

import (
	"os"
	"github.com/pwrlabs/pwrgo/encode"
)

func (w *PWRWallet) GetAddress() string {
	return w.address
}

func (w *PWRWallet) GetPrivateKey() string {
	return w.privateKeyStr
}

func (w *PWRWallet) GetPublicKey() string {
	return w.publicKey
}

func (w *PWRWallet) GetNonce() int {
	nonce := w.rpc.GetNonceOfAddress(w.address)
	return nonce
}

func (w *PWRWallet) GetBalance() int {
	nonce := w.rpc.GetBalanceOfAddress(w.address)
	return nonce
}

func (w *PWRWallet) StoreWallet(path string, password string) error {
	privateKeyBytes := w.privateKey.D.Bytes()
	encryptedData, err := encode.Encrypt(privateKeyBytes, password)
	if err != nil {
		return err
	}

	return os.WriteFile(path, encryptedData, 0600)
}
