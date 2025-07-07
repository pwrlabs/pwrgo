package transactions

import (
	"encoding/hex"
)

// Falcon transaction bytes
func SetPublicKeyBytes(
	publicKey []byte, nonce int, address []byte, feePerByte int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1001, nonce, address, feePerByte)

	publicKeyLenBytes := DecToBytes(len(publicKey), 2)

	txnBytes = append(txnBytes, publicKeyLenBytes...)
	txnBytes = append(txnBytes, publicKey...)

	return txnBytes, nil
}

func JoinAsValidatorBytes(ip string, nonce int, address []byte, feePerByte int) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1002, nonce, address, feePerByte)

	ipBytes := []byte(ip)
	ipLenBytes := DecToBytes(len(ipBytes), 2)

	txnBytes = append(txnBytes, ipLenBytes...)
	txnBytes = append(txnBytes, ipBytes...)
	return txnBytes, nil
}

func DelegateTxBytes(to string, amount int, nonce int, address []byte, feePerByte int) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1003, nonce, address, feePerByte)

	recipientBytes, _ := hex.DecodeString(to[2:])
	amountBytes := decToBytes64(int64(amount), 8)

	paddedRecipient := make([]byte, 20)
	copy(paddedRecipient[20-len(recipientBytes):], recipientBytes)

	txnBytes = append(txnBytes, paddedRecipient...)
	txnBytes = append(txnBytes, amountBytes...)
	return txnBytes, nil
}

func ChangeIpBytes(newIp string, nonce int, address []byte, feePerByte int) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1004, nonce, address, feePerByte)

	newIpBytes := []byte(newIp)
	newIpLenBytes := DecToBytes(len(newIpBytes), 2)

	txnBytes = append(txnBytes, newIpLenBytes...)
	txnBytes = append(txnBytes, newIpBytes...)
	return txnBytes, nil
}

func ClaimActiveNodeSpotBytes(nonce int, address []byte, feePerByte int) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1005, nonce, address, feePerByte)

	return txnBytes, nil
}

func TransferTxBytes(
	amount int, recipient string, nonce int, address []byte, feePerByte int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1006, nonce, address, feePerByte)

	recipientBytes, err := hex.DecodeString(recipient[2:])
	if err != nil {
		return nil, err
	}

	paddedRecipient := make([]byte, 20)
	copy(paddedRecipient[20-len(recipientBytes):], recipientBytes)
	amountBytes := decToBytes64(int64(amount), 8)

	txnBytes = append(txnBytes, paddedRecipient...)
	txnBytes = append(txnBytes, amountBytes...)

	return txnBytes, nil
}

func EarlyWithdrawPenaltyProposal(
	title string, description string, withdraw_penalty_time int,
	withdraw_penalty int, nonce int, address []byte, feePerByte int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1009, nonce, address, feePerByte)

	titleBytes := []byte(title)
	titleLength := len(titleBytes)
	titleLengthBytes := DecToBytes(titleLength, 4)
	descriptionBytes := []byte(description)
	withdrawPenaltyTimeBytes := decToBytes64(int64(withdraw_penalty_time), 8)
	withdrawPenaltyBytes := DecToBytes(withdraw_penalty, 4)

	paddedTitleLength := make([]byte, 4)
	copy(paddedTitleLength[4-len(titleLengthBytes):], titleLengthBytes)
	paddedWithdrawPenalty := make([]byte, 4)
	copy(paddedWithdrawPenalty[4-len(withdrawPenaltyBytes):], withdrawPenaltyBytes)

	txnBytes = append(txnBytes, paddedTitleLength...)
	txnBytes = append(txnBytes, titleBytes...)
	txnBytes = append(txnBytes, withdrawPenaltyTimeBytes...)
	txnBytes = append(txnBytes, paddedWithdrawPenalty...)
	txnBytes = append(txnBytes, descriptionBytes...)

	return txnBytes, nil
}

