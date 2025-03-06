package encode

const (
    transferType                              = 0
    joinType                                  = 1
    claimSpotType                             = 2
    delegateType                              = 3
    withdrawType                              = 4
    vmDataType                                = 5
    claimVmIdType                             = 6
    validatorRemoveType                       = 7
    setGuardianType                           = 8
    removeGuardianType                        = 9
    guardianApprovalType                      = 10
    payableVmDataType                         = 11
    setConduitsType                           = 13
    addConduitsType                           = 14
    moveStakeType                             = 16
    changeEarlyWithdrawPenaltyProposalType    = 17
    changeFeePerByteProposalType              = 18
    changeMaxBlockSizeProposalType            = 19
    changeMaxTxnSizeProposalType              = 20
    changeOverallBurnPercentageProposalType   = 21
    changeRewardPerYearProposalType           = 22
    changeValidatorCountLimitProposalType     = 23
    changeValidatorJoiningFeeProposalType     = 24
    changeVmIdClaimingFeeProposalType         = 25
    changeVmOwnerTxnFeeShareProposalType      = 26
    otherProposalType                         = 27
    voteOnProposalType                        = 28

    // Falcon transaction types
    falconSetPublicKeyType                    = 1001
    falconJoinAsValidatorType                 = 1002
    falconDelegateType                        = 1003
    falconChangeIpType                        = 1004
    falconClaimActiveNodeSpotType             = 1005
    falconTransferType                        = 1006
    falconVmDataType                          = 1007
)
