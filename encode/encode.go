package encode

import (
	"encoding/hex"
	"errors"
)

func ClaimVMIdBytes(vmId int, nonce int) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(claimVmIdType, nonce)
	vmIdBytes := decToBytes64(int64(vmId), 8)

	txnBytes = append(txnBytes, vmIdBytes...)

	return txnBytes, nil
}

func ClaimActiveNodeSpotBytes(nonce int) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(claimSpotType, nonce)
	return txnBytes, nil
}

func JoinBytes(ip string, nonce int) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(joinType, nonce)

	ipBytes := []byte(ip)
	txnBytes = append(txnBytes, ipBytes...)
	return txnBytes, nil
}

func TransferTxBytes(amount int, recipient string, nonce int) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(transferType, nonce)

	amountBytes := decToBytes64(int64(amount), 8)
	recipientBytes, err := hex.DecodeString(recipient[2:])
	if err != nil {
		return nil, err
	}

	paddedRecipient := make([]byte, 20)
	copy(paddedRecipient[20-len(recipientBytes):], recipientBytes)

	txnBytes = append(txnBytes, amountBytes...)
	txnBytes = append(txnBytes, paddedRecipient...)

	return txnBytes, nil
}

func VmDataBytes(vmId int, data []byte, nonce int) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(vmDataType, nonce)
	vmIdBytes := decToBytes64(int64(vmId), 8)
	dataLengthBytes := DecToBytes(int(len(data)), 4)

	txnBytes = append(txnBytes, vmIdBytes...)
	txnBytes = append(txnBytes, dataLengthBytes...)
	txnBytes = append(txnBytes, data...)

	return txnBytes, nil
}

func ValidatorRemoveBytes(validatorAddress string, nonce int) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(validatorRemoveType, nonce)

	validatorBytes, _ := hex.DecodeString(validatorAddress[2:])
	paddedValidator := make([]byte, 20)
	copy(paddedValidator[20-len(validatorBytes):], validatorBytes)

	txnBytes = append(txnBytes, paddedValidator...)
	return txnBytes, nil
}

func GuardianApprovalTransaction(txn []byte, nonce int) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(guardianApprovalType, nonce)

	txnBytes = append(txnBytes, txn...)
	return txnBytes, nil
}

func PayableVmDataBytes(vmId int, data []byte, amount int, nonce int) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(payableVmDataType, nonce)
	vmIdBytes := decToBytes64(int64(vmId), 8)
	amountBytes := decToBytes64(int64(amount), 8)
	dataLengthBytes := DecToBytes(int(len(data)), 4)

	txnBytes = append(txnBytes, vmIdBytes...)
	txnBytes = append(txnBytes, dataLengthBytes...)
	txnBytes = append(txnBytes, data...)
	txnBytes = append(txnBytes, amountBytes...)

	return txnBytes, nil
}

func DelegateTxBytes(to string, amount int, nonce int) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(delegateType, nonce)

	amountBytes := decToBytes64(int64(amount), 8)
	recipientBytes, _ := hex.DecodeString(to[2:])

	paddedRecipient := make([]byte, 20)
	copy(paddedRecipient[20-len(recipientBytes):], recipientBytes)

	txnBytes = append(txnBytes, amountBytes...)
	txnBytes = append(txnBytes, paddedRecipient...)
	return txnBytes, nil
}

func WithdrawTxBytes(from string, amount int, nonce int) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(withdrawType, nonce)

	amountBytes := decToBytes64(int64(amount), 8)

	recipientBytes, _ := hex.DecodeString(from[2:])
	paddedRecipient := make([]byte, 20)
	copy(paddedRecipient[20-len(recipientBytes):], recipientBytes)

	txnBytes = append(txnBytes, amountBytes...)
	txnBytes = append(txnBytes, paddedRecipient...)
	return txnBytes, nil
}

func SetGuardianTxBytes(guardian string, expiration int, nonce int) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(setGuardianType, nonce)

	expirationBytes := decToBytes64(int64(expiration), 8)
	recipientBytes, _ := hex.DecodeString(guardian[2:])

	paddedRecipient := make([]byte, 20)
	copy(paddedRecipient[20-len(recipientBytes):], recipientBytes)

	txnBytes = append(txnBytes, expirationBytes...)
	txnBytes = append(txnBytes, paddedRecipient...)
	return txnBytes, nil
}

func RemoveGuardianTxBytes(nonce int) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(removeGuardianType, nonce)
	return txnBytes, nil
}

