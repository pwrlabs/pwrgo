package falwallet

import (
	"github.com/pwrlabs/pwrgo/rpc"
)

type Falcon512Wallet struct {
	PublicKey  []byte
	PrivateKey []byte
	Address    []byte
	rpc        *rpc.RPC
}
