package rpc

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
)

func (r *RPC) GetNonceOfAddress(address string) int {
	var response = get(r.GetRpcNodeUrl() + "/nonceOfUser?userAddress=" + address)
	var resp = parseRPCResponse(response)
	return resp.Nonce
}

func (r *RPC) GetBalanceOfAddress(address string) int {
	var response = get(r.GetRpcNodeUrl() + "/balanceOf?userAddress=" + address)
	var resp = parseRPCResponse(response)
	return resp.Balance
}

func (r *RPC) GetBlocksCount() int {
	var response = get(r.GetRpcNodeUrl() + "/blocksCount")
	var resp = parseRPCResponse(response)
	return resp.BlocksCount
}

func (r *RPC) GetValidatorsCount() int {
	var response = get(r.GetRpcNodeUrl() + "/totalValidatorsCount")
	var resp = parseRPCResponse(response)
	return resp.ValidatorsCount
}

func (r *RPC) GetBlockByNumber(blockNumber int) Block {
	var blockNumberStr = strconv.Itoa(blockNumber)
	var response = get(r.GetRpcNodeUrl() + "/block?blockNumber=" + blockNumberStr)
	var resp = parseRPCResponse(response)
	return resp.Block
}

func (r *RPC) GetBlockchainVersion() int64 {
	var response = get(r.GetRpcNodeUrl() + "/blockchainVersion")
	num, _ := strconv.ParseInt(response, 10, 64)
	return num
}

// func GetFee() int64 {

// }

func (r *RPC) GetEcdsaVerificationFee() int {
	var response = get(r.GetRpcNodeUrl() + "/ecdsaVerificationFee")
	var resp = parseRPCResponse(response)
	return resp.ECDSAVerificationFee
}

func (r *RPC) GetBurnPercentage() int {
	var response = get(r.GetRpcNodeUrl() + "/burnPercentage")
	var resp = parseRPCResponse(response)
	return resp.BurnPercentage
}

func (r *RPC) GetTotalVotingPower() int {
	var response = get(r.GetRpcNodeUrl() + "/totalVotingPower")
	var resp = parseRPCResponse(response)
	return resp.TotalVotingPower
}

func (r *RPC) GetPwrRewardsPerYear() int {
	var response = get(r.GetRpcNodeUrl() + "/pwrRewardsPerYear")
	var resp = parseRPCResponse(response)
	return resp.PwrRewardsPerYear
}

func (r *RPC) GetWithdrawalLockTime() int {
	var response = get(r.GetRpcNodeUrl() + "/withdrawalLockTime")
	var resp = parseRPCResponse(response)
	return resp.WithdrawalLockTime
}

