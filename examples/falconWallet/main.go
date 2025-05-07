package main

import (
	"fmt"

	"github.com/pwrlabs/pwrgo/rpc"
	"github.com/pwrlabs/pwrgo/wallet"
)

func main() {
	wallets, _ := wallet.New("your seed phrase here")
	wallets.StoreWallet("example_wallet.dat", "your_password_here")

	wallet, _ := wallet.LoadWallet("example_wallet.dat", "your_password_here")
	fmt.Println("Address:", wallet.GetAddress())
	fmt.Println("Balance:", wallet.GetBalance())

	var tx rpc.BroadcastResponse
	tx = wallet.TransferPWR("0x2bf2ff7d9ace8aef8a21726242e7e2d323f0d5d5", 1, wallet.GetRpc().GetFeeBerByte())
	if tx.Success {
		fmt.Println("TX HASH:", tx.Hash)
	} else {
		fmt.Println("TX Error:", tx.Error)
	}
}
