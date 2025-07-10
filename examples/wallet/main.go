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

    feePerByte := wallet.GetRpc().GetFeePerByte()

	var tx rpc.BroadcastResponse
	tx = wallet.TransferPWR("0x2bf2ff7d9ace8aef8a21726242e7e2d323f0d5d5", 1, feePerByte)
	if tx.Success {
		fmt.Println("Transfer tx hash:", tx.Hash)
	} else {
		fmt.Println("Transfer tx Error:", tx.Error)
	}

	vidaId := 123
    data := []byte("Hello world")

    tx = wallet.SendVidaData(vidaId, data, feePerByte)
    if tx.Success {
        fmt.Printf("Send VIDA Data tx hash: %s\n", tx.Hash)
    } else {
        fmt.Println("Error sending VIDA data tx:", tx.Error)
    }

	tx = wallet.SendPayableVidaData(vidaId, data, 1000, feePerByte)
    if tx.Success {
        fmt.Printf("Send Payable VIDA Data tx hash: %s\n", tx.Hash)
    } else {
        fmt.Println("Error sending payable VIDA data tx:", tx.Error)
    }
}
