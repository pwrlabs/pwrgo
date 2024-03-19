package pwrgo

import (
   "encoding/json"
   "log"
   "strconv"
   "crypto/ecdsa"
   "github.com/ethereum/go-ethereum/common/hexutil"
   "github.com/ethereum/go-ethereum/crypto"
   "math/big"
)

var ReturnBlockNumberOnTx = false

type Transaction struct {
	PositionInTheBlock int `json:"positionInTheBlock"`
	NonceOrValidationHash string `json:"nonceOrValidationHash"`
	Size int `json:"size"`
	RawTxn string `json:"rawTxn"`
	Data string `json:"data"`
	VmId int `json:"vmId"`
	Fee int `json:"fee"`
	From string `json:"from"`
	To string `json:"to"`
	TxnFee int `json:"txnFee"`
	Type string `json:"type"`
	Hash string `json:"hash"`
}

type Block struct {
	BlockHash string `json:"blockHash"`
	Success bool `json:"success"`
	BlockNumber int `json:"blockNumber"`
	BlockReward int `json:"blockReward"`
	TransactionCount int `json:"transactionCount"`
	Transactions   []Transaction `json:"transactions"`
	BlockSubmitter string `json:"blockSubmitter"`
	BlockSize int `json:"blockSize"`
	Timestamp int `json:"timestamp"`
}

type RPCResponse struct {
  Message string `json:"message,omitempty"`
  Nonce int `json:"nonce,omitempty"`
  Balance int `json:"balance,omitempty"`
  BlocksCount int `json:"blocksCount,omitempty"`
  ValidatorsCount int `json:"validatorsCount,omitempty"`
  Block Block `json:"block,omitempty"`
  Success bool
  TxHash string
  BlockNumber int
  Error string
}

func SetRpcNodeUrl(url string) {
   RPC_ENDPOINT = url;
}

func parseRPCResponse(responseStr string) (response RPCResponse) {
    err := json.Unmarshal([]byte(responseStr), &response)
    if err != nil {
        log.Fatalf("Error unmarshaling %s", err)
    }
    return
}

func NonceOfUser(address string) (int) {
    var response = get(RPC_ENDPOINT + "/nonceOfUser/?userAddress=" + address)
    var resp = parseRPCResponse(response)
    return resp.Nonce
}

func BalanceOf(address string) (int) {
    var response = get(RPC_ENDPOINT + "/balanceOf/?userAddress=" + address)
    var resp = parseRPCResponse(response)
    return resp.Balance
}

func BlocksCount() (int) {
    var response = get(RPC_ENDPOINT + "/blocksCount/")
    var resp = parseRPCResponse(response)
    return resp.BlocksCount
}

func ValidatorsCount() (int) {
    var response = get(RPC_ENDPOINT + "/totalValidatorsCount/")
    var resp = parseRPCResponse(response)
    return resp.ValidatorsCount
}

func GetBlock(blockNumber int) (Block) {
    var blockNumberStr = strconv.Itoa(blockNumber)
    var response = get(RPC_ENDPOINT + "/block/?blockNumber=" + blockNumberStr)
	var resp = parseRPCResponse(response)
    return resp.Block
}

func TransferPWR(to string, amount string, nonce int, privateKey *ecdsa.PrivateKey) (RPCResponse) {
    if len(to) != 42 {
        log.Fatalf("Invalid address ", to)
    }
    if nonce < 0 {
        log.Fatalf("Nonce cannot be negative ", nonce)
    }

    amt := new(big.Int)
    amt.SetString(amount, 10)
    var buffer []byte
    buffer,err := transferTxBytes(nonce, amt, to)

    if err != nil {
        log.Fatalf("Failed to get tx bytes ", err.Error())
    }

    signature, err := signMessage(buffer, privateKey)
    if err != nil {
        log.Fatalf("Failed to sign message ", err.Error())
    }

	var blockNumber = 0
	if ReturnBlockNumberOnTx {
		blockNumber = BlocksCount()
	}

    finalTxn := append(buffer, signature...)
    var transferTx = hexutil.Encode(finalTxn)
    var transferTxn = `{"txn":"` + transferTx[2:] + `"}`
    var result = post(RPC_ENDPOINT + "/broadcast/", transferTxn)

	hash := crypto.Keccak256Hash(finalTxn)

	transferResponse := parseRPCResponse(result)
	
	if transferResponse.Message == "Txn broadcasted to validator nodes" {
		transferResponse.Success = true
	} else {
		transferResponse.Success = false
		transferResponse.Error = transferResponse.Message
	}

	transferResponse.TxHash = hash.Hex()
	transferResponse.BlockNumber = blockNumber
    return transferResponse
}

