package wallet

import (
	"net/http"
    "io"
    "bytes"
	"encoding/json"
	"crypto/ecdsa"
	"log"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pwrlabs/pwrgo/rpc"
)

func SignMessage(message []byte, account *PWRWallet) ([]byte, error) {
	messageHash := crypto.Keccak256(message)
	signature, err := crypto.Sign(messageHash, account.privateKey)
 
	if err != nil {
	   	return nil, err
	}
	
	if signature[64] == 0 || signature[64] ==  1 {
	  	signature[64] += 27
	} 
	
	return signature, nil
}

func SignAndBroadcast(buffer []byte, account *PWRWallet) (BroadcastResponse) {
	signature, err := SignMessage(buffer, account)
    if err != nil {
        log.Fatal("Failed to sign message: ", err.Error())
    }

    finalTxn := append(buffer, signature...)
    var transferTx = hexutil.Encode(finalTxn)
    var transferTxn = `{"txn":"` + transferTx[2:] + `"}`
    result := post(rpc.GetRpcNodeUrl() + "/broadcast/", transferTxn)

	hash := crypto.Keccak256Hash(finalTxn)
	txResponse := parseBroadcastResponse(result)
	
	if txResponse.Message == "Txn broadcast to validator nodes" {
		txResponse.Success = true
	} else {
		txResponse.Success = false
		txResponse.Error = txResponse.Message
	}

	txResponse.TxHash = hash.Hex()
    return txResponse
}

func FromPrivateKey(privateKeyStr string) *PWRWallet {
	if privateKeyStr[0:2] == "0x" {
		privateKeyStr = privateKeyStr[2:]
	}

	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		log.Fatal(err.Error())
	}

	return privateKeyToWallet(privateKey)
}

func NewWallet() *PWRWallet {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err.Error())
	}
	return privateKeyToWallet(privateKey)
}

func post(url string, jsonStr string) string {
    var jsonBytes = []byte(jsonStr)
    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))

    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)
    return string(body)
}


func privateKeyToWallet (privateKey *ecdsa.PrivateKey) *PWRWallet {
	publicKey := &privateKey.PublicKey
	publicKeyStr := hexutil.Encode(crypto.FromECDSAPub(publicKey))
	privateKeyStr := hexutil.Encode(crypto.FromECDSA(privateKey))
	address := crypto.PubkeyToAddress(*publicKey)
		
	var wallet = new(PWRWallet)
	wallet.privateKey = privateKey
	wallet.publicKey = publicKeyStr
	wallet.address = address.Hex()
	wallet.privateKeyStr = privateKeyStr
	return wallet
}

func parseBroadcastResponse(responseStr string) (response BroadcastResponse) {
    err := json.Unmarshal([]byte(responseStr), &response)
    if err != nil {
        log.Fatalf("Error unmarshaling: %s", err)
    }
    return
}