func SetConduitTxBytes(vmId int, conduits []string, nonce int) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(setConduitsType, nonce)

	vmIdBytes := decToBytes64(int64(vmId), 8)
	txnBytes = append(txnBytes, vmIdBytes...)

	var lengthsBytes []byte
	var conduitsBytes []byte

	for _, c := range conduits {
		decoded, err := hex.DecodeString(c[2:])

		if err != nil {
			return nil, errors.New("invalid conduit address")
		}

		lengthBytes := DecToBytes(len(decoded), 4)
		paddedLength := make([]byte, 4)
		copy(paddedLength[4-len(lengthBytes):], lengthBytes)

		lengthsBytes = append(lengthsBytes, paddedLength...)
		conduitsBytes = append(conduitsBytes, decoded...)
	}

	txnBytes = append(txnBytes, lengthsBytes...)
	txnBytes = append(txnBytes, conduitsBytes...)
	return txnBytes, nil
}

func AddConduitTxBytes(vmId int, transaction []byte, nonce int) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(addConduitsType, nonce)

	vmIdBytes := decToBytes64(int64(vmId), 8)

	txnBytes = append(txnBytes, vmIdBytes...)
	txnBytes = append(txnBytes, transaction...)

	return txnBytes, nil
}

func MoveStakeTxBytes(shares int, from_validator string, to_validator string, nonce int) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(moveStakeType, nonce)

	sharesBytes := decToBytes64(int64(shares), 8)
	fromValidatorBytes, err1 := hex.DecodeString(from_validator[2:])
	toValidatorBytes, err2 := hex.DecodeString(to_validator[2:])

	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}

	paddedFromValidator := make([]byte, 20)
	copy(paddedFromValidator[20-len(fromValidatorBytes):], fromValidatorBytes)

	paddedToValidator := make([]byte, 20)
	copy(paddedToValidator[20-len(toValidatorBytes):], toValidatorBytes)

	txnBytes = append(txnBytes, sharesBytes...)
	txnBytes = append(txnBytes, paddedFromValidator...)
	txnBytes = append(txnBytes, paddedToValidator...)

	return txnBytes, nil
}