func ChangeFeePerByteProposalTx(
	title string, description string, feePerByte int, nonce int, address []byte, _feePerByte int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1010, nonce, address, _feePerByte)

	titleBytes := []byte(title)
	titleLengthBytes := DecToBytes(len(titleBytes), 4)
	descriptionBytes := []byte(description)
	feePerByteBytes := decToBytes64(int64(feePerByte), 8)

	paddedTitleLength := make([]byte, 4)
	copy(paddedTitleLength[4-len(titleLengthBytes):], titleLengthBytes)

	txnBytes = append(txnBytes, paddedTitleLength...)
	txnBytes = append(txnBytes, titleBytes...)
	txnBytes = append(txnBytes, feePerByteBytes...)
	txnBytes = append(txnBytes, descriptionBytes...)

	return txnBytes, nil
}

func ChangeMaxBlockSizeProposalTx(
	title string, description string, max_block_size int, nonce int, address []byte, feePerByte int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1011, nonce, address, feePerByte)

	titleBytes := []byte(title)
	titleLength := len(titleBytes)
	titleLengthBytes := DecToBytes(titleLength, 4)
	descriptionBytes := []byte(description)
	maxBlockSizeBytes := DecToBytes(max_block_size, 4)

	paddedTitleLength := make([]byte, 4)
	copy(paddedTitleLength[4-len(titleLengthBytes):], titleLengthBytes)
	paddedMaxBlockSize := make([]byte, 4)
	copy(paddedMaxBlockSize[4-len(maxBlockSizeBytes):], maxBlockSizeBytes)

	txnBytes = append(txnBytes, paddedTitleLength...)
	txnBytes = append(txnBytes, titleBytes...)
	txnBytes = append(txnBytes, paddedMaxBlockSize...)
	txnBytes = append(txnBytes, descriptionBytes...)

	return txnBytes, nil
}

func ChangeMaxTxnSizeProposalTx(
	title string, description string, maxTxnSize int, nonce int, address []byte, feePerByte int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1012, nonce, address, feePerByte)

	titleBytes := []byte(title)
	titleLengthBytes := DecToBytes(len(titleBytes), 4)
	maxTxnSizeBytes := DecToBytes(maxTxnSize, 4)
	descriptionBytes := []byte(description)

	txnBytes = append(txnBytes, titleLengthBytes...)
	txnBytes = append(txnBytes, titleBytes...)
	txnBytes = append(txnBytes, maxTxnSizeBytes...)
	txnBytes = append(txnBytes, descriptionBytes...)

	return txnBytes, nil
}

func ChangeOverallBurnPercentageProposalTx(
	title string, description string, burnPercentage int, nonce int, address []byte, feePerByte int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1013, nonce, address, feePerByte)

	titleBytes := []byte(title)
	titleLengthBytes := DecToBytes(len(titleBytes), 4)
	burnPercentageBytes := DecToBytes(burnPercentage, 4)
	descriptionBytes := []byte(description)

	txnBytes = append(txnBytes, titleLengthBytes...)
	txnBytes = append(txnBytes, titleBytes...)
	txnBytes = append(txnBytes, burnPercentageBytes...)
	txnBytes = append(txnBytes, descriptionBytes...)

	return txnBytes, nil
}

func ChangeRewardPerYearProposalTx(
	title string, description string, rewardPerYear int64, nonce int, address []byte, feePerByte int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1014, nonce, address, feePerByte)

	titleBytes := []byte(title)
	titleLengthBytes := DecToBytes(len(titleBytes), 4)
	rewardPerYearBytes := decToBytes64(rewardPerYear, 8)
	descriptionBytes := []byte(description)

	txnBytes = append(txnBytes, titleLengthBytes...)
	txnBytes = append(txnBytes, titleBytes...)
	txnBytes = append(txnBytes, rewardPerYearBytes...)
	txnBytes = append(txnBytes, descriptionBytes...)

	return txnBytes, nil
}

func ChangeValidatorCountLimitProposalTx(
	title string, description string, validatorCountLimit int, nonce int, address []byte, feePerByte int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1015, nonce, address, feePerByte)

	titleBytes := []byte(title)
	titleLengthBytes := DecToBytes(len(titleBytes), 4)
	validatorCountLimitBytes := DecToBytes(validatorCountLimit, 4)
	descriptionBytes := []byte(description)

	txnBytes = append(txnBytes, titleLengthBytes...)
	txnBytes = append(txnBytes, titleBytes...)
	txnBytes = append(txnBytes, validatorCountLimitBytes...)
	txnBytes = append(txnBytes, descriptionBytes...)

	return txnBytes, nil
}

