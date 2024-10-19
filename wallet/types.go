package wallet

import (
	"crypto/ecdsa"
)

type PWRWallet struct {
    privateKey *ecdsa.PrivateKey
	privateKeyStr string
	publicKey string
	address string
}

type BroadcastResponse struct {
    Message string   `json:"message,omitempty"`
    Success bool
    TxHash string
    Error string
}
