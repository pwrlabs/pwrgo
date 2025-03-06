package wallet

import (
	"os"
	"encoding/hex"
	"crypto/ecdsa"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pwrlabs/pwrgo/encode"
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

func SignTx(buffer []byte, account *PWRWallet) ([]byte, error) {
	signature, err := SignMessage(buffer, account)
	if err != nil {
		return nil, err
	}
	txn_bytes := append(buffer, signature...)
	return txn_bytes, nil
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

func LoadWallet(path string, password string) (*PWRWallet, error) {
	encryptedData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	privateKeyBytes, err := encode.Decrypt(encryptedData, password)
	if err != nil {
		return nil, err
	}

	privateKey := hex.EncodeToString(privateKeyBytes)
	return FromPrivateKey(privateKey), nil
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