func ChangeValidatorJoiningFeeProposalTx(
	title string, description string, joiningFee int64, nonce int, address []byte, feePerByte int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1016, nonce, address, feePerByte)

	titleBytes := []byte(title)
	titleLengthBytes := DecToBytes(len(titleBytes), 4)
	joiningFeeBytes := decToBytes64(joiningFee, 8)
	descriptionBytes := []byte(description)

	txnBytes = append(txnBytes, titleLengthBytes...)
	txnBytes = append(txnBytes, titleBytes...)
	txnBytes = append(txnBytes, joiningFeeBytes...)
	txnBytes = append(txnBytes, descriptionBytes...)

	return txnBytes, nil
}

func ChangeVidaIdClaimingFeeProposalTx(
	title string, description string, vidaIdClaimingFee int64, nonce int, address []byte, feePerByte int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1017, nonce, address, feePerByte)

	titleBytes := []byte(title)
	titleLengthBytes := DecToBytes(len(titleBytes), 4)
	vidaIdClaimingFeeBytes := decToBytes64(vidaIdClaimingFee, 8)
	descriptionBytes := []byte(description)

	txnBytes = append(txnBytes, titleLengthBytes...)
	txnBytes = append(txnBytes, titleBytes...)
	txnBytes = append(txnBytes, vidaIdClaimingFeeBytes...)
	txnBytes = append(txnBytes, descriptionBytes...)

	return txnBytes, nil
}

func ChangeVmOwnerTxnFeeShareProposalTx(
	title string, description string, vmOwnerTxnFeeShare int, nonce int, address []byte, feePerByte int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1018, nonce, address, feePerByte)

	titleBytes := []byte(title)
	titleLengthBytes := DecToBytes(len(titleBytes), 4)
	vmOwnerTxnFeeShareBytes := DecToBytes(vmOwnerTxnFeeShare, 4)
	descriptionBytes := []byte(description)

	txnBytes = append(txnBytes, titleLengthBytes...)
	txnBytes = append(txnBytes, titleBytes...)
	txnBytes = append(txnBytes, vmOwnerTxnFeeShareBytes...)
	txnBytes = append(txnBytes, descriptionBytes...)

	return txnBytes, nil
}

func OtherProposalTx(
	title string, description string, nonce int, address []byte, feePerByte int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1019, nonce, address, feePerByte)

	titleBytes := []byte(title)
	titleLengthBytes := DecToBytes(len(titleBytes), 4)
	descriptionBytes := []byte(description)

	txnBytes = append(txnBytes, titleLengthBytes...)
	txnBytes = append(txnBytes, titleBytes...)
	txnBytes = append(txnBytes, descriptionBytes...)

	return txnBytes, nil
}

func VoteOnProposalTx(
	proposalHash []byte, vote byte, nonce int, address []byte, feePerByte int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1020, nonce, address, feePerByte)

	txnBytes = append(txnBytes, proposalHash...)
	txnBytes = append(txnBytes, vote)

	return txnBytes, nil
}

func GuardianApprovalTransaction(
	wrappedTxns [][]byte, nonce int, address []byte, feePerByte int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1021, nonce, address, feePerByte)

	// Calculate total size needed for all wrapped transactions
	totalWrappedSize := 0
	for _, wrappedTxn := range wrappedTxns {
		totalWrappedSize += 4 + len(wrappedTxn) // 4 bytes for length + txn size
	}

	buffer := make([]byte, len(txnBytes)+4+totalWrappedSize)
	copy(buffer, txnBytes)

	// Add number of wrapped transactions
	numTxnsBytes := DecToBytes(len(wrappedTxns), 4)
	copy(buffer[len(txnBytes):], numTxnsBytes)

	// Add each wrapped transaction with its length prefix
	offset := len(txnBytes) + 4
	for _, wrappedTxn := range wrappedTxns {
		lengthBytes := DecToBytes(len(wrappedTxn), 4)
		copy(buffer[offset:], lengthBytes)
		offset += 4
		copy(buffer[offset:], wrappedTxn)
		offset += len(wrappedTxn)
	}

	return buffer, nil
}

