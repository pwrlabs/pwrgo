package rpc

type Transaction struct {
	ActionFee          int   	 `json:"actionFee"`
	Fee				   int 	 	 `json:"fee"`
	Nonce  	   	       int   	 `json:"nonce"`
	PositionInTheBlock int   	 `json:"positionInTheBlock"`
	Size 			   int   	 `json:"size"`
    Timestamp          int 	 	 `json:"timestamp"`
	BlockNumber		   int		 `json:"blockNumber"`
	Value 			   int		 `json:"value"`
	Success 		   bool	  	 `json:"success"`
	Hash 			   string    `json:"hash"`
	Type			   string 	 `json:"type"`
	Sender			   string 	 `json:"sender"`
	Receiver           string 	 `json:"receiver"`
	Data               string 	 `json:"data"`
}

type Block struct {
	BlockNumber 	 int		   `json:"blockNumber"`
	BlockReward      int 		   `json:"blockReward"`
	TransactionCount int		   `json:"transactionCount"`
	BlockSize 		 int 		   `json:"blockSize"`
	Timestamp 		 int		   `json:"timestamp"`
	Success 		 bool 		   `json:"success"`
	BlockHash		 string 	   `json:"blockHash"`
	BlockSubmitter 	 string 	   `json:"blockSubmitter"`
	Transactions     []Transaction `json:"transactions"`
}

type RPCResponse struct {
	Message                    string                       `json:"message,omitempty"`
	Nonce                      int                          `json:"nonce,omitempty"`
	Balance                    int                          `json:"balance,omitempty"`
	BlocksCount                int                          `json:"blocksCount,omitempty"`
	ValidatorsCount            int                          `json:"validatorsCount,omitempty"`
	ECDSAVerificationFee       int                          `json:"ecdsaVerificationFee,omitempty"`
	BurnPercentage             int                          `json:"burnPercentage,omitempty"`
	TotalVotingPower           int                          `json:"totalVotingPower,omitempty"`
	PwrRewardsPerYear          int                          `json:"pwrRewardsPerYear,omitempty"`
	WithdrawalLockTime         int                          `json:"withdrawalLockTime,omitempty"`
	MaxBlockSize               int                          `json:"maxBlockSize,omitempty"`
	MaxTransactionSize         int                          `json:"maxTransactionSize,omitempty"`
	BlockTimestamp             int                          `json:"blockTimestamp,omitempty"`
	ProposalFee                int                          `json:"proposalFee,omitempty"`
	ValidatorCountLimit        int                          `json:"validatorCountLimit,omitempty"`
	ProposalValidityTime       int                          `json:"proposalValidityTime,omitempty"`
	ValidatorSlashingFee       int                          `json:"validatorSlashingFee,omitempty"`
	ValidatorJoiningFee        int                          `json:"validatorJoiningFee,omitempty"`
	ValidatorOperationalFee    int                          `json:"validatorOperationalFee,omitempty"`
	MinimumDelegatingAmount    int                          `json:"minimumDelegatingAmount,omitempty"`
	DelegatorsCount            int                          `json:"delegatorsCount,omitempty"`
	DelegatedPWR               int                          `json:"delegatedPWR,omitempty"`
	SharesOfDelegator          int                          `json:"shares,omitempty"`
	VmIdClaimingFee            int                          `json:"vmIdClaimingFee,omitempty"`
	VmOwnerTransactionFeeShare int                          `json:"vmOwnerTransactionFeeShare,omitempty"`
	OwnerOfVmIds               string                       `json:"owner,omitempty"`
	GuardianOfAddress          string                       `json:"guardian,omitempty"`
	ShareValue                 float32                      `json:"shareValue,omitempty"`
	Delegators                 map[string]int64             `json:"delegators,omitempty"`
	AllEarlyWithdrawPenalties  map[string]string            `json:"allEarlyWithdrawPenalties,omitempty"`
	Block                      Block                        `json:"block,omitempty"`
	Validator                  Validator                    `json:"validator,omitempty"`
	Validators                 []Validator                  `json:"validators,omitempty"`
	ConduitsOfVm               []Validator                  `json:"conduits,omitempty"`
	VMDataTransaction          []VMDataTransaction          `json:"transactions,omitempty"`
	MaxGuardianTime            int                          `json:"maxGuardianTime,omitempty"`
	ActiveVotingPower          int                          `json:"activeVotingPower,omitempty"`
	FeePerByte                 int                          `json:"feePerByte,omitempty"`
	Success                    bool                         `json:"success"`
	BlockNumber                int                          `json:"blockNumber"`
	Transaction                Transaction                  `json:"transaction"`
}

type Delegator struct {
    Address string    `json:"address"`
    Shares  int64  	  `json:"shares"`
}

type VMDataTransaction struct {
    Transaction
	VmId int   `json:"vmId"`
}

type Validator struct {
    Address         string    `json:"address"`
    IP              string    `json:"ip"`
    IsBadActor      bool  	  `json:"badActor"`
    VotingPower     int 	  `json:"votingPower"`
    Shares          int  	  `json:"totalShares"`
    DelegatorsCount int  	  `json:"delegatorsCount"`
    Status          string    `json:"status"`
}

type Penalty struct {
    WithdrawTime int64 	  `json:"withdrawTime"`
    Penalty      string   `json:"penalty"`
}