func ClaimVMId(vmId int64, nonce int, privateKey *ecdsa.PrivateKey) (RPCResponse) {
	if nonce < 0 {
        log.Fatalf("Nonce cannot be negative ", nonce)
	}

    var buffer []byte
    buffer, err := claimVMIdBytes(vmId, nonce)
	if err != nil {
        log.Fatalf("Failed to get vm data bytes ", err.Error())
    }

	signature, err := signMessage(buffer, privateKey)
    if err != nil {
        log.Fatalf("Failed to sign message ", err.Error())
    }

	var blockNumber = 0
	if ReturnBlockNumberOnTx {
		blockNumber = BlocksCount()
	}

    finalTxn := append(buffer, signature...)
    var dataTx = hexutil.Encode(finalTxn)
    var dataTxn = `{"txn":"` + dataTx[2:] + `"}`
    var result = post(RPC_ENDPOINT + "/broadcast/", dataTxn)
	hash := crypto.Keccak256Hash(finalTxn)

	claimTxResponse := parseRPCResponse(result)

	if claimTxResponse.Message == "Txn broadcasted to validator nodes" {
		claimTxResponse.Success = true
	} else {
		claimTxResponse.Success = false
		claimTxResponse.Error = claimTxResponse.Message
	}
	
	claimTxResponse.TxHash = hash.Hex()
	claimTxResponse.BlockNumber = blockNumber
    return claimTxResponse
}

func GetTransferTxBytes(to string, amount string, nonce int) ([]byte) {
	amt := new(big.Int)
    amt.SetString(amount, 10)
    var buffer []byte
    buffer,_ = transferTxBytes(nonce, amt, to)
	return buffer
}

func Delegate(to string, amount string, nonce int, privateKey *ecdsa.PrivateKey) (RPCResponse) {
	if nonce < 0 {
        log.Fatalf("Nonce cannot be negative ", nonce)
    }

    amt := new(big.Int)
    amt.SetString(amount, 10)
	var buffer []byte
    buffer, err := delegateTxBytes(to, amt, nonce)
    if err != nil {
        log.Fatalf("Failed to get DelegateTx bytes ", err.Error())
    }

	signature, err := signMessage(buffer, privateKey)
    if err != nil {
        log.Fatalf("Failed to sign message ", err.Error())
    }

	var blockNumber = 0
	if ReturnBlockNumberOnTx {
		blockNumber = BlocksCount()
	}

	finalTxn := append(buffer, signature...)
    var dataTx = hexutil.Encode(finalTxn)
    var dataTxn = `{"txn":"` + dataTx[2:] + `"}`
    var result = post(RPC_ENDPOINT + "/broadcast/", dataTxn)
	hash := crypto.Keccak256Hash(finalTxn)

	delegateTxResponse := parseRPCResponse(result)

	if delegateTxResponse.Message == "Txn broadcasted to validator nodes" {
		delegateTxResponse.Success = true
	} else {
		delegateTxResponse.Success = false
		delegateTxResponse.Error = delegateTxResponse.Message
	}
	
	delegateTxResponse.TxHash = hash.Hex()
	delegateTxResponse.BlockNumber = blockNumber
    return delegateTxResponse
}

func Withdraw(from string, amount string, nonce int, privateKey *ecdsa.PrivateKey) (RPCResponse) {
	if nonce < 0 {
        log.Fatalf("Nonce cannot be negative ", nonce)
    }

    amt := new(big.Int)
    amt.SetString(amount, 10)

	var buffer []byte
    buffer, err := withdrawTxBytes(from, amt, nonce)
    if err != nil {
        log.Fatalf("Failed to get withdrawTx bytes ", err.Error())
    }

	signature, err := signMessage(buffer, privateKey)
    if err != nil {
        log.Fatalf("Failed to sign message ", err.Error())
    }

	var blockNumber = 0
	if ReturnBlockNumberOnTx {
		blockNumber = BlocksCount()
	}

	finalTxn := append(buffer, signature...)
    var dataTx = hexutil.Encode(finalTxn)
    var dataTxn = `{"txn":"` + dataTx[2:] + `"}`
    var result = post(RPC_ENDPOINT + "/broadcast/", dataTxn)
	hash := crypto.Keccak256Hash(finalTxn)

	withdrawTxResponse := parseRPCResponse(result)

	if withdrawTxResponse.Message == "Txn broadcasted to validator nodes" {
		withdrawTxResponse.Success = true
	} else {
		withdrawTxResponse.Success = false
		withdrawTxResponse.Error = withdrawTxResponse.Message
	}
	
	withdrawTxResponse.TxHash = hash.Hex()
	withdrawTxResponse.BlockNumber = blockNumber
    return withdrawTxResponse
}


func SetGuardian(guardian string, expiration string, nonce int, privateKey *ecdsa.PrivateKey) (RPCResponse) {
	if nonce < 0 {
        log.Fatalf("Nonce cannot be negative ", nonce)
    }

    exp := new(big.Int)
    exp.SetString(expiration, 10)

	var buffer []byte
    buffer, err := setGuardianTxBytes(guardian, exp, nonce)
    if err != nil {
        log.Fatalf("Failed to get setGuardian bytes ", err.Error())
    }

	signature, err := signMessage(buffer, privateKey)
    if err != nil {
        log.Fatalf("Failed to sign message ", err.Error())
    }

	var blockNumber = 0
	if ReturnBlockNumberOnTx {
		blockNumber = BlocksCount()
	}

	finalTxn := append(buffer, signature...)
    var dataTx = hexutil.Encode(finalTxn)
    var dataTxn = `{"txn":"` + dataTx[2:] + `"}`
    var result = post(RPC_ENDPOINT + "/broadcast/", dataTxn)
	hash := crypto.Keccak256Hash(finalTxn)

	setGuardianTxResponse := parseRPCResponse(result)

	if setGuardianTxResponse.Message == "Txn broadcasted to validator nodes" {
		setGuardianTxResponse.Success = true
	} else {
		setGuardianTxResponse.Success = false
		setGuardianTxResponse.Error = setGuardianTxResponse.Message
	}
	
	setGuardianTxResponse.TxHash = hash.Hex()
	setGuardianTxResponse.BlockNumber = blockNumber
    return setGuardianTxResponse
}