func RemoveGuardianTransaction(nonce int, address []byte, feePerByte int) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1022, nonce, address, feePerByte)
	return txnBytes, nil
}

func SetGuardianTransaction(
	expiryDate int, guardianAddress string, nonce int, address []byte, feePerByte int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1023, nonce, address, feePerByte)

	expiryDateBytes := decToBytes64(int64(expiryDate), 8)
	guardianAddressBytes, _ := hex.DecodeString(guardianAddress[2:])

	buffer := make([]byte, len(txnBytes)+8+20)
	copy(buffer, txnBytes)
	copy(buffer[len(txnBytes):], expiryDateBytes)
	copy(buffer[len(txnBytes)+8:], guardianAddressBytes)

	return buffer, nil
}

func MoveStakeTransaction(
	sharesAmount int64, fromValidator string, toValidator string, nonce int, address []byte, feePerByte int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1024, nonce, address, feePerByte)

	sharesAmountBytes := decToBytes64(sharesAmount, 8)
	fromValidatorBytes, _ := hex.DecodeString(fromValidator[2:])
	toValidatorBytes, _ := hex.DecodeString(toValidator[2:])

	buffer := make([]byte, len(txnBytes)+2+len(sharesAmountBytes)+20+20)
	copy(buffer, txnBytes)
	copy(buffer[len(txnBytes):], DecToBytes(len(sharesAmountBytes), 2))
	copy(buffer[len(txnBytes)+2:], sharesAmountBytes)
	copy(buffer[len(txnBytes)+2+len(sharesAmountBytes):], fromValidatorBytes)
	copy(buffer[len(txnBytes)+2+len(sharesAmountBytes)+20:], toValidatorBytes)

	return buffer, nil
}

func RemoveValidatorTransaction(
	validatorAddress string, nonce int, address []byte, feePerByte int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1025, nonce, address, feePerByte)

	validatorAddressBytes, _ := hex.DecodeString(validatorAddress[2:])

	buffer := make([]byte, len(txnBytes)+20)
	copy(buffer, txnBytes)
	copy(buffer[len(txnBytes):], validatorAddressBytes)

	return buffer, nil
}

func WithdrawTransaction(
	sharesAmount int, validator string, nonce int, address []byte, feePerByte int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1026, nonce, address, feePerByte)

	sharesAmountBytes := decToBytes64(int64(sharesAmount), 8)
	validatorBytes, _ := hex.DecodeString(validator[2:])

	buffer := make([]byte, len(txnBytes)+2+len(sharesAmountBytes)+20)
	copy(buffer, txnBytes)
	copy(buffer[len(txnBytes):], DecToBytes(len(sharesAmountBytes), 2))
	copy(buffer[len(txnBytes)+2:], sharesAmountBytes)
	copy(buffer[len(txnBytes)+2+len(sharesAmountBytes):], validatorBytes)

	return buffer, nil
}

func ClaimVidaIdTransaction(
	vidaId int64, nonce int, address []byte, feePerByte int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1028, nonce, address, feePerByte)

	buffer := make([]byte, len(txnBytes)+8)
	copy(buffer, txnBytes)
	copy(buffer[len(txnBytes):], decToBytes64(vidaId, 8))

	return buffer, nil
}

func ConduitApprovalTransaction(
	vidaId int64, wrappedTxns [][]byte, nonce int, address []byte, feePerByte int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1029, nonce, address, feePerByte)

	// Calculate total size needed
	totalWrappedSize := 8 // vidaId
	for _, wrappedTxn := range wrappedTxns {
		totalWrappedSize += 4 + len(wrappedTxn)
	}

	buffer := make([]byte, len(txnBytes)+4+totalWrappedSize)
	copy(buffer, txnBytes)
	copy(buffer[len(txnBytes):], decToBytes64(vidaId, 8))
	copy(buffer[len(txnBytes)+8:], DecToBytes(len(wrappedTxns), 4))

	offset := len(txnBytes) + 12
	for _, wrappedTxn := range wrappedTxns {
		lengthBytes := DecToBytes(len(wrappedTxn), 4)
		copy(buffer[offset:], lengthBytes)
		offset += 4
		copy(buffer[offset:], wrappedTxn)
		offset += len(wrappedTxn)
	}

	return buffer, nil
}

