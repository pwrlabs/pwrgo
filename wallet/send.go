package wallet

import (
	"log"

	"github.com/pwrlabs/pwrgo/encode"
	"github.com/pwrlabs/pwrgo/rpc"
)

func (w *PWRWallet) SetPublicKey(publicKey []byte, feePerByte int) rpc.BroadcastResponse {
	var buffer []byte
	buffer, err := encode.SetPublicKeyBytes(publicKey, w.GetNonce(), w.Address, feePerByte)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *PWRWallet) JoinAsValidator(ip string, feePerByte int) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := encode.JoinAsValidatorBytes(ip, w.GetNonce(), w.Address, feePerByte)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *PWRWallet) Delegate(to string, amount int, feePerByte int) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := encode.DelegateTxBytes(to, amount, w.GetNonce(), w.Address, feePerByte)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *PWRWallet) ChangeIp(newIp string, feePerByte int) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := encode.ChangeIpBytes(newIp, w.GetNonce(), w.Address, feePerByte)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *PWRWallet) ClaimActiveNodeSpot(feePerByte int) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := encode.ClaimActiveNodeSpotBytes(w.GetNonce(), w.Address, feePerByte)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *PWRWallet) TransferPWR(to string, amount int, feePerByte int) rpc.BroadcastResponse {
	if len(to) != 42 {
		log.Fatal("Invalid address: ", to)
	}

	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := encode.TransferTxBytes(amount, to, w.GetNonce(), w.Address, feePerByte)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *PWRWallet) ProposeChangeEarlyWithdrawPenalty(
	title string, description string, withdraw_penalty_time int,
	withdraw_penalty int, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := encode.EarlyWithdrawPenaltyProposal(
		title, description, withdraw_penalty_time, withdraw_penalty, w.GetNonce(), w.Address, feePerByte,
	)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *PWRWallet) ProposeChangeFeePerByte(
	title string, description string, feePerByte int, _feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(_feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := encode.ChangeFeePerByteProposalTx(
		title, description, feePerByte, w.GetNonce(), w.Address, _feePerByte,
	)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *PWRWallet) ProposeChangeMaxBlockSize(
	title string, description string, maxBlockSize int, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := encode.ChangeMaxBlockSizeProposalTx(
		title, description, maxBlockSize, w.GetNonce(), w.Address, feePerByte,
	)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *PWRWallet) ProposeChangeMaxTxnSize(
	title string, description string, maxTxnSize int, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := encode.ChangeMaxTxnSizeProposalTx(
		title, description, maxTxnSize, w.GetNonce(), w.Address, feePerByte,
	)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *PWRWallet) ProposeChangeOverallBurnPercentage(
	title string, description string, burnPercentage int, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := encode.ChangeOverallBurnPercentageProposalTx(
		title, description, burnPercentage, w.GetNonce(), w.Address, feePerByte,
	)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *PWRWallet) ProposeChangeRewardPerYear(
	title string, description string, rewardPerYear int64, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := encode.ChangeRewardPerYearProposalTx(
		title, description, rewardPerYear, w.GetNonce(), w.Address, feePerByte,
	)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *PWRWallet) ProposeChangeValidatorCountLimit(
	title string, description string, validatorCountLimit int, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := encode.ChangeValidatorCountLimitProposalTx(
		title, description, validatorCountLimit, w.GetNonce(), w.Address, feePerByte,
	)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *PWRWallet) ProposeChangeValidatorJoiningFee(
	title string, description string, joiningFee int64, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := encode.ChangeValidatorJoiningFeeProposalTx(
		title, description, joiningFee, w.GetNonce(), w.Address, feePerByte,
	)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *PWRWallet) ProposeChangeVidaIdClaimingFee(
	title string, description string, vidaIdClaimingFee int64, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := encode.ChangeVidaIdClaimingFeeProposalTx(
		title, description, vidaIdClaimingFee, w.GetNonce(), w.Address, feePerByte,
	)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *PWRWallet) ProposeChangeVmOwnerTxnFeeShare(
	title string, description string, vmOwnerTxnFeeShare int, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := encode.ChangeVmOwnerTxnFeeShareProposalTx(
		title, description, vmOwnerTxnFeeShare, w.GetNonce(), w.Address, feePerByte,
	)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *PWRWallet) ProposeOther(
	title string, description string, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := encode.OtherProposalTx(
		title, description, w.GetNonce(), w.Address, feePerByte,
	)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *PWRWallet) VoteOnProposal(
	proposalHash []byte, vote byte, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := encode.VoteOnProposalTx(
		proposalHash, vote, w.GetNonce(), w.Address, feePerByte,
	)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *PWRWallet) GuardianApproval(
	wrappedTxns [][]byte, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := encode.GuardianApprovalTransaction(
		wrappedTxns, w.GetNonce(), w.Address, feePerByte,
	)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *PWRWallet) RemoveGuardian(feePerByte int) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := encode.RemoveGuardianTransaction(w.GetNonce(), w.Address, feePerByte)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *PWRWallet) SetGuardian(
	expiryDate int, guardianAddress string, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := encode.SetGuardianTransaction(
		expiryDate, guardianAddress, w.GetNonce(), w.Address, feePerByte,
	)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *PWRWallet) MoveStake(
	sharesAmount int64, fromValidator string, toValidator string, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := encode.MoveStakeTransaction(
		sharesAmount, fromValidator, toValidator, w.GetNonce(), w.Address, feePerByte,
	)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *PWRWallet) RemoveValidator(
	validatorAddress string, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := encode.RemoveValidatorTransaction(
		validatorAddress, w.GetNonce(), w.Address, feePerByte,
	)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *PWRWallet) Withdraw(
	sharesAmount int, validator string, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := encode.WithdrawTransaction(
		sharesAmount, validator, w.GetNonce(), w.Address, feePerByte,
	)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *PWRWallet) ClaimVidaId(
	vidaId int64, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := encode.ClaimVidaIdTransaction(
		vidaId, w.GetNonce(), w.Address, feePerByte,
	)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *PWRWallet) ConduitApproval(
	vidaId int64, wrappedTxns [][]byte, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := encode.ConduitApprovalTransaction(
		vidaId, wrappedTxns, w.GetNonce(), w.Address, feePerByte,
	)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *PWRWallet) PayableVidaData(
	vidaId int64, data []byte, value int64, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := encode.PayableVidaDataTransaction(
		vidaId, data, value, w.GetNonce(), w.Address, feePerByte,
	)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *PWRWallet) RemoveConduits(
	vidaId int64, conduits [][]byte, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := encode.RemoveConduitsTransaction(
		vidaId, conduits, w.GetNonce(), w.Address, feePerByte,
	)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *PWRWallet) SetConduitMode(
	vidaId int64, mode byte, conduitThreshold int, conduits [][]byte,
	conduitsWithVotingPower map[string]int64, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := encode.SetConduitModeTransaction(
		vidaId, mode, conduitThreshold, conduits, conduitsWithVotingPower,
		w.GetNonce(), w.Address, feePerByte,
	)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *PWRWallet) SetVidaPrivateState(
	vidaId int64, privateState bool, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := encode.SetVidaPrivateStateTransaction(
		vidaId, privateState, w.GetNonce(), w.Address, feePerByte,
	)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *PWRWallet) SetVidaToAbsolutePublic(
	vidaId int64, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := encode.SetVidaToAbsolutePublicTransaction(
		vidaId, w.GetNonce(), w.Address, feePerByte,
	)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *PWRWallet) AddVidaSponsoredAddresses(
	vidaId int64, sponsoredAddresses []string, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := encode.AddVidaSponsoredAddressesTransaction(
		vidaId, sponsoredAddresses, w.GetNonce(), w.Address, feePerByte,
	)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *PWRWallet) AddVidaAllowedSenders(
	vidaId int64, allowedSenders []string, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := encode.AddVidaAllowedSendersTransaction(
		vidaId, allowedSenders, w.GetNonce(), w.Address, feePerByte,
	)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *PWRWallet) RemoveVidaAllowedSenders(
	vidaId int64, allowedSenders []string, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := encode.RemoveVidaAllowedSendersTransaction(
		vidaId, allowedSenders, w.GetNonce(), w.Address, feePerByte,
	)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *PWRWallet) RemoveSponsoredAddresses(
	vidaId int64, sponsoredAddresses []string, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := encode.RemoveSponsoredAddressesTransaction(
		vidaId, sponsoredAddresses, w.GetNonce(), w.Address, feePerByte,
	)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *PWRWallet) SetPWRTransferRights(
	vidaId int64, ownerCanTransferPWR bool, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := encode.SetPWRTransferRightsTransaction(
		vidaId, ownerCanTransferPWR, w.GetNonce(), w.Address, feePerByte,
	)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *PWRWallet) TransferPWRFromVida(
	vidaId int64, receiver string, amount int64, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := encode.TransferPWRFromVidaTransaction(
		vidaId, receiver, amount, w.GetNonce(), w.Address, feePerByte,
	)
	if err != nil {
		log.Fatal("Failed to get tx bytes: ", err.Error())
	}

	txn_bytes, err := w.SignTx(buffer)
	if err != nil {
		log.Fatal("Failed to sign message: ", err.Error())
	}

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func makeSurePublicKeyIsSet(feePerByte int, w *PWRWallet) *rpc.BroadcastResponse {
	nonce := w.GetNonce()

	if nonce == 0 {
		tx := w.SetPublicKey(w.PublicKey, feePerByte)
		return &tx
	}

	return nil
}
