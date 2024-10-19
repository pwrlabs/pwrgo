package main

import (
    "fmt"
    "github.com/pwrlabs/pwrgo/rpc"
)

func main() {
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
	
	var tx = rpc.GetTransactionByHash("0xe8be5c174a3457ba1015bef3399bb71a4081b8c24bdbebb05a8aea746ac32486")
	fmt.Println("Transfer TX: ", tx)
}