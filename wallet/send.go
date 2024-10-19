package wallet

import (
	"log"
	"github.com/pwrlabs/pwrgo/encode"
)

func (w *PWRWallet) TransferPWR(to string, amount int) (BroadcastResponse) {
    if len(to) != 42 {
        log.Fatal("Invalid address: ", to)
    }
    nonce := w.GetNonce()
    if nonce < 0 {
        log.Fatal("Nonce cannot be negative: ", nonce)
	}

    var buffer []byte
    buffer,err := encode.TransferTxBytes(amount, to, nonce)

    if err != nil {
        log.Fatal("Failed to get tx bytes: ", err.Error())
    }

	return SignAndBroadcast(buffer, w)
}

func (w *PWRWallet) ClaimVMId(vmId int) (BroadcastResponse) {
    nonce := w.GetNonce()
	if nonce < 0 {
        log.Fatal("Nonce cannot be negative: ", nonce)
	}

    var buffer []byte
    buffer, err := encode.ClaimVMIdBytes(vmId, nonce)
	if err != nil {
        log.Fatal("Failed to get claimVMIdBytes: ", err.Error())
    }

	return SignAndBroadcast(buffer, w)
}

func (w *PWRWallet) Join(ipAddress string) (BroadcastResponse) {
    nonce := w.GetNonce()
	if nonce < 0 {
        log.Fatal("Nonce cannot be negative: ", nonce)
	}

    var buffer []byte
    buffer, err := encode.JoinBytes(ipAddress, nonce)
	if err != nil {
        log.Fatal("Failed to get joinBytes: ", err.Error())
    }

	return SignAndBroadcast(buffer, w)
}


func (w *PWRWallet) ValidatorRemove(validatorAddress string) (BroadcastResponse) {
    nonce := w.GetNonce()
	if nonce < 0 {
        log.Fatal("Nonce cannot be negative: ", nonce)
	}

    var buffer []byte
    buffer, err := encode.ValidatorRemoveBytes(validatorAddress, nonce)
	if err != nil {
        log.Fatal("Failed to get validatorRemoveBytes: ", err.Error())
    }

	return SignAndBroadcast(buffer, w)
}

func (w *PWRWallet) ClaimActiveNodeSpot() (BroadcastResponse) {
    nonce := w.GetNonce()
	if nonce < 0 {
        log.Fatal("Nonce cannot be negative: ", nonce)
	}

    var buffer []byte
    buffer, err := encode.ClaimActiveNodeSpotBytes(nonce)
	if err != nil {
        log.Fatal("Failed to get claimActiveNodeSpotBytes: ", err.Error())
    }
	
	return SignAndBroadcast(buffer, w)
}

func (w *PWRWallet) Delegate(validator string, amount int) (BroadcastResponse) {
    nonce := w.GetNonce()
	if nonce < 0 {
        log.Fatal("Nonce cannot be negative: ", nonce)
    }

	var buffer []byte
    buffer, err := encode.DelegateTxBytes(validator, amount, nonce)
    if err != nil {
        log.Fatal("Failed to get DelegateTx bytes: ", err.Error())
    }

	return SignAndBroadcast(buffer, w)
}

func (w *PWRWallet) Withdraw(validator string, amount int) (BroadcastResponse) {
    nonce := w.GetNonce()
	if nonce < 0 {
        log.Fatal("Nonce cannot be negative: ", nonce)
    }

	var buffer []byte
    buffer, err := encode.WithdrawTxBytes(validator, amount, nonce)
    if err != nil {
        log.Fatal("Failed to get withdrawTx bytes: ", err.Error())
    }

	return SignAndBroadcast(buffer, w)
}


func (w *PWRWallet) SetGuardian(guardian string, expiration int) (BroadcastResponse) {
    nonce := w.GetNonce()
	if nonce < 0 {
        log.Fatal("Nonce cannot be negative: ", nonce)
    }

	var buffer []byte
    buffer, err := encode.SetGuardianTxBytes(guardian, expiration, nonce)
    if err != nil {
        log.Fatal("Failed to get setGuardian bytes: ", err.Error())
    }

	return SignAndBroadcast(buffer, w)
}

func (w *PWRWallet) RemoveGuardian() BroadcastResponse {
    nonce := w.GetNonce()
	if nonce < 0 {
        log.Fatal("Nonce cannot be negative: ", nonce)
    }

	var buffer []byte
    buffer, err := encode.RemoveGuardianTxBytes(nonce)
    if err != nil {
        log.Fatal("Failed to get txnBaseBytes: ", err.Error())
    }

	return SignAndBroadcast(buffer, w)
}

func (w *PWRWallet) SetConduits(vmId int, conduits []string) (BroadcastResponse) {
    nonce := w.GetNonce()
    if nonce < 0 {
        log.Fatal("Nonce cannot be negative: ", nonce)
    }

    var buffer []byte
    buffer, err := encode.SetConduitTxBytes(vmId, conduits, nonce)
    if err != nil {
        log.Fatal("Failed to get vm data bytes: ", err.Error())
    }

	return SignAndBroadcast(buffer, w)
}

