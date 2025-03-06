package wallet

import (
	"crypto/ecdsa"
)

type PWRWallet struct {
    privateKey    *ecdsa.PrivateKey
	privateKeyStr string
	publicKey     string
	address       string
}
