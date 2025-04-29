package main

import (
    "fmt"
    "encoding/json"
    "encoding/hex"

    "github.com/pwrlabs/pwrgo/rpc"
)

var pwr = rpc.SetRpcNodeUrl("https://pwrrpc.pwrlabs.io")

// Handler for incoming messages
func messageHandler(transaction rpc.VMDataTransaction) {
    sender := transaction.Sender
    data := transaction.Data
    
    dataBytes, _ := hex.DecodeString(data)
    var obj map[string]interface{}

    if err := json.Unmarshal(dataBytes, &obj); err != nil {
        fmt.Println("Error parsing JSON:", err)
    }

    if action, _ := obj["action"].(string); action == "send-message-v1" {
        message, _ := obj["message"].(string)
        fmt.Printf("Message from %s: %s\n", sender, message)
    }
}

func Vidas() {
    vidaId := 1
    startingBlock := pwr.GetLatestBlockNumber()

    subscription := pwr.SubscribeToVidaTransactions(
        vidaId,
        startingBlock,
        messageHandler,
    )

    subscription.Pause()
    subscription.Resume()
    // subscription.Stop()

    fmt.Println("Latest checked blocked:", subscription.GetLatestCheckedBlock())

    if (subscription.IsRunning()) {
        fmt.Println("Press Enter to exit...")
        fmt.Scanln()
    }
}

func RpcCall() {
    var rpcUrl = pwr.GetRpcNodeUrl()
    fmt.Println("RPC URL:", rpcUrl)

    // var delegators = pwr.GetDelegatorsOfValidator("0x0EA54532D7CA460083E547910D5C30C5896967C9")
    // fmt.Println("Delegators:", delegators)
    
    var blocksCount = pwr.GetBlocksCount()
    fmt.Println("Blocks count:", blocksCount)

    var latestBlockCount = pwr.GetLatestBlockNumber()
    fmt.Println("Validators count:", latestBlockCount)

    var block = pwr.GetBlockByNumber(10)
    fmt.Println("Block:", block)

    var balance = pwr.GetBalanceOfAddress("0xf8d42a75bdb93769c44e2b6c53c82fa77804cdd4")
    fmt.Println("Balance:", balance)

    var startBlcok = 1176
    var endBlock = 1179
    var vmId = 1234
    var transactions = pwr.GetVidaDataTransactions(startBlcok, endBlock, vmId)
    fmt.Println("VM Data:", transactions)

    // var guardian = pwr.GetGuardianOfAddress("0xD97C25C0842704588DD70A061C09A522699E2B9C")
    // fmt.Println("Guardian:", guardian)

    var activeVotingPower = pwr.GetActiveVotingPower()
    fmt.Println("ActiveVotingPower:", activeVotingPower)

    var allValidators = pwr.GetAllValidators()
    fmt.Println("AllValidators:", allValidators)

    // var conduitsVm = pwr.GetConduitsOfVmId(101001)
    // fmt.Println("ConduitsVm:", conduitsVm)

    var totalValidatorsCount = pwr.GetValidatorsCount()
    fmt.Println("TotalValidatorsCount:", totalValidatorsCount)

    var tx = pwr.GetTransactionByHash("0x0075C4D04BC18586CF6F6ECE8783E922F4AAC8D58C71D429848AD637F2DAC33A")
    fmt.Println("Transfer TX: ", tx)
}

func main() {
    RpcCall()
    // Vidas()
}
