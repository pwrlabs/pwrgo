package main

import (
    "fmt"
    "github.com/pwrlabs/pwrgo/rpc"
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

	var delegators = rpc.GetDelegatorsOfValidator("0x7b6F32435084Cab827f0ce7Af1C0D48600CE3CaD")
	fmt.Println("Delegators:", delegators)
    
    var blocksCount = rpc.GetBlocksCount()
    fmt.Println("Blocks count:", blocksCount)

    var latestBlockCount = rpc.GetLatestBlockNumber()
    fmt.Println("Validators count:", latestBlockCount)

	var startBlcok = 843500
	var endBlock = 843750
	var vmId = 123
	var transactions = rpc.GetVmDataTransactions(startBlcok, endBlock, vmId)
	fmt.Println("VM Data:", transactions)

	var guardian = rpc.GetGuardianOfAddress("0xD97C25C0842704588DD70A061C09A522699E2B9C")
    fmt.Println("Guardian:", guardian)

	var block = rpc.GetBlockByNumber(836599)
    fmt.Println("Block:", block)

	var activeVotingPower = rpc.GetActiveVotingPower()
	fmt.Println("ActiveVotingPower:", activeVotingPower)

	var conduitsVm = rpc.GetConduitsOfVmId(101001)
	fmt.Println("ConduitsVm:", conduitsVm)

	var totalValidatorsCount = rpc.GetValidatorsCount()
	fmt.Println("TotalValidatorsCount:", totalValidatorsCount)

    var transferTx = wallet.TransferPWR("0x3B3B69093879E7B6F28366FA3C32762590FF547E", 1000)
	if transferTx.Success {
		fmt.Printf("Transfer tx hash: %s\n", transferTx.TxHash)
	} else {
		fmt.Println("Error sending Transfer tx:", transferTx.Error)
	}
	
	var tx = rpc.GetTransactionByHash("0xe8be5c174a3457ba1015bef3399bb71a4081b8c24bdbebb05a8aea746ac32486")
	fmt.Println("Transfer TX: ", tx)

    // var data = []byte("Hello world")
    // var vmTxResponses = wallet.SendVMData(1234, data)
	// if vmTxResponses.Success {
	// 	fmt.Printf("Sending tx hash: %s\n", vmTxResponses.TxHash)
	// } else {
	// 	fmt.Println("Error sending VM data tx:", vmTxResponses.Error)
	// }

	// var data = []byte("Hello world")
	// var amount = 1000
    // var vmTxResponses = wallet.SendPayableVMData(vmId, amount, data)
	// if vmTxResponses.Success {
	// 	fmt.Printf("Sending tx hash: %s\n", vmTxResponses.TxHash)
	// } else {
	// 	fmt.Println("Error sending VM data tx:", vmTxResponses.Error)
	// }

	// conduits := []string{"0x7EbFBd2BABA5F68F720C059d62eFc4aaCFA66513"}
	// var vmIds = 101001
	// var setConduit = wallet.SetConduits(vmIds, conduits)
	// if setConduit.Success {
	// 	fmt.Printf("Sending tx hash: %s\n", setConduit.TxHash)
	// } else {
	// 	fmt.Println("Error sending VM data tx:", setConduit.Error)
	// }
}