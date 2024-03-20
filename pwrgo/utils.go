package pwrgo

import (
   "crypto/ecdsa"
   "github.com/ethereum/go-ethereum/crypto"
   "encoding/hex"
   "math/big"
)

var RPC_ENDPOINT = "https://pwrrpc.pwrlabs.io"

var TRANSFER_TYPE = 0
var JOIN_TYPE = 1
var CLAIM_ACTIVE_NODE_SPOT_TYPE = 2
var DELEGATE_TYPE = 3
var WITHDRAW_TYPE = 4
var VM_DATA_TX_TYPE = 5
var CLAIM_VM_ID_TYPE = 6
var VALIDATOR_REMOVE_TYPE = 7
var SET_GUARDIAN_TYPE = 8
var SEND_GUARDIAN_WRAPPED_TX_TYPE = 10
var SEND_CONDUIT_TX_TYPE = 11

func txnBaseBytes(txType int, nonce int) ([]byte, error){
   typeByte := decToBytes(txType, 1)
   chainByte := decToBytes(0, 1)
   nonceBytes := decToBytes(nonce, 4)

   paddedNonce := make([]byte, 4)
   copy(paddedNonce[4-len(nonceBytes):], nonceBytes)

   var txnBytes []byte
   txnBytes = append(txnBytes, typeByte...)
   txnBytes = append(txnBytes, chainByte...)
   txnBytes = append(txnBytes, paddedNonce...)

   return txnBytes, nil
}

func claimVMIdBytes(vmId int64, nonce int) ([]byte, error) {
   txnBytes, _ := txnBaseBytes(CLAIM_VM_ID_TYPE, nonce);
   vmIdBytes := decToBytes64(vmId, 8)

   txnBytes = append(txnBytes, vmIdBytes...)

   return txnBytes, nil
}

func claimActiveNodeSpotBytes(nonce int) ([]byte, error) {
   txnBytes, _ := txnBaseBytes(CLAIM_ACTIVE_NODE_SPOT_TYPE, nonce);
   return txnBytes, nil
}

func joinBytes(ip string, nonce int) ([]byte, error) {
   txnBytes, _ := txnBaseBytes(JOIN_TYPE, nonce);
   
   ipBytes := []byte(ip)
   txnBytes = append(txnBytes, ipBytes...)
   return txnBytes, nil
}


func validatorRemoveBytes(validatorAddress string, nonce int) ([]byte, error) {
   txnBytes, _ := txnBaseBytes(VALIDATOR_REMOVE_TYPE, nonce);
   
   validatorBytes, _ := hex.DecodeString(validatorAddress[2:])
   paddedValidator := make([]byte, 20)
   copy(paddedValidator[20-len(validatorBytes):], validatorBytes)

   txnBytes = append(txnBytes, paddedValidator...)
   return txnBytes, nil
}

func guardianWrappedTxBytes(txn []byte, nonce int) ([]byte, error) {
   txnBytes, _ := txnBaseBytes(SEND_GUARDIAN_WRAPPED_TX_TYPE, nonce);
   
   txnBytes = append(txnBytes, txn...)
   return txnBytes, nil
}

func delegateTxBytes(to string, amount *big.Int, nonce int) ([]byte, error) {
   txnBytes, _ := txnBaseBytes(DELEGATE_TYPE, nonce);

   amountBytes := amount.Bytes()
   paddedAmount := make([]byte, 8)
   copy(paddedAmount[8-len(amountBytes):], amountBytes)
   
   recipientBytes, _ := hex.DecodeString(to[2:])
   paddedRecipient := make([]byte, 20)
   copy(paddedRecipient[20-len(recipientBytes):], recipientBytes)
   
   txnBytes = append(txnBytes, paddedAmount...)
   txnBytes = append(txnBytes, paddedRecipient...)
   return txnBytes, nil
}

func withdrawTxBytes(from string, amount *big.Int, nonce int) ([]byte, error) {
   txnBytes, _ := txnBaseBytes(WITHDRAW_TYPE, nonce);

   amountBytes := amount.Bytes()
   paddedAmount := make([]byte, 8)
   copy(paddedAmount[8-len(amountBytes):], amountBytes)
   
   recipientBytes, _ := hex.DecodeString(from[2:])
   paddedRecipient := make([]byte, 20)
   copy(paddedRecipient[20-len(recipientBytes):], recipientBytes)
   
   txnBytes = append(txnBytes, paddedAmount...)
   txnBytes = append(txnBytes, paddedRecipient...)
   return txnBytes, nil
}

func setGuardianTxBytes(guardian string, expiration *big.Int, nonce int) ([]byte, error) {
   txnBytes, _ := txnBaseBytes(SET_GUARDIAN_TYPE, nonce);

   expirationBytes := expiration.Bytes()
   paddedExpiration := make([]byte, 8)
   copy(paddedExpiration[8-len(expirationBytes):], expirationBytes)
   
   recipientBytes, _ := hex.DecodeString(guardian[2:])
   paddedRecipient := make([]byte, 20)
   copy(paddedRecipient[20-len(recipientBytes):], recipientBytes)
   
   txnBytes = append(txnBytes, paddedExpiration...)
   txnBytes = append(txnBytes, paddedRecipient...)
   return txnBytes, nil
}

func sendConduitTxBytes(vmId int64, nonce int, txn []byte) ([]byte, error) {
   txnBytes, _ := txnBaseBytes(SEND_CONDUIT_TX_TYPE, nonce);
   vmIdBytes := decToBytes64(vmId, 8)

   txnBytes = append(txnBytes, vmIdBytes...)
   txnBytes = append(txnBytes, txn...)

   return txnBytes, nil
}

func transferTxBytes(nonce int, amount *big.Int, recipient string) ([]byte, error) { // TransferPWR()
   txnBytes,_ := txnBaseBytes(TRANSFER_TYPE, nonce);

   amountBytes := amount.Bytes()
   recipientBytes, err := hex.DecodeString(recipient[2:])
   if err != nil {
      return nil, err
   }

   paddedAmount := make([]byte, 8)
   copy(paddedAmount[8-len(amountBytes):], amountBytes)
   
   paddedRecipient := make([]byte, 20)
   copy(paddedRecipient[20-len(recipientBytes):], recipientBytes)
   
   txnBytes = append(txnBytes, paddedAmount...)
   txnBytes = append(txnBytes, paddedRecipient...)
   
   return txnBytes, nil
}

func vmDataBytes(vmId int64, nonce int, data []byte) ([]byte, error) {
   txnBytes,_ := txnBaseBytes(VM_DATA_TX_TYPE, nonce)
   vmIdBytes := decToBytes64(vmId, 8)

   txnBytes = append(txnBytes, vmIdBytes...)
   txnBytes = append(txnBytes, data...)
   
   return txnBytes, nil
}

func decToBytes(value, length int) []byte {
   result := make([]byte, length)
   for i := 0; i < length; i++ {
      result[length-1-i] = byte(value >> (8 * i))
   }
   return result
}

func decToBytes64(value int64, length int) []byte {
   result := make([]byte, length)
   for i := 0; i < length; i++ {
      result[length-1-i] = byte(value >> (8 * i))
   }
   return result
}

func signMessage(message []byte, privateKey *ecdsa.PrivateKey) ([]byte, error) {
   messageHash := crypto.Keccak256(message)
   signature, err := crypto.Sign(messageHash, privateKey)
   if err != nil {
      return nil, err
   }
   
   if signature[64] == 0 || signature[64] ==  1 {
     signature[64] += 27
   } 
   
   return signature, nil
}