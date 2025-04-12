package wallet

import (
	"crypto/ecdsa"

	"github.com/pwrlabs/pwrgo/rpc"
)

type PWRWallet struct {
    privateKey    *ecdsa.PrivateKey
	privateKeyStr string
	publicKey     string
	address       string
	rpc           *rpc.RPC
}