func PayableVidaDataTransaction(
	vidaId int64, data []byte, value int64, nonce int, address []byte, feePerByte int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1030, nonce, address, feePerByte)

	buffer := make([]byte, len(txnBytes)+8+4+len(data)+8)
	copy(buffer, txnBytes)
	copy(buffer[len(txnBytes):], decToBytes64(vidaId, 8))
	copy(buffer[len(txnBytes)+8:], DecToBytes(len(data), 4))
	copy(buffer[len(txnBytes)+12:], data)
	copy(buffer[len(txnBytes)+12+len(data):], decToBytes64(value, 8))

	return buffer, nil
}

func RemoveConduitsTransaction(
	vidaId int64, conduits [][]byte, nonce int, address []byte, feePerByte int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1031, nonce, address, feePerByte)

	buffer := make([]byte, len(txnBytes)+8+(len(conduits)*20))
	copy(buffer, txnBytes)
	copy(buffer[len(txnBytes):], decToBytes64(vidaId, 8))

	offset := len(txnBytes) + 8
	for _, conduit := range conduits {
		copy(buffer[offset:], conduit)
		offset += 20
	}

	return buffer, nil
}

func SetConduitModeTransaction(
	vidaId int64, mode byte, conduitThreshold int, conduits [][]byte, conduitsWithVotingPower map[string]int64,
	nonce int, address []byte, feePerByte int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1033, nonce, address, feePerByte)

	// Calculate total size
	totalSize := len(txnBytes) + 8 + 1 + 4 + 4
	if conduits != nil {
		totalSize += len(conduits) * 20
	}
	if conduitsWithVotingPower != nil {
		totalSize += len(conduitsWithVotingPower) * 28 // 20 bytes address + 8 bytes voting power
	}

	buffer := make([]byte, totalSize)
	copy(buffer, txnBytes)
	copy(buffer[len(txnBytes):], decToBytes64(vidaId, 8))
	buffer[len(txnBytes)+8] = mode
	copy(buffer[len(txnBytes)+9:], DecToBytes(conduitThreshold, 4))

	offset := len(txnBytes) + 13
	if conduits != nil {
		copy(buffer[offset:], DecToBytes(len(conduits), 4))
		offset += 4
		for _, conduit := range conduits {
			copy(buffer[offset:], conduit)
			offset += 20
		}
	}

	if conduitsWithVotingPower != nil && len(conduitsWithVotingPower) > 0 {
		copy(buffer[offset:], DecToBytes(len(conduitsWithVotingPower), 4))
		offset += 4
		for addr, power := range conduitsWithVotingPower {
			addrBytes, _ := hex.DecodeString(addr[2:])
			copy(buffer[offset:], addrBytes)
			offset += 20
			copy(buffer[offset:], decToBytes64(power, 8))
			offset += 8
		}
	}

	if (conduits == nil || len(conduits) == 0) && (conduitsWithVotingPower == nil || len(conduitsWithVotingPower) == 0) {
		copy(buffer[offset:], DecToBytes(0, 4))
	}

	return buffer, nil
}

func SetVidaPrivateStateTransaction(
	vidaId int64, privateState bool, nonce int, address []byte, feePerByte int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1034, nonce, address, feePerByte)

	buffer := make([]byte, len(txnBytes)+8+1)
	copy(buffer, txnBytes)
	copy(buffer[len(txnBytes):], decToBytes64(vidaId, 8))
	if privateState {
		buffer[len(txnBytes)+8] = 1
	} else {
		buffer[len(txnBytes)+8] = 0
	}

	return buffer, nil
}

func SetVidaToAbsolutePublicTransaction(
	vidaId int64, nonce int, address []byte, feePerByte int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1035, nonce, address, feePerByte)

	buffer := make([]byte, len(txnBytes)+8)
	copy(buffer, txnBytes)
	copy(buffer[len(txnBytes):], decToBytes64(vidaId, 8))

	return buffer, nil
}

