package wallet

import (
	"github.com/pwrlabs/pwrgo/rpc"
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
	nonce := rpc.GetNonceOfAddress(w.address)
	return nonce
}

func (w *PWRWallet) GetBalance() int {
	nonce := rpc.GetBalanceOfAddress(w.address)
	return nonce
}
