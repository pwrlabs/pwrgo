package main

import (
	"fmt"

	"github.com/pwrlabs/pwrgo/rpc"
	"github.com/pwrlabs/pwrgo/wallet"
)

func main() {
	wallets, _ := wallet.New("demand april length soap cash concert shuffle result force mention fringe slim")
	wallets.StoreWallet("example_wallet.dat", "your_password_here")

	wallet, _ := wallet.LoadWallet("example_wallet.dat", "your_password_here")
	fmt.Println("Address:", wallet.GetAddress())
	fmt.Println("Nonce:", wallet.GetNonce())
	fmt.Println("Balance:", wallet.GetBalance())

	var tx rpc.BroadcastResponse
	tx = wallet.TransferPWR("0x2bf2ff7d9ace8aef8a21726242e7e2d323f0d5d5", 1, wallet.GetRpc().GetFeeBerByte())
	if tx.Success {
		fmt.Println("TX HASH:", tx.Hash)
	} else {
		fmt.Println("TX Error:", tx.Error)
	}

	vidaId := 123
    data := []byte("Hello world")
    feePerByte := wallet.GetRpc().GetFeeBerByte()

    tx = wallet.SendVidaData(vidaId, data, feePerByte)
    if tx.Success {
        fmt.Printf("Sending tx hash: %s\n", tx.Hash)
    } else {
        fmt.Println("Error sending VIDA data tx:", tx.Error)
    }

	tx = wallet.SendPayableVidaData(vidaId, data, 1000, feePerByte)
    if tx.Success {
        fmt.Printf("Sending tx hash: %s\n", tx.Hash)
    } else {
        fmt.Println("Error sending VIDA data tx:", tx.Error)
    }
}