func AddVidaSponsoredAddressesTransaction(
	vidaId int64, sponsoredAddresses []string, nonce int, address []byte, feePerByte int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1036, nonce, address, feePerByte)

	sponsoredAddressesBytes := make([][]byte, len(sponsoredAddresses))
	for i, addr := range sponsoredAddresses {
		sponsoredAddressesBytes[i], _ = hex.DecodeString(addr[2:])
	}

	buffer := make([]byte, len(txnBytes)+8+(len(sponsoredAddressesBytes)*20))
	copy(buffer, txnBytes)
	copy(buffer[len(txnBytes):], decToBytes64(vidaId, 8))

	offset := len(txnBytes) + 8
	for _, addr := range sponsoredAddresses {
		copy(buffer[offset:], addr)
		offset += 20
	}

	return buffer, nil
}

func AddVidaAllowedSendersTransaction(
	vidaId int64, allowedSenders []string, nonce int, address []byte, feePerByte int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1037, nonce, address, feePerByte)

	allowedSendersBytes := make([][]byte, len(allowedSenders))
	for i, addr := range allowedSenders {
		allowedSendersBytes[i], _ = hex.DecodeString(addr[2:])
	}

	buffer := make([]byte, len(txnBytes)+8+(len(allowedSendersBytes)*20))
	copy(buffer, txnBytes)
	copy(buffer[len(txnBytes):], decToBytes64(vidaId, 8))

	offset := len(txnBytes) + 8
	for _, addr := range allowedSendersBytes {
		copy(buffer[offset:], addr)
		offset += 20
	}

	return buffer, nil
}

func RemoveVidaAllowedSendersTransaction(
	vidaId int64, allowedSenders []string, nonce int, address []byte, feePerByte int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1038, nonce, address, feePerByte)

	allowedSendersBytes := make([][]byte, len(allowedSenders))
	for i, addr := range allowedSenders {
		allowedSendersBytes[i], _ = hex.DecodeString(addr[2:])
	}

	buffer := make([]byte, len(txnBytes)+8+(len(allowedSendersBytes)*20))
	copy(buffer, txnBytes)
	copy(buffer[len(txnBytes):], decToBytes64(vidaId, 8))

	offset := len(txnBytes) + 8
	for _, addr := range allowedSendersBytes {
		copy(buffer[offset:], addr)
		offset += 20
	}

	return buffer, nil
}

func RemoveSponsoredAddressesTransaction(
	vidaId int64, sponsoredAddresses []string, nonce int, address []byte, feePerByte int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1039, nonce, address, feePerByte)

	sponsoredAddressesBytes := make([][]byte, len(sponsoredAddresses))
	for i, addr := range sponsoredAddresses {
		sponsoredAddressesBytes[i], _ = hex.DecodeString(addr[2:])
	}

	buffer := make([]byte, len(txnBytes)+8+(len(sponsoredAddressesBytes)*20))
	copy(buffer, txnBytes)
	copy(buffer[len(txnBytes):], decToBytes64(vidaId, 8))

	offset := len(txnBytes) + 8
	for _, addr := range sponsoredAddressesBytes {
		copy(buffer[offset:], addr)
		offset += 20
	}

	return buffer, nil
}

func SetPWRTransferRightsTransaction(
	vidaId int64, ownerCanTransferPWR bool, nonce int, address []byte, feePerByte int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1040, nonce, address, feePerByte)

	buffer := make([]byte, len(txnBytes)+8+1)
	copy(buffer, txnBytes)
	copy(buffer[len(txnBytes):], decToBytes64(vidaId, 8))
	if ownerCanTransferPWR {
		buffer[len(txnBytes)+8] = 1
	} else {
		buffer[len(txnBytes)+8] = 0
	}

	return buffer, nil
}

func TransferPWRFromVidaTransaction(
	vidaId int64, receiver string, amount int64, nonce int, address []byte, feePerByte int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(1041, nonce, address, feePerByte)

	receiverBytes, _ := hex.DecodeString(receiver[2:])

	buffer := make([]byte, len(txnBytes)+8+20+8)
	copy(buffer, txnBytes)
	copy(buffer[len(txnBytes):], decToBytes64(vidaId, 8))
	copy(buffer[len(txnBytes)+8:], receiverBytes)
	copy(buffer[len(txnBytes)+28:], decToBytes64(amount, 8))

	return buffer, nil
}
