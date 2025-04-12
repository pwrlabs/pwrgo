package main

import (
	"fmt"

	"github.com/pwrlabs/pwrgo/falwallet"
	"github.com/pwrlabs/pwrgo/rpc"
)

func main() {
	// wallet, _ := falwallet.New()
	// wallet.StoreWallet("new_wallet.dat")

	wallet, _ := falwallet.LoadWallet("new_wallet.dat")
	fmt.Println("Address:", wallet.GetAddress())

	pwr := rpc.SetRpcNodeUrl("https://pwrrpc.pwrlabs.io")

	var tx rpc.BroadcastResponse
	tx = wallet.TransferPWR("0x2bf2ff7d9ace8aef8a21726242e7e2d323f0d5d5", 1, pwr.GetFeeBerByte())
	if tx.Success {
		fmt.Println("TX HASH:", tx.Hash)
	} else {
		fmt.Println("TX Error:", tx.Error)
	}
}