func (w *PWRWallet) AddConduitTransaction(vmId int, transaction []byte) (BroadcastResponse) {
    nonce := w.GetNonce()
    if nonce < 0 {
        log.Fatal("Nonce cannot be negative: ", nonce)
    }

    var buffer []byte
    buffer, err := encode.AddConduitTxBytes(vmId, transaction, nonce)
    if err != nil {
        log.Fatal("Failed to get vm data bytes: ", err.Error())
    }

	return SignAndBroadcast(buffer, w)
}

func (w *PWRWallet) SendGuardianApprovalTransaction(tx []byte) (BroadcastResponse) {
    nonce := w.GetNonce()
    if nonce < 0 {
        log.Fatal("Nonce cannot be negative: ", nonce)
    }

    var buffer []byte
    buffer, err := encode.GuardianApprovalTransaction(tx, nonce)
    if err != nil {
        log.Fatal("Failed to get guardianWrappedTxBytes: ", err.Error())
    }

	return SignAndBroadcast(buffer, w)
}

func (w *PWRWallet) SendVMData(vmId int, data []byte) (BroadcastResponse) {
    nonce := w.GetNonce()
    if nonce < 0 {
        log.Fatal("Nonce cannot be negative: ", nonce)
    }

    var buffer []byte
    buffer, err := encode.VmDataBytes(vmId, data, nonce)
    if err != nil {
        log.Fatal("Failed to get vm data bytes: ", err.Error())
    }
    
	return SignAndBroadcast(buffer, w)
}

func (w *PWRWallet) SendPayableVMData(vmId int, amount int, data []byte) (BroadcastResponse) {
    nonce := w.GetNonce()
    if nonce < 0 {
        log.Fatal("Nonce cannot be negative: ", nonce)
    }

    var buffer []byte
    buffer, err := encode.PayableVmDataBytes(vmId, data, amount, nonce)
    if err != nil {
        log.Fatal("Failed to get vm data bytes: ", err.Error())
    }
    
	return SignAndBroadcast(buffer, w)
}

func (w *PWRWallet) MoveStake(shares int, from_validator string, to_validator string) (BroadcastResponse) {
    nonce := w.GetNonce()
    if nonce < 0 {
        log.Fatal("Nonce cannot be negative: ", nonce)
    }

    var buffer []byte
    buffer, err := encode.MoveStakeTxBytes(shares, from_validator, to_validator, nonce)
    if err != nil {
        log.Fatal("Failed to get vm data bytes: ", err.Error())
    }
    
	return SignAndBroadcast(buffer, w)
}

func (w *PWRWallet) ChangeEarlyWithdrawPenaltyProposal(
    withdraw_penalty int, withdraw_penalty_time int, title string, description string,
) (BroadcastResponse) {
    nonce := w.GetNonce()
    if nonce < 0 {
        log.Fatal("Nonce cannot be negative: ", nonce)
    }

    var buffer []byte
    buffer, err := encode.ChangeEarlyWithdrawPenaltyProposalTxBytes(
        title, description, withdraw_penalty_time, withdraw_penalty, nonce,
    )
    if err != nil {
        log.Fatal("Failed to get vm data bytes: ", err.Error())
    }
    
	return SignAndBroadcast(buffer, w)
}

func (w *PWRWallet) ChangeFeePerByteProposal(
    fee_per_byte int, title string, description string,
) (BroadcastResponse) {
    nonce := w.GetNonce()
    if nonce < 0 {
        log.Fatal("Nonce cannot be negative: ", nonce)
    }

    var buffer []byte
    buffer, err := encode.ChangeFeePerByteProposalTxBytes(
        title, description, fee_per_byte, nonce,
    )
    if err != nil {
        log.Fatal("Failed to get vm data bytes: ", err.Error())
    }
    
	return SignAndBroadcast(buffer, w)
}

func (w *PWRWallet) ChangeMaxBlockSizeProposal(
    max_block_size int, title string, description string,
) (BroadcastResponse) {
    nonce := w.GetNonce()
    if nonce < 0 {
        log.Fatal("Nonce cannot be negative: ", nonce)
    }

    var buffer []byte
    buffer, err := encode.ChangeMaxBlockSizeProposalTxBytes(
        title, description, max_block_size, nonce,
    )
    if err != nil {
        log.Fatal("Failed to get vm data bytes: ", err.Error())
    }
    
	return SignAndBroadcast(buffer, w)
}

func (w *PWRWallet) ChangeMaxTxnSizeProposal(
    max_txn_size int, title string, description string,
) (BroadcastResponse) {
    nonce := w.GetNonce()
    if nonce < 0 {
        log.Fatal("Nonce cannot be negative: ", nonce)
    }

    var buffer []byte
    buffer, err := encode.ChangeMaxTxnSizeProposalTxBytes(
        title, description, max_txn_size, nonce,
    )
    if err != nil {
        log.Fatal("Failed to get vm data bytes: ", err.Error())
    }
    
	return SignAndBroadcast(buffer, w)
}

