package main

import (
	"fmt"

	"github.com/pwrlabs/pwrgo/rpc"
	"github.com/pwrlabs/pwrgo/wallet"
)

func main() {
	// wallet, _ := wallet.New()
	// wallet.StoreWallet("new_wallet.dat")

	wallets, _ := wallet.LoadWallet("new_wallet.dat")
	fmt.Println("Address:", wallets.GetAddress())

	var tx rpc.BroadcastResponse
	tx = wallets.TransferPWR("0x2bf2ff7d9ace8aef8a21726242e7e2d323f0d5d5", 1, wallets.GetRpc().GetFeeBerByte())
	if tx.Success {
		fmt.Println("TX HASH:", tx.Hash)
	} else {
		fmt.Println("TX Error:", tx.Error)
	}
}
