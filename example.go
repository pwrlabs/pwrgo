package main

import (
    "fmt"
    "github.com/pwrlabs/pwrgo/pwrgo"
)

func main() {
    // Update RPC URL (default is mainnet)
	pwrgo.SetRpcNodeUrl("https://pwrrpc.pwrlabs.io")

    // Import wallet by private key
    privateKeyHex := "0x9d4428c6e0638331b4866b70c831f8ba51c11b031f4b55eed4087bbb8ef0151f"
    var wallet = pwrgo.FromPrivateKey(privateKeyHex)
    
    fmt.Printf("Public key: %s\n", wallet.PublicKey)
    fmt.Printf("Address: %s\n", wallet.Address)
    
    // Get nonce for address
    var nonce = pwrgo.NonceOfUser(wallet.Address)
    fmt.Println("Nonce: ", nonce)
    
    // Get PWR balance of address
    var balance = pwrgo.BalanceOf(wallet.Address)
    fmt.Println("Balance: ", balance)
    
    // Get total blocks count
    var blocksCount = pwrgo.BlocksCount()
    fmt.Println("Blocks count: ", blocksCount)

    // Get total validators count
    var validatorsCount = pwrgo.ValidatorsCount()
    fmt.Println("Validators count: ", validatorsCount)
    
    // Get block info by Block Number
    var latestBlock = pwrgo.GetBlock(blocksCount - 1)
    fmt.Println("Latest block hash: ", latestBlock.BlockHash)
    fmt.Println("Latest block timestamp: ", latestBlock.Timestamp)
    fmt.Println("Latest block tx count: ", latestBlock.TransactionCount)
    fmt.Println("Latest block submitter: ", latestBlock.BlockSubmitter)
    
	pwrgo.ReturnBlockNumberOnTx = true // automatically calls blocksCount from RPC and returns BlockNumber on tx response

	// Transfer PWR
    //var transferTx = pwrgo.TransferPWR("0x61bd8fc1e30526aaf1c4706ada595d6d236d9883", "1", nonce, wallet.PrivateKey) // send 1 PWR
	//if transferTx.Success {
	//	fmt.Printf("[Block #%d] Transfer tx hash: %s\n", transferTx.BlockNumber, transferTx.TxHash)
	//	nonce = nonce + 1 // increment nonce since we just Transferred PWR
	//} else {
	//	fmt.Println("Error sending Transfer tx: ", transferTx.Error)
	//	fmt.Println("Error sending ", transferTx.TxHash)
	//}

    // Create new wallet and print address and keys
    var newWallet = pwrgo.NewWallet()
    fmt.Println("New wallet address: ", newWallet.Address)
    fmt.Println("New wallet private key: ", newWallet.PrivateKeyStr)
    fmt.Println("New wallet public key: ", newWallet.PublicKey)
    
	// // Claim VM ID
	// var vmId = int64(1337)
	// var claimTxResponse = pwrgo.ClaimVMId(vmId, nonce, wallet.PrivateKey)
	// if claimTxResponse.Success {
	//    fmt.Printf("[Block #%d] Claimed VM ID tx hash: %s", claimTxResponse.BlockNumber, claimTxResponse.TxHash)
	//    nonce = nonce + 1
	// } else {
	//    fmt.Println("Error claiming VM ID tx: ", claimTxResponse.Error)
	//    fmt.Println("Error sending ", claimTxResponse.TxHash)
	// }

    // // Send data to VM 1337
    // var data = []byte("Hello world")
    // var vmTxResponse = pwrgo.SendVMDataTx(1337, data, nonce, wallet.PrivateKey)
	// if vmTxResponse.Success {
	// 	fmt.Printf("[Block #%d] VM data tx hash: %s", vmTxResponse.BlockNumber, vmTxResponse.TxHash)
	// } else {
	// 	fmt.Println("Error sending VM data tx: ", vmTxResponse.Error)
	// 	fmt.Println("Error sending ", vmTxResponse.TxHash)
	// }


	// Get Transfer TX Bytes
	// var transferBytes = pwrgo.GetTransferTxBytes("0x61bd8fc1e30526aaf1c4706ada595d6d236d9883", "1", nonce)
	
	// Send Conduit TX
	// var conduitTxResponse = pwrgo.SendConduitTx(1337, transferBytes, nonce, wallet.PrivateKey)
	// if conduitTxResponse.Success {
	// 	fmt.Printf("[Block #%d] Conduit tx hash: %s", conduitTxResponse.BlockNumber, conduitTxResponse.TxHash)
	// } else {
	// 	fmt.Println("Error sending conduit tx: ", conduitTxResponse.Error)
	// 	fmt.Println("Error sending ", conduitTxResponse.TxHash)
	// }

	// Delegate PWR
	// var delegateTxResponse = pwrgo.Delegate("0x61bd8fc1e30526aaf1c4706ada595d6d236d9883", "899886346", nonce, wallet.PrivateKey)
	// if delegateTxResponse.Success {
	// 	fmt.Printf("[Block #%d] Delegate tx hash: %s", delegateTxResponse.BlockNumber, delegateTxResponse.TxHash)
	// } else {
	// 	fmt.Println("Error sending delegate tx: ", delegateTxResponse.Error)
	// 	fmt.Println("Error sending ", delegateTxResponse.TxHash)
	// }

	// Withdraw PWR
	// var withdrawTxResponse = pwrgo.Withdraw("0x61bd8fc1e30526aaf1c4706ada595d6d236d9883", "899886346", nonce, wallet.PrivateKey)
	// if withdrawTxResponse.Success {
	// 	fmt.Printf("[Block #%d] Withdraw tx hash: %s", withdrawTxResponse.BlockNumber, withdrawTxResponse.TxHash)
	// } else {
	// 	fmt.Println("Error sending withdraw tx: ", withdrawTxResponse.Error)
	// 	fmt.Println("Error sending ", withdrawTxResponse.TxHash)
	// }
	
	// Set Guardian
	// var setGuardianTxResponse = pwrgo.SetGuardian("0x61bd8fc1e30526aaf1c4706ada595d6d236d9883", "1710969140", nonce, wallet.PrivateKey)
	// if setGuardianTxResponse.Success {
	// 	fmt.Printf("[Block #%d] SetGuardian tx hash: %s", setGuardianTxResponse.BlockNumber, setGuardianTxResponse.TxHash)
	// } else {
	// 	fmt.Println("Error sending SetGuardian tx: ", setGuardianTxResponse.Error)
	// 	fmt.Println("Error sending ", setGuardianTxResponse.TxHash)
	// }

	// Remove Guardian
	// var removeGuardianTxResponse = pwrgo.RemoveGuardian(nonce, wallet.PrivateKey)
	// if removeGuardianTxResponse.Success {
	// 	fmt.Printf("[Block #%d] RemoveGuardian tx hash: %s", removeGuardianTxResponse.BlockNumber, removeGuardianTxResponse.TxHash)
	// } else {
	// 	fmt.Println("Error sending RemoveGuardian tx: ", removeGuardianTxResponse.Error)
	// 	fmt.Println("Error sending ", removeGuardianTxResponse.TxHash)
	// }

}