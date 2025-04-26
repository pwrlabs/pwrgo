package wallet

import (
	"github.com/pwrlabs/pwrgo/rpc"
)

type PWRWallet struct {
	PublicKey  []byte
	PrivateKey []byte
	Address    []byte
	rpc        *rpc.RPC
}
