package main

import (
    "fmt"
    "encoding/json"
    "encoding/hex"

    "github.com/pwrlabs/pwrgo/rpc"
)

// Handler for incoming messages
func messageHandler(transaction rpc.VMDataTransaction) {
    sender := transaction.Sender
    data := transaction.Data
    
    dataBytes, _ := hex.DecodeString(data[2:])
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
    startingBlock := rpc.GetLatestBlockNumber()

    subscription := rpc.SubscribeToVidaTransactions(
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
    var delegators = rpc.GetDelegatorsOfValidator("0x7b6F32435084Cab827f0ce7Af1C0D48600CE3CaD")
    fmt.Println("Delegators:", delegators)
    
    var blocksCount = rpc.GetBlocksCount()
    fmt.Println("Blocks count:", blocksCount)

    var latestBlockCount = rpc.GetLatestBlockNumber()
    fmt.Println("Validators count:", latestBlockCount)

    var startBlcok = 65208
    var endBlock = 65210
    var vmId = 1234
    var transactions = rpc.GetVmDataTransactions(startBlcok, endBlock, vmId)
    fmt.Println("VM Data:", transactions)

    var guardian = rpc.GetGuardianOfAddress("0xD97C25C0842704588DD70A061C09A522699E2B9C")
    fmt.Println("Guardian:", guardian)

    var block = rpc.GetBlockByNumber(10)
    fmt.Println("Block:", block)

    var activeVotingPower = rpc.GetActiveVotingPower()
    fmt.Println("ActiveVotingPower:", activeVotingPower)

    var allValidators = rpc.GetAllValidators()
    fmt.Println("AllValidators:", allValidators)

    var conduitsVm = rpc.GetConduitsOfVmId(101001)
    fmt.Println("ConduitsVm:", conduitsVm)

    var totalValidatorsCount = rpc.GetValidatorsCount()
    fmt.Println("TotalValidatorsCount:", totalValidatorsCount)

    var tx = rpc.GetTransactionByHash("0x82c856bce3fb7ce2a504e8d108ed0ee59e5f8c5fc2c0002e94f9ef774da01911")
    fmt.Println("Transfer TX: ", tx)
}

func main() {
    RpcCall()
    // Vidas()
}