func (w *PWRWallet) ChangeOverallBurnPercentageProposal(
    burn_percentage int, title string, description string,
) (BroadcastResponse) {
    nonce := w.GetNonce()
    if nonce < 0 {
        log.Fatal("Nonce cannot be negative: ", nonce)
    }

    var buffer []byte
    buffer, err := encode.ChangeOverallBurnPercentageProposalTxBytes(
        title, description, burn_percentage, nonce,
    )
    if err != nil {
        log.Fatal("Failed to get vm data bytes: ", err.Error())
    }
    
	return SignAndBroadcast(buffer, w)
}

func (w *PWRWallet) ChangeRewardPerYearProposal(
    burn_percentage int, title string, description string,
) (BroadcastResponse) {
    nonce := w.GetNonce()
    if nonce < 0 {
        log.Fatal("Nonce cannot be negative: ", nonce)
    }

    var buffer []byte
    buffer, err := encode.ChangeRewardPerYearProposalTxBytes(
        title, description, burn_percentage, nonce,
    )
    if err != nil {
        log.Fatal("Failed to get vm data bytes: ", err.Error())
    }
    
	return SignAndBroadcast(buffer, w)
}

func (w *PWRWallet) ChangeValidatorCountLimitProposal(
    validator_count_limit int, title string, description string,
) (BroadcastResponse) {
    nonce := w.GetNonce()
    if nonce < 0 {
        log.Fatal("Nonce cannot be negative: ", nonce)
    }

    var buffer []byte
    buffer, err := encode.ChangeValidatorCountLimitProposalTxBytes(
        title, description, validator_count_limit, nonce,
    )
    if err != nil {
        log.Fatal("Failed to get vm data bytes: ", err.Error())
    }
    
	return SignAndBroadcast(buffer, w)
}

func (w *PWRWallet) ChangeValidatorJoiningFeeProposal(
    joining_fee int, title string, description string,
) (BroadcastResponse) {
    nonce := w.GetNonce()
    if nonce < 0 {
        log.Fatal("Nonce cannot be negative: ", nonce)
    }

    var buffer []byte
    buffer, err := encode.ChangeValidatorJoiningFeeProposalTxBytes(
        title, description, joining_fee, nonce,
    )
    if err != nil {
        log.Fatal("Failed to get vm data bytes: ", err.Error())
    }
    
	return SignAndBroadcast(buffer, w)
}

func (w *PWRWallet) ChangeVmIdClaimingFeeProposal(
    claiming_fee int, title string, description string,
) (BroadcastResponse) {
    nonce := w.GetNonce()
    if nonce < 0 {
        log.Fatal("Nonce cannot be negative: ", nonce)
    }

    var buffer []byte
    buffer, err := encode.ChangeVmIdClaimingFeeProposalTxBytes(
        title, description, claiming_fee, nonce,
    )
    if err != nil {
        log.Fatal("Failed to get vm data bytes: ", err.Error())
    }
    
	return SignAndBroadcast(buffer, w)
}

func (w *PWRWallet) ChangeVmOwnerTxnFeeShareProposal(
    fee_share int, title string, description string,
) (BroadcastResponse) {
    nonce := w.GetNonce()
    if nonce < 0 {
        log.Fatal("Nonce cannot be negative: ", nonce)
    }

    var buffer []byte
    buffer, err := encode.ChangeVmOwnerTxnFeeShareProposalTxBytes(
        title, description, fee_share, nonce,
    )
    if err != nil {
        log.Fatal("Failed to get vm data bytes: ", err.Error())
    }
    
	return SignAndBroadcast(buffer, w)
}

func (w *PWRWallet) OtherProposal(
    title string, description string,
) (BroadcastResponse) {
    nonce := w.GetNonce()
    if nonce < 0 {
        log.Fatal("Nonce cannot be negative: ", nonce)
    }

    var buffer []byte
    buffer, err := encode.OtherProposalTxBytes(
        title, description, nonce,
    )
    if err != nil {
        log.Fatal("Failed to get vm data bytes: ", err.Error())
    }
    
	return SignAndBroadcast(buffer, w)
}

func (w *PWRWallet) VoteOnProposal(
    proposal_hash string, vote int,
) (BroadcastResponse) {
    nonce := w.GetNonce()
    if nonce < 0 {
        log.Fatal("Nonce cannot be negative: ", nonce)
    }

    var buffer []byte
    buffer, err := encode.VoteOnProposalTxBytes(
        proposal_hash, vote, nonce,
    )
    if err != nil {
        log.Fatal("Failed to get vm data bytes: ", err.Error())
    }
    
	return SignAndBroadcast(buffer, w)
}