func ChangeEarlyWithdrawPenaltyProposalTxBytes(
	title string, description string, withdraw_penalty_time int, withdraw_penalty int, nonce int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(changeEarlyWithdrawPenaltyProposalType, nonce)

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

func ChangeFeePerByteProposalTxBytes(
	title string, description string, fee_per_byte int, nonce int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(changeFeePerByteProposalType, nonce)

	titleBytes := []byte(title)
	titleLengthBytes := DecToBytes(len(titleBytes), 4)
	descriptionBytes := []byte(description)
	feePerByteBytes := decToBytes64(int64(fee_per_byte), 8)

	paddedTitleLength := make([]byte, 4)
	copy(paddedTitleLength[4-len(titleLengthBytes):], titleLengthBytes)

	txnBytes = append(txnBytes, paddedTitleLength...)
	txnBytes = append(txnBytes, titleBytes...)
	txnBytes = append(txnBytes, feePerByteBytes...)
	txnBytes = append(txnBytes, descriptionBytes...)

	return txnBytes, nil
}

func ChangeMaxBlockSizeProposalTxBytes(
	title string, description string, max_block_size int, nonce int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(changeMaxBlockSizeProposalType, nonce)

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

func ChangeMaxTxnSizeProposalTxBytes(
	title string, description string, max_txn_size int, nonce int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(changeMaxTxnSizeProposalType, nonce)

	titleBytes := []byte(title)
	titleLength := len(titleBytes)
	titleLengthBytes := DecToBytes(titleLength, 4)
	descriptionBytes := []byte(description)
	maxTxnSizeBytes := DecToBytes(max_txn_size, 4)

	paddedTitleLength := make([]byte, 4)
	copy(paddedTitleLength[4-len(titleLengthBytes):], titleLengthBytes)
	paddedMaxTxnSize := make([]byte, 4)
	copy(paddedMaxTxnSize[4-len(maxTxnSizeBytes):], maxTxnSizeBytes)

	txnBytes = append(txnBytes, paddedTitleLength...)
	txnBytes = append(txnBytes, titleBytes...)
	txnBytes = append(txnBytes, paddedMaxTxnSize...)
	txnBytes = append(txnBytes, descriptionBytes...)

	return txnBytes, nil
}

func ChangeOverallBurnPercentageProposalTxBytes(
	title string, description string, burn_percentage int, nonce int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(changeOverallBurnPercentageProposalType, nonce)

	titleBytes := []byte(title)
	titleLength := len(titleBytes)
	titleLengthBytes := DecToBytes(titleLength, 4)
	descriptionBytes := []byte(description)
	burnPercentageBytes := DecToBytes(burn_percentage, 4)

	paddedTitleLength := make([]byte, 4)
	copy(paddedTitleLength[4-len(titleLengthBytes):], titleLengthBytes)
	paddedBurnPercentage := make([]byte, 4)
	copy(paddedBurnPercentage[4-len(burnPercentageBytes):], burnPercentageBytes)

	txnBytes = append(txnBytes, paddedTitleLength...)
	txnBytes = append(txnBytes, titleBytes...)
	txnBytes = append(txnBytes, paddedBurnPercentage...)
	txnBytes = append(txnBytes, descriptionBytes...)

	return txnBytes, nil
}

func ChangeRewardPerYearProposalTxBytes(
	title string, description string, reward_per_year int, nonce int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(changeRewardPerYearProposalType, nonce)

	titleBytes := []byte(title)
	titleLengthBytes := DecToBytes(len(titleBytes), 4)
	descriptionBytes := []byte(description)
	rewardPerYearBytes := decToBytes64(int64(reward_per_year), 8)

	paddedTitleLength := make([]byte, 4)
	copy(paddedTitleLength[4-len(titleLengthBytes):], titleLengthBytes)

	txnBytes = append(txnBytes, paddedTitleLength...)
	txnBytes = append(txnBytes, titleBytes...)
	txnBytes = append(txnBytes, rewardPerYearBytes...)
	txnBytes = append(txnBytes, descriptionBytes...)

	return txnBytes, nil
}

func ChangeValidatorCountLimitProposalTxBytes(
	title string, description string, validator_count_limit int, nonce int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(changeValidatorCountLimitProposalType, nonce)

	titleBytes := []byte(title)
	titleLengthBytes := DecToBytes(len(titleBytes), 4)
	descriptionBytes := []byte(description)
	validatorCountLimitBytes := DecToBytes(validator_count_limit, 4)

	paddedTitleLength := make([]byte, 4)
	copy(paddedTitleLength[4-len(titleLengthBytes):], titleLengthBytes)
	paddedValidatorCountLimit := make([]byte, 4)
	copy(paddedValidatorCountLimit[4-len(validatorCountLimitBytes):], validatorCountLimitBytes)

	txnBytes = append(txnBytes, paddedTitleLength...)
	txnBytes = append(txnBytes, titleBytes...)
	txnBytes = append(txnBytes, paddedValidatorCountLimit...)
	txnBytes = append(txnBytes, descriptionBytes...)

	return txnBytes, nil
}

func ChangeValidatorJoiningFeeProposalTxBytes(
	title string, description string, joining_fee int, nonce int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(changeValidatorJoiningFeeProposalType, nonce)

	titleBytes := []byte(title)
	titleLengthBytes := DecToBytes(len(titleBytes), 4)
	descriptionBytes := []byte(description)
	joiningFeeBytes := decToBytes64(int64(joining_fee), 8)

	paddedTitleLength := make([]byte, 4)
	copy(paddedTitleLength[4-len(titleLengthBytes):], titleLengthBytes)

	txnBytes = append(txnBytes, paddedTitleLength...)
	txnBytes = append(txnBytes, titleBytes...)
	txnBytes = append(txnBytes, joiningFeeBytes...)
	txnBytes = append(txnBytes, descriptionBytes...)

	return txnBytes, nil
}

func ChangeVmIdClaimingFeeProposalTxBytes(
	title string, description string, claiming_fee int, nonce int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(changeVmIdClaimingFeeProposalType, nonce)

	titleBytes := []byte(title)
	titleLengthBytes := DecToBytes(len(titleBytes), 4)
	descriptionBytes := []byte(description)
	claimingFeeBytes := decToBytes64(int64(claiming_fee), 8)

	paddedTitleLength := make([]byte, 4)
	copy(paddedTitleLength[4-len(titleLengthBytes):], titleLengthBytes)

	txnBytes = append(txnBytes, paddedTitleLength...)
	txnBytes = append(txnBytes, titleBytes...)
	txnBytes = append(txnBytes, claimingFeeBytes...)
	txnBytes = append(txnBytes, descriptionBytes...)

	return txnBytes, nil
}

func ChangeVmOwnerTxnFeeShareProposalTxBytes(
	title string, description string, fee_share int, nonce int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(changeVmOwnerTxnFeeShareProposalType, nonce)

	titleBytes := []byte(title)
	titleLength := len(titleBytes)
	titleLengthBytes := DecToBytes(titleLength, 4)
	descriptionBytes := []byte(description)
	feeShareBytes := DecToBytes(fee_share, 4)

	paddedTitleLength := make([]byte, 4)
	copy(paddedTitleLength[4-len(titleLengthBytes):], titleLengthBytes)
	paddedFeeShare := make([]byte, 4)
	copy(paddedFeeShare[4-len(feeShareBytes):], feeShareBytes)

	txnBytes = append(txnBytes, paddedTitleLength...)
	txnBytes = append(txnBytes, titleBytes...)
	txnBytes = append(txnBytes, paddedFeeShare...)
	txnBytes = append(txnBytes, descriptionBytes...)

	return txnBytes, nil
}

func OtherProposalTxBytes(
	title string, description string, nonce int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(otherProposalType, nonce)

	titleBytes := []byte(title)
	titleLength := len(titleBytes)
	titleLengthBytes := DecToBytes(titleLength, 4)
	descriptionBytes := []byte(description)

	paddedTitleLength := make([]byte, 4)
	copy(paddedTitleLength[4-len(titleLengthBytes):], titleLengthBytes)

	txnBytes = append(txnBytes, paddedTitleLength...)
	txnBytes = append(txnBytes, titleBytes...)
	txnBytes = append(txnBytes, descriptionBytes...)

	return txnBytes, nil
}

func VoteOnProposalTxBytes(
	proposal_hash string, vote int, nonce int,
) ([]byte, error) {
	txnBytes, _ := txnBaseBytes(voteOnProposalType, nonce)

	proposalHashBytes, err := hex.DecodeString(proposal_hash[2:])
	if err != nil {
		return nil, err
	}
	voteBytes := DecToBytes(vote, 1)

	paddedProposalHash := make([]byte, 32)
	copy(paddedProposalHash[32-len(proposalHashBytes):], proposalHashBytes)
	paddedVote := make([]byte, 1)
	copy(paddedVote[1-len(voteBytes):], voteBytes)

	txnBytes = append(txnBytes, paddedProposalHash...)
	txnBytes = append(txnBytes, paddedVote...)

	return txnBytes, nil
}

// Falcon transaction bytes
func FalconSetPublicKeyBytes(
	publicKey []byte, nonce int, address []byte, feePerByte int,
) ([]byte, error) {
	txnBytes, _ := falconTxnBaseBytes(falconSetPublicKeyType, nonce, address, feePerByte)

	publicKeyLenBytes := DecToBytes(len(publicKey), 2)

	txnBytes = append(txnBytes, publicKeyLenBytes...)
	txnBytes = append(txnBytes, publicKey...)

	return txnBytes, nil
}

func FalconJoinAsValidatorBytes(ip string, nonce int, address []byte, feePerByte int) ([]byte, error) {
	txnBytes, _ := falconTxnBaseBytes(falconJoinAsValidatorType, nonce, address, feePerByte)

	ipBytes := []byte(ip)
	ipLenBytes := DecToBytes(len(ipBytes), 2)

	txnBytes = append(txnBytes, ipLenBytes...)
	txnBytes = append(txnBytes, ipBytes...)
	return txnBytes, nil
}

func FalconDelegateTxBytes(to string, amount int, nonce int, address []byte, feePerByte int) ([]byte, error) {
	txnBytes, _ := falconTxnBaseBytes(falconDelegateType, nonce, address, feePerByte)

	recipientBytes, _ := hex.DecodeString(to[2:])
	amountBytes := decToBytes64(int64(amount), 8)

	paddedRecipient := make([]byte, 20)
	copy(paddedRecipient[20-len(recipientBytes):], recipientBytes)

	txnBytes = append(txnBytes, paddedRecipient...)
	txnBytes = append(txnBytes, amountBytes...)
	return txnBytes, nil
}

func FalconChangeIpBytes(newIp string, nonce int, address []byte, feePerByte int) ([]byte, error) {
	txnBytes, _ := falconTxnBaseBytes(falconChangeIpType, nonce, address, feePerByte)

	newIpBytes := []byte(newIp)
	newIpLenBytes := DecToBytes(len(newIpBytes), 2)

	txnBytes = append(txnBytes, newIpLenBytes...)
	txnBytes = append(txnBytes, newIpBytes...)
	return txnBytes, nil
}

func FalconClaimActiveNodeSpotBytes(nonce int, address []byte, feePerByte int) ([]byte, error) {
	txnBytes, _ := falconTxnBaseBytes(falconClaimActiveNodeSpotType, nonce, address, feePerByte)

	return txnBytes, nil
}

func FalconTransferTxBytes(
	amount int, recipient string, nonce int, address []byte, feePerByte int,
) ([]byte, error) {
	txnBytes, _ := falconTxnBaseBytes(falconTransferType, nonce, address, feePerByte)

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

func FalconVmDataBytes(
	vmId int, data []byte, nonce int, address []byte, feePerByte int,
) ([]byte, error) {
	txnBytes, _ := falconTxnBaseBytes(falconVmDataType, nonce, address, feePerByte)

	vmIdBytes := decToBytes64(int64(vmId), 8)
	dataLengthBytes := DecToBytes(int(len(data)), 4)

	txnBytes = append(txnBytes, vmIdBytes...)
	txnBytes = append(txnBytes, dataLengthBytes...)
	txnBytes = append(txnBytes, data...)

	return txnBytes, nil
}
