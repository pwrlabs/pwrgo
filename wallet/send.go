package wallet

import (
	"log"

	"github.com/pwrlabs/pwrgo/config/transactions"
	"github.com/pwrlabs/pwrgo/rpc"
)

func (w *PWRWallet) SetPublicKey(publicKey []byte, feePerByte int) rpc.BroadcastResponse {
	var buffer []byte
	buffer, err := transactions.SetPublicKeyBytes(publicKey, w.GetNonce(), w.Address, feePerByte)
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
	buffer, err := transactions.JoinAsValidatorBytes(ip, w.GetNonce(), w.Address, feePerByte)
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
	buffer, err := transactions.DelegateTxBytes(to, amount, w.GetNonce(), w.Address, feePerByte)
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
	buffer, err := transactions.ChangeIpBytes(newIp, w.GetNonce(), w.Address, feePerByte)
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
	buffer, err := transactions.ClaimActiveNodeSpotBytes(w.GetNonce(), w.Address, feePerByte)
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
	buffer, err := transactions.TransferTxBytes(amount, to, w.GetNonce(), w.Address, feePerByte)
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
	buffer, err := transactions.EarlyWithdrawPenaltyProposal(
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
	buffer, err := transactions.ChangeFeePerByteProposalTx(
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
	buffer, err := transactions.ChangeMaxBlockSizeProposalTx(
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
	buffer, err := transactions.ChangeMaxTxnSizeProposalTx(
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
	buffer, err := transactions.ChangeOverallBurnPercentageProposalTx(
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
	buffer, err := transactions.ChangeRewardPerYearProposalTx(
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
	buffer, err := transactions.ChangeValidatorCountLimitProposalTx(
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
	buffer, err := transactions.ChangeValidatorJoiningFeeProposalTx(
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
	buffer, err := transactions.ChangeVidaIdClaimingFeeProposalTx(
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
	buffer, err := transactions.ChangeVmOwnerTxnFeeShareProposalTx(
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
	buffer, err := transactions.OtherProposalTx(
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
	buffer, err := transactions.VoteOnProposalTx(
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
	buffer, err := transactions.GuardianApprovalTransaction(
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
	buffer, err := transactions.RemoveGuardianTransaction(w.GetNonce(), w.Address, feePerByte)
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
	buffer, err := transactions.SetGuardianTransaction(
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
	buffer, err := transactions.MoveStakeTransaction(
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
	buffer, err := transactions.RemoveValidatorTransaction(
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
	buffer, err := transactions.WithdrawTransaction(
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
	vidaId int, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := transactions.ClaimVidaIdTransaction(
		int64(vidaId), w.GetNonce(), w.Address, feePerByte,
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
	vidaId int, wrappedTxns [][]byte, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := transactions.ConduitApprovalTransaction(
		int64(vidaId), wrappedTxns, w.GetNonce(), w.Address, feePerByte,
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

func (w *PWRWallet) SendPayableVidaData(
	vidaId int, data []byte, value int64, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := transactions.PayableVidaDataTransaction(
		int64(vidaId), data, value, w.GetNonce(), w.Address, feePerByte,
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

func (w *PWRWallet) SendVidaData(
	vidaId int, data []byte, feePerByte int,
) rpc.BroadcastResponse {
	return w.SendPayableVidaData(vidaId, data, 0, feePerByte)
}

func (w *PWRWallet) RemoveConduits(
	vidaId int, conduits [][]byte, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := transactions.RemoveConduitsTransaction(
		int64(vidaId), conduits, w.GetNonce(), w.Address, feePerByte,
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
	vidaId int, mode byte, conduitThreshold int, conduits [][]byte,
	conduitsWithVotingPower map[string]int64, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := transactions.SetConduitModeTransaction(
		int64(vidaId), mode, conduitThreshold, conduits, conduitsWithVotingPower,
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
	vidaId int, privateState bool, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := transactions.SetVidaPrivateStateTransaction(
		int64(vidaId), privateState, w.GetNonce(), w.Address, feePerByte,
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
	vidaId int, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := transactions.SetVidaToAbsolutePublicTransaction(
		int64(vidaId), w.GetNonce(), w.Address, feePerByte,
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
	vidaId int, sponsoredAddresses []string, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := transactions.AddVidaSponsoredAddressesTransaction(
		int64(vidaId), sponsoredAddresses, w.GetNonce(), w.Address, feePerByte,
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
	vidaId int, allowedSenders []string, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := transactions.AddVidaAllowedSendersTransaction(
		int64(vidaId), allowedSenders, w.GetNonce(), w.Address, feePerByte,
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
	vidaId int, allowedSenders []string, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := transactions.RemoveVidaAllowedSendersTransaction(
		int64(vidaId), allowedSenders, w.GetNonce(), w.Address, feePerByte,
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
	vidaId int, sponsoredAddresses []string, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := transactions.RemoveSponsoredAddressesTransaction(
		int64(vidaId), sponsoredAddresses, w.GetNonce(), w.Address, feePerByte,
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
	vidaId int, ownerCanTransferPWR bool, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := transactions.SetPWRTransferRightsTransaction(
		int64(vidaId), ownerCanTransferPWR, w.GetNonce(), w.Address, feePerByte,
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
	vidaId int, receiver string, amount int64, feePerByte int,
) rpc.BroadcastResponse {
	response := makeSurePublicKeyIsSet(feePerByte, w)
	if response != nil && !response.Success {
		return *response
	}

	var buffer []byte
	buffer, err := transactions.TransferPWRFromVidaTransaction(
		int64(vidaId), receiver, amount, w.GetNonce(), w.Address, feePerByte,
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