func (r *RPC) GetAllEarlyWithdrawPenalties() []Penalty {
	var response = get(r.GetRpcNodeUrl() + "/allEarlyWithdrawPenalties")
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

func (r *RPC) GetMaxBlockSize() int {
	var response = get(r.GetRpcNodeUrl() + "/maxBlockSize")
	var resp = parseRPCResponse(response)
	return resp.MaxBlockSize
}

func (r *RPC) GetMaxTransactionSize() int {
	var response = get(r.GetRpcNodeUrl() + "/maxTransactionSize")
	var resp = parseRPCResponse(response)
	return resp.MaxTransactionSize
}

func (r *RPC) getBlockNumber() int {
	var response = get(r.GetRpcNodeUrl() + "/blockNumber")
	var resp = parseRPCResponse(response)
	return resp.BlockNumber
}

func (r *RPC) GetBlockTimestamp() int {
	var response = get(r.GetRpcNodeUrl() + "/blockTimestamp")
	var resp = parseRPCResponse(response)
	return resp.BlockTimestamp
}

func (r *RPC) GetLatestBlockNumber() int {
	return r.getBlockNumber()
}

func (r *RPC) GetProposalFee() int {
	var response = get(r.GetRpcNodeUrl() + "/proposalFee")
	var resp = parseRPCResponse(response)
	return resp.ProposalFee
}

func (r *RPC) GetProposalValidityTime() int {
	var response = get(r.GetRpcNodeUrl() + "/proposalValidityTime")
	var resp = parseRPCResponse(response)
	return resp.ProposalValidityTime
}

func (r *RPC) GetValidatorCountLimit() int {
	var response = get(r.GetRpcNodeUrl() + "/validatorCountLimit")
	var resp = parseRPCResponse(response)
	return resp.ValidatorCountLimit
}

func (r *RPC) GetValidatorSlashingFee() int {
	var response = get(r.GetRpcNodeUrl() + "/validatorSlashingFee")
	var resp = parseRPCResponse(response)
	return resp.ValidatorSlashingFee
}

func (r *RPC) GetValidatorOperationalFee() int {
	var response = get(r.GetRpcNodeUrl() + "/validatorOperationalFee")
	var resp = parseRPCResponse(response)
	return resp.ValidatorOperationalFee
}

func (r *RPC) GetValidatorJoiningFee() int {
	var response = get(r.GetRpcNodeUrl() + "/validatorJoiningFee")
	var resp = parseRPCResponse(response)
	return resp.ValidatorJoiningFee
}

func (r *RPC) GetMinimumDelegatingAmount() int {
	var response = get(r.GetRpcNodeUrl() + "/minimumDelegatingAmount")
	var resp = parseRPCResponse(response)
	return resp.MinimumDelegatingAmount
}

func (r *RPC) GetDelegatorsCount() int {
	var response = get(r.GetRpcNodeUrl() + "/totalDelegatorsCount")
	var resp = parseRPCResponse(response)
	return resp.DelegatorsCount
}

func (r *RPC) GetValidator(address string) Validator {
	var response = get(r.GetRpcNodeUrl() + "/validator?validatorAddress=" + address)
	var resp = parseRPCResponse(response)
	return resp.Validator
}

func (r *RPC) GetDelegatorsOfPwr(delegatorAddress string, validatorAddress string) int {
	var response = get(
		r.GetRpcNodeUrl() + "/validator/delegator/delegatedPWROfAddress?userAddress=" +
			delegatorAddress + "&validatorAddress=" + validatorAddress,
	)
	var resp = parseRPCResponse(response)
	return resp.DelegatedPWR
}

func (r *RPC) GetSharesOfDelegator(delegatorAddress string, validatorAddress string) int {
	var response = get(
		r.GetRpcNodeUrl() + "/validator/delegator/sharesOfAddress?userAddress=" +
			delegatorAddress + "&validatorAddress=" + validatorAddress,
	)
	var resp = parseRPCResponse(response)
	return resp.SharesOfDelegator
}

func (r *RPC) GetShareValue(validatorAddress string) float32 {
	var response = get(
		r.GetRpcNodeUrl() + "/validator/shareValue?validatorAddress=" + validatorAddress,
	)
	var resp = parseRPCResponse(response)
	return resp.ShareValue
}

func (r *RPC) GetVmOwnerTransactionFeeShare() int {
	var response = get(r.GetRpcNodeUrl() + "/vmOwnerTransactionFeeShare")
	var resp = parseRPCResponse(response)
	return resp.VmOwnerTransactionFeeShare
}

func (r *RPC) GetVidaIdClaimingFee() int {
	var response = get(r.GetRpcNodeUrl() + "/vidaIdClaimingFee")
	var resp = parseRPCResponse(response)
	return resp.VmIdClaimingFee
}

func (r *RPC) GetVidaDataTransactions(startingBlock int, endingBlock int, vidaId int) []VMDataTransaction {
	startingBlockStr := strconv.Itoa(startingBlock)
	endingBlockStr := strconv.Itoa(endingBlock)
	vidaIdStr := strconv.Itoa(vidaId)
	var response = get(
		r.GetRpcNodeUrl() + "/getVidaTransactions?startingBlock=" + startingBlockStr +
			"&endingBlock=" + endingBlockStr + "&vidaId=" + vidaIdStr,
	)

	var resp struct {
		Transactions []string `json:"transactions"`
	}

	if err := json.Unmarshal([]byte(response), &resp); err != nil {
		log.Fatal("Error unmarshaling response: ", err)
	}

	var transactions []VMDataTransaction
	for _, txStr := range resp.Transactions {
		var tx VMDataTransaction
		if err := json.Unmarshal([]byte(txStr), &tx); err != nil {
			log.Fatal("Error unmarshaling transaction: ", err)
		}
		transactions = append(transactions, tx)
	}

	return transactions
}

func (r *RPC) GetVidaIdAddress(vidaId int) string {
	hexAddress := "0"
	if vidaId >= 0 {
		hexAddress = "1"
	}
	if vidaId < 0 {
		vidaId = -vidaId
	}

	vmIdString := strconv.Itoa(vidaId)
	padding := 39 - len(vmIdString)
	if padding > 0 {
		hexAddress += strings.Repeat("0", padding)
	}
	hexAddress += vmIdString

	return "0x" + hexAddress
}

func (r *RPC) GetOwnerOfVidaIds(vidaId int) string {
	vidaIdStr := strconv.Itoa(vidaId)
	var response = get(r.GetRpcNodeUrl() + "/ownerOfVidaId?vidaId=" + vidaIdStr)
	var resp = parseRPCResponse(response)
	return "0x" + resp.OwnerOfVidaIds
}

func (r *RPC) GetConduitsOfVmId(vidaId int) []Validator {
	vidaIdStr := strconv.Itoa(vidaId)
	var response = get(r.GetRpcNodeUrl() + "/conduitsOfVida?vidaId=" + vidaIdStr)
	var resp = parseRPCResponse(response)
	return resp.ConduitsOfVm
}

func (r *RPC) GetMaxGuardianTime() int {
	var response = get(r.GetRpcNodeUrl() + "/maxGuardianTime")
	var resp = parseRPCResponse(response)
	return resp.MaxGuardianTime
}

func (r *RPC) GetGuardianOfAddress(address string) string {
	var response = get(r.GetRpcNodeUrl() + "/guardianOf?userAddress=" + address)
	var resp = parseRPCResponse(response)
	return resp.GuardianOfAddress
}

func (r *RPC) GetActiveVotingPower() int {
	var response = get(r.GetRpcNodeUrl() + "/activeVotingPower")
	var resp = parseRPCResponse(response)
	return resp.ActiveVotingPower
}

func (r *RPC) GetStandbyValidatorsCount() int {
	var response = get(r.GetRpcNodeUrl() + "/standbyValidatorsCount")
	var resp = parseRPCResponse(response)
	return resp.ValidatorsCount
}

func (r *RPC) GetActiveValidatorsCount() int {
	var response = get(r.GetRpcNodeUrl() + "/activeValidatorsCount")
	var resp = parseRPCResponse(response)
	return resp.ValidatorsCount
}

func (r *RPC) GetAllValidators() []Validator {
	var response = get(r.GetRpcNodeUrl() + "/allValidators")
	var resp = parseRPCResponse(response)
	return resp.Validators
}

func (r *RPC) GetStandbyValidators() []Validator {
	var response = get(r.GetRpcNodeUrl() + "/standbyValidators")
	var resp = parseRPCResponse(response)
	return resp.Validators
}

func (r *RPC) GetActiveValidators() []Validator {
	var response = get(r.GetRpcNodeUrl() + "/activeValidators")
	var resp = parseRPCResponse(response)
	return resp.Validators
}

func (r *RPC) GetTransactionByHash(txHash string) Transaction {
	var response = get(r.GetRpcNodeUrl() + "/transactionByHash?transactionHash=" + txHash)
	var resp = parseRPCResponse(response)
	return resp.Transaction
}

func (r *RPC) GetDelegatorsOfValidator(validatorAddress string) []Delegator {
	var response = get(
		r.GetRpcNodeUrl() + "/validator/delegatorsOfValidator?validatorAddress=" + validatorAddress,
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

func (r *RPC) GetFeeBerByte() int {
	var response = get(r.GetRpcNodeUrl() + "/feePerByte")
	var resp = parseRPCResponse(response)
	return resp.FeePerByte
}
