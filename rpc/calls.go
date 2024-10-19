package rpc

import (
	"log"
	"strconv"
	"strings"
)

func GetNonceOfAddress(address string) int {
	var response = get(GetRpcNodeUrl() + "/nonceOfUser/?userAddress=" + address)
	var resp = parseRPCResponse(response)
	return resp.Nonce
}

func GetBalanceOfAddress(address string) int {
	var response = get(GetRpcNodeUrl() + "/balanceOf/?userAddress=" + address)
	var resp = parseRPCResponse(response)
	return resp.Balance
}

func GetBlocksCount() int {
	var response = get(GetRpcNodeUrl() + "/blocksCount/")
	var resp = parseRPCResponse(response)
	return resp.BlocksCount
}

func GetValidatorsCount() int {
	var response = get(GetRpcNodeUrl() + "/totalValidatorsCount/")
	var resp = parseRPCResponse(response)
	return resp.ValidatorsCount
}

func GetBlockByNumber(blockNumber int) Block {
	var blockNumberStr = strconv.Itoa(blockNumber)
	var response = get(GetRpcNodeUrl() + "/block/?blockNumber=" + blockNumberStr)
	var resp = parseRPCResponse(response)
	return resp.Block
}

func GetBlockchainVersion() int64 {
	var response = get(GetRpcNodeUrl() + "/blockchainVersion/")
	num, _ := strconv.ParseInt(response, 10, 64)
	return num
}

// func GetFee() int64 {

// }

func GetEcdsaVerificationFee() int {
	var response = get(GetRpcNodeUrl() + "/ecdsaVerificationFee/")
	var resp = parseRPCResponse(response)
	return resp.ECDSAVerificationFee
}

func GetBurnPercentage() int {
	var response = get(GetRpcNodeUrl() + "/burnPercentage/")
	var resp = parseRPCResponse(response)
	return resp.BurnPercentage
}

func GetTotalVotingPower() int {
	var response = get(GetRpcNodeUrl() + "/totalVotingPower/")
	var resp = parseRPCResponse(response)
	return resp.TotalVotingPower
}

func GetPwrRewardsPerYear() int {
	var response = get(GetRpcNodeUrl() + "/pwrRewardsPerYear/")
	var resp = parseRPCResponse(response)
	return resp.PwrRewardsPerYear
}

func GetWithdrawalLockTime() int {
	var response = get(GetRpcNodeUrl() + "/withdrawalLockTime/")
	var resp = parseRPCResponse(response)
	return resp.WithdrawalLockTime
}

func GetAllEarlyWithdrawPenalties() []Penalty {
	var response = get(GetRpcNodeUrl() + "/allEarlyWithdrawPenalties/")
	var resp = parseRPCResponse(response)

	var penalties []Penalty
	for key, value := range resp.AllEarlyWithdrawPenalties {
		withdrawTime, err := strconv.ParseInt(key, 10, 64)
		if err != nil {
			log.Printf("Skipping invalid key: %v", key)
			continue
		}
		penalties = append(penalties, Penalty{WithdrawTime: withdrawTime, Penalty: value})
	}

	return penalties
}

func GetMaxBlockSize() int {
	var response = get(GetRpcNodeUrl() + "/maxBlockSize/")
	var resp = parseRPCResponse(response)
	return resp.MaxBlockSize
}

func GetMaxTransactionSize() int {
	var response = get(GetRpcNodeUrl() + "/maxTransactionSize/")
	var resp = parseRPCResponse(response)
	return resp.MaxTransactionSize
}

func getBlockNumber() int {
	var response = get(GetRpcNodeUrl() + "/blockNumber/")
	var resp = parseRPCResponse(response)
	return resp.BlockNumber
}

func GetBlockTimestamp() int {
	var response = get(GetRpcNodeUrl() + "/blockTimestamp/")
	var resp = parseRPCResponse(response)
	return resp.BlockTimestamp
}

func GetLatestBlockNumber() int {
	return getBlockNumber()
}

func GetProposalFee() int {
	var response = get(GetRpcNodeUrl() + "/proposalFee/")
	var resp = parseRPCResponse(response)
	return resp.ProposalFee
}

func GetProposalValidityTime() int {
	var response = get(GetRpcNodeUrl() + "/proposalValidityTime/")
	var resp = parseRPCResponse(response)
	return resp.ProposalValidityTime
}

func GetValidatorCountLimit() int {
	var response = get(GetRpcNodeUrl() + "/validatorCountLimit/")
	var resp = parseRPCResponse(response)
	return resp.ValidatorCountLimit
}

func GetValidatorSlashingFee() int {
	var response = get(GetRpcNodeUrl() + "/validatorSlashingFee/")
	var resp = parseRPCResponse(response)
	return resp.ValidatorSlashingFee
}

func GetValidatorOperationalFee() int {
	var response = get(GetRpcNodeUrl() + "/validatorOperationalFee/")
	var resp = parseRPCResponse(response)
	return resp.ValidatorOperationalFee
}

func GetValidatorJoiningFee() int {
	var response = get(GetRpcNodeUrl() + "/validatorJoiningFee/")
	var resp = parseRPCResponse(response)
	return resp.ValidatorJoiningFee
}

func GetMinimumDelegatingAmount() int {
	var response = get(GetRpcNodeUrl() + "/minimumDelegatingAmount/")
	var resp = parseRPCResponse(response)
	return resp.MinimumDelegatingAmount
}

func GetDelegatorsCount() int {
	var response = get(GetRpcNodeUrl() + "/totalDelegatorsCount/")
	var resp = parseRPCResponse(response)
	return resp.DelegatorsCount
}

func GetValidator(address string) Validator {
	var response = get(GetRpcNodeUrl() + "/validator/?validatorAddress=" + address)
	var resp = parseRPCResponse(response)
	return resp.Validator
}

func GetDelegatorsOfPwr(delegatorAddress string, validatorAddress string) int {
	var response = get(
		GetRpcNodeUrl() + "/validator/delegator/delegatedPWROfAddress/?userAddress=" +
			delegatorAddress + "&validatorAddress=" + validatorAddress,
	)
	var resp = parseRPCResponse(response)
	return resp.DelegatedPWR
}

func GetSharesOfDelegator(delegatorAddress string, validatorAddress string) int {
	var response = get(
		GetRpcNodeUrl() + "/validator/delegator/sharesOfAddress/?userAddress=" +
			delegatorAddress + "&validatorAddress=" + validatorAddress,
	)
	var resp = parseRPCResponse(response)
	return resp.SharesOfDelegator
}

func GetShareValue(validatorAddress string) float32 {
	var response = get(
		GetRpcNodeUrl() + "/validator/shareValue/?validatorAddress=" + validatorAddress,
	)
	var resp = parseRPCResponse(response)
	return resp.ShareValue
}

func GetVmOwnerTransactionFeeShare() int {
	var response = get(GetRpcNodeUrl() + "/vmOwnerTransactionFeeShare/")
	var resp = parseRPCResponse(response)
	return resp.VmOwnerTransactionFeeShare
}

func GetVmIdClaimingFee() int {
	var response = get(GetRpcNodeUrl() + "/vmIdClaimingFee/")
	var resp = parseRPCResponse(response)
	return resp.VmIdClaimingFee
}

func GetVmDataTransactions(startingBlock int, endingBlock int, vmId int) []VMDataTransaction {
	startingBlockStr := strconv.Itoa(startingBlock)
	endingBlockStr := strconv.Itoa(endingBlock)
	vmIdStr := strconv.Itoa(vmId)
	var response = get(
		GetRpcNodeUrl() + "/getVmTransactions/?startingBlock=" + startingBlockStr +
			"&endingBlock=" + endingBlockStr + "&vmId=" + vmIdStr,
	)
	var resp = parseRPCResponse(response)
	return resp.VMDataTransaction
}

func GetVmIdAddress(vmId int) string {
	hexAddress := "0"
	if vmId >= 0 {
		hexAddress = "1"
	}
	if vmId < 0 {
		vmId = -vmId
	}

	vmIdString := strconv.Itoa(vmId)
	padding := 39 - len(vmIdString)
	if padding > 0 {
		hexAddress += strings.Repeat("0", padding)
	}
	hexAddress += vmIdString

	return "0x" + hexAddress
}

func GetOwnerOfVmIds(vmId int) string {
	vmIdStr := strconv.Itoa(vmId)
	var response = get(GetRpcNodeUrl() + "/ownerOfVmId/?vmId=" + vmIdStr)
	var resp = parseRPCResponse(response)
	return "0x" + resp.OwnerOfVmIds
}

func GetConduitsOfVmId(vmId int) []Validator {
	vmIdStr := strconv.Itoa(vmId)
	var response = get(GetRpcNodeUrl() + "/conduitsOfVm/?vmId=" + vmIdStr)
	var resp = parseRPCResponse(response)
	return resp.ConduitsOfVm
}

func GetMaxGuardianTime() int {
	var response = get(GetRpcNodeUrl() + "/maxGuardianTime/")
	var resp = parseRPCResponse(response)
	return resp.MaxGuardianTime
}

func GetGuardianOfAddress(address string) string {
	var response = get(GetRpcNodeUrl() + "/guardianOf/?userAddress=" + address)
	var resp = parseRPCResponse(response)
	return resp.GuardianOfAddress
}

func GetActiveVotingPower() int {
	var response = get(GetRpcNodeUrl() + "/activeVotingPower/")
	var resp = parseRPCResponse(response)
	return resp.ActiveVotingPower
}

func GetStandbyValidatorsCount() int {
	var response = get(GetRpcNodeUrl() + "/standbyValidatorsCount/")
	var resp = parseRPCResponse(response)
	return resp.ValidatorsCount
}

func GetActiveValidatorsCount() int {
	var response = get(GetRpcNodeUrl() + "/activeValidatorsCount/")
	var resp = parseRPCResponse(response)
	return resp.ValidatorsCount
}

func GetValidators() []Validator {
	var response = get(GetRpcNodeUrl() + "/allValidators/")
	var resp = parseRPCResponse(response)
	return resp.Validators
}

func GetStandbyValidators() []Validator {
	var response = get(GetRpcNodeUrl() + "/standbyValidators/")
	var resp = parseRPCResponse(response)
	return resp.Validators
}

func GetActiveValidators() []Validator {
	var response = get(GetRpcNodeUrl() + "/activeValidators/")
	var resp = parseRPCResponse(response)
	return resp.Validators
}

func GetTransactionByHash(txHash string) Transaction {
	var response = get(GetRpcNodeUrl() + "/transactionByHash/?transactionHash=" + txHash)
	var resp = parseRPCResponse(response)

	return resp.Transaction
}

func GetDelegatorsOfValidator(validatorAddress string) []Delegator {
	var response = get(
		GetRpcNodeUrl() + "/validator/delegatorsOfValidator/?validatorAddress=" + validatorAddress,
	)
	var resp = parseRPCResponse(response)

	delegators := []Delegator{}
	if resp.Delegators != nil {
		for address, shares := range resp.Delegators {
			delegators = append(delegators, Delegator{
				Address: address,
				Shares:  shares,
			})
		}
	}

	return delegators
}

func GetFeeBerByte() int {
	var response = get(GetRpcNodeUrl() + "/feePerByte/")
	var resp = parseRPCResponse(response)
	return resp.FeePerByte
}
