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

func SignAndBroadcast(buffer []byte, privateKey *ecdsa.PrivateKey) (RPCResponse){
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

	txResponse := parseRPCResponse(result)
	
	if txResponse.Message == "Txn broadcasted to validator nodes" {
		txResponse.Success = true
	} else {
		txResponse.Success = false
		txResponse.Error = txResponse.Message
	}

	txResponse.TxHash = hash.Hex()
	txResponse.BlockNumber = blockNumber
    return txResponse
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

	return SignAndBroadcast(buffer, privateKey)
}

func ClaimVMId(vmId int64, nonce int, privateKey *ecdsa.PrivateKey) (RPCResponse) {
	if nonce < 0 {
        log.Fatalf("Nonce cannot be negative ", nonce)
	}

    var buffer []byte
    buffer, err := claimVMIdBytes(vmId, nonce)
	if err != nil {
        log.Fatalf("Failed to get claimVMIdBytes ", err.Error())
    }

	return SignAndBroadcast(buffer, privateKey)
}

func Join(ipAddress string, nonce int, privateKey *ecdsa.PrivateKey) (RPCResponse) {
	if nonce < 0 {
        log.Fatalf("Nonce cannot be negative ", nonce)
	}

    var buffer []byte
    buffer, err := joinBytes(ipAddress, nonce)
	if err != nil {
        log.Fatalf("Failed to get joinBytes ", err.Error())
    }

	return SignAndBroadcast(buffer, privateKey)
}


func ValidatorRemove(validatorAddress string, nonce int, privateKey *ecdsa.PrivateKey) (RPCResponse) {
	if nonce < 0 {
        log.Fatalf("Nonce cannot be negative ", nonce)
	}

    var buffer []byte
    buffer, err := validatorRemoveBytes(validatorAddress, nonce)
	if err != nil {
        log.Fatalf("Failed to get validatorRemoveBytes ", err.Error())
    }

	return SignAndBroadcast(buffer, privateKey)
}

func ClaimActiveNodeSpot(nonce int, privateKey *ecdsa.PrivateKey) (RPCResponse) {
	if nonce < 0 {
        log.Fatalf("Nonce cannot be negative ", nonce)
	}

    var buffer []byte
    buffer, err := claimActiveNodeSpotBytes(nonce)
	if err != nil {
        log.Fatalf("Failed to get claimActiveNodeSpotBytes ", err.Error())
    }
	
	return SignAndBroadcast(buffer, privateKey)
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

	return SignAndBroadcast(buffer, privateKey)
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

	return SignAndBroadcast(buffer, privateKey)
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

	return SignAndBroadcast(buffer, privateKey)
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

	return SignAndBroadcast(buffer, privateKey)
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

	return SignAndBroadcast(buffer, privateKey)
}

func SendGuardianWrappedTx(tx []byte, nonce int, privateKey *ecdsa.PrivateKey) (RPCResponse) {
    if nonce < 0 {
        log.Fatalf("Nonce cannot be negative ", nonce)
    }

    var buffer []byte
    buffer, err := guardianWrappedTxBytes(tx, nonce)
    if err != nil {
        log.Fatalf("Failed to get guardianWrappedTxBytes ", err.Error())
    }

	return SignAndBroadcast(buffer, privateKey)
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
    
	return SignAndBroadcast(buffer, privateKey)
}