func RemoveGuardian(nonce int, privateKey *ecdsa.PrivateKey) (RPCResponse) {
	if nonce < 0 {
        log.Fatalf("Nonce cannot be negative ", nonce)
    }

	var buffer []byte
    buffer, err := txnBaseBytes(9, nonce)
    if err != nil {
        log.Fatalf("Failed to get txnBaseBytes ", err.Error())
    }

	signature, err := signMessage(buffer, privateKey)
    if err != nil {
        log.Fatalf("Failed to sign message ", err.Error())
    }

	var blockNumber = 0
	if ReturnBlockNumberOnTx {
		blockNumber = BlocksCount()
	}

	finalTxn := append(buffer, signature...)
    var dataTx = hexutil.Encode(finalTxn)
    var dataTxn = `{"txn":"` + dataTx[2:] + `"}`
    var result = post(RPC_ENDPOINT + "/broadcast/", dataTxn)
	hash := crypto.Keccak256Hash(finalTxn)

	removeGuardianTxResponse := parseRPCResponse(result)

	if removeGuardianTxResponse.Message == "Txn broadcasted to validator nodes" {
		removeGuardianTxResponse.Success = true
	} else {
		removeGuardianTxResponse.Success = false
		removeGuardianTxResponse.Error = removeGuardianTxResponse.Message
	}
	
	removeGuardianTxResponse.TxHash = hash.Hex()
	removeGuardianTxResponse.BlockNumber = blockNumber
    return removeGuardianTxResponse
}

// TO-DO: Fix/test this function
func SendConduitTx(vmId int64, tx []byte, nonce int, privateKey *ecdsa.PrivateKey) (RPCResponse) {
    if nonce < 0 {
        log.Fatalf("Nonce cannot be negative ", nonce)
    }

    var buffer []byte
    buffer, err := sendConduitTxBytes(vmId, nonce, tx)
    if err != nil {
        log.Fatalf("Failed to get vm data bytes ", err.Error())
    }

	signature, err := signMessage(buffer, privateKey)
    if err != nil {
        log.Fatalf("Failed to sign message ", err.Error())
    }

	var blockNumber = 0
	if ReturnBlockNumberOnTx {
		blockNumber = BlocksCount()
	}

    finalTxn := append(buffer, signature...)
    var dataTx = hexutil.Encode(finalTxn)
    var dataTxn = `{"txn":"` + dataTx[2:] + `"}`
    var result = post(RPC_ENDPOINT + "/broadcast/", dataTxn)
	hash := crypto.Keccak256Hash(finalTxn)

	conduitTxResponse := parseRPCResponse(result)

	if conduitTxResponse.Message == "Txn broadcasted to validator nodes" {
		conduitTxResponse.Success = true
	} else {
		conduitTxResponse.Success = false
		conduitTxResponse.Error = conduitTxResponse.Message
	}
	
	conduitTxResponse.TxHash = hash.Hex()
	conduitTxResponse.BlockNumber = blockNumber
    return conduitTxResponse
}

func SendVMDataTx(vmId int64, data []byte, nonce int, privateKey *ecdsa.PrivateKey) (RPCResponse) {
    if nonce < 0 {
        log.Fatalf("Nonce cannot be negative ", nonce)
    }

    var buffer []byte
    buffer, err := vmDataBytes(vmId, nonce, data)
    if err != nil {
        log.Fatalf("Failed to get vm data bytes ", err.Error())
    }

    signature, err := signMessage(buffer, privateKey)
    if err != nil {
        log.Fatalf("Failed to sign message ", err.Error())
    }

	var blockNumber = 0
	if ReturnBlockNumberOnTx {
		blockNumber = BlocksCount()
	}

    finalTxn := append(buffer, signature...)
    var dataTx = hexutil.Encode(finalTxn)
    var dataTxn = `{"txn":"` + dataTx[2:] + `"}`
    var result = post(RPC_ENDPOINT + "/broadcast/", dataTxn)
	hash := crypto.Keccak256Hash(finalTxn)

	vmDataTxResponse := parseRPCResponse(result)

	if vmDataTxResponse.Message == "Txn broadcasted to validator nodes" {
		vmDataTxResponse.Success = true
	} else {
		vmDataTxResponse.Success = false
		vmDataTxResponse.Error = vmDataTxResponse.Message
	}
	
	vmDataTxResponse.TxHash = hash.Hex()
	vmDataTxResponse.BlockNumber = blockNumber
    return vmDataTxResponse
}
