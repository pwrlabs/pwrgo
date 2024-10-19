package main

import (
    "fmt"
    "github.com/pwrlabs/pwrgo/wallet"
)

func main() {
    // Import wallet by private key
    var privateKeyHex = "0x04828e90065864c111871769c601d7de2246570b39dd37c19ccac16c14b18f72"
    var wallet = wallet.FromPrivateKey(privateKeyHex)
    
	var address = wallet.GetAddress()
    fmt.Printf("Address: %s\n", address)
    
    var nonce = wallet.GetNonce()
    fmt.Println("Nonce:", nonce)

	var balance = wallet.GetBalance()
    fmt.Println("Balance:", balance)

	var data = []byte("Hello world")
	var vmId = 1234
	var amount = 1000

    var transferTx = wallet.TransferPWR("0x3B3B69093879E7B6F28366FA3C32762590FF547E", amount)
	if transferTx.Success {
		fmt.Printf("Transfer tx hash: %s\n", transferTx.TxHash)
	} else {
		fmt.Println("Error sending Transfer tx:", transferTx.Error)
	}

    var vmTx = wallet.SendVMData(vmId, data)
	if vmTx.Success {
		fmt.Printf("Sending tx hash: %s\n", vmTx.TxHash)
	} else {
		fmt.Println("Error sending VM data tx:", vmTx.Error)
	}
	
    var vmTxResponses = wallet.SendPayableVMData(vmId, amount, data)
	if vmTxResponses.Success {
		fmt.Printf("Sending tx hash: %s\n", vmTxResponses.TxHash)
	} else {
		fmt.Println("Error sending VM data tx:", vmTxResponses.Error)
	}

	// conduits := []string{"0x7EbFBd2BABA5F68F720C059d62eFc4aaCFA66513"}
	// var vmIds = 101001
	// var setConduit = wallet.SetConduits(vmIds, conduits)
	// if setConduit.Success {
	// 	fmt.Printf("Sending tx hash: %s\n", setConduit.TxHash)
	// } else {
	// 	fmt.Println("Error sending VM data tx:", setConduit.Error)
	// }
}