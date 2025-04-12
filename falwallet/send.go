package falwallet

import (
	"log"

	"github.com/pwrlabs/pwrgo/encode"
	"github.com/pwrlabs/pwrgo/rpc"
)

func (w *Falcon512Wallet) SetPublicKey(publicKey []byte, feePerByte int) (rpc.BroadcastResponse) {
    var buffer []byte
    buffer, err := encode.FalconSetPublicKeyBytes(publicKey, w.GetNonce(), w.Address, feePerByte)
    if err != nil {
        log.Fatal("Failed to get tx bytes: ", err.Error())
    }

    txn_bytes, err := w.SignTx(buffer)
    if err != nil {
        log.Fatal("Failed to sign message: ", err.Error())
    }

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *Falcon512Wallet) JoinAsValidator(ip string, feePerByte int) (rpc.BroadcastResponse) {
    response := makeSurePublicKeyIsSet(feePerByte, w)
    if response != nil && !response.Success { return *response }

    var buffer []byte
    buffer, err := encode.FalconJoinAsValidatorBytes(ip, w.GetNonce(), w.Address, feePerByte)
    if err != nil {
        log.Fatal("Failed to get tx bytes: ", err.Error())
    }

    txn_bytes, err := w.SignTx(buffer)
    if err != nil {
        log.Fatal("Failed to sign message: ", err.Error())
    }

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *Falcon512Wallet) Delegate(to string, amount int, feePerByte int) (rpc.BroadcastResponse) {
    response := makeSurePublicKeyIsSet(feePerByte, w)
    if response != nil && !response.Success { return *response }

    var buffer []byte
    buffer, err := encode.FalconDelegateTxBytes(to, amount, w.GetNonce(), w.Address, feePerByte)
    if err != nil {
        log.Fatal("Failed to get tx bytes: ", err.Error())
    }

    txn_bytes, err := w.SignTx(buffer)
    if err != nil {
        log.Fatal("Failed to sign message: ", err.Error())
    }

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *Falcon512Wallet) ChangeIp(newIp string, feePerByte int) (rpc.BroadcastResponse) {
    response := makeSurePublicKeyIsSet(feePerByte, w)
    if response != nil && !response.Success { return *response }

    var buffer []byte
    buffer, err := encode.FalconChangeIpBytes(newIp, w.GetNonce(), w.Address, feePerByte)
    if err != nil {
        log.Fatal("Failed to get tx bytes: ", err.Error())
    }

    txn_bytes, err := w.SignTx(buffer)
    if err != nil {
        log.Fatal("Failed to sign message: ", err.Error())
    }

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *Falcon512Wallet) ClaimActiveNodeSpot(feePerByte int) (rpc.BroadcastResponse) {
    response := makeSurePublicKeyIsSet(feePerByte, w)
    if response != nil && !response.Success { return *response }

    var buffer []byte
    buffer, err := encode.FalconClaimActiveNodeSpotBytes(w.GetNonce(), w.Address, feePerByte)
    if err != nil {
        log.Fatal("Failed to get tx bytes: ", err.Error())
    }

    txn_bytes, err := w.SignTx(buffer)
    if err != nil {
        log.Fatal("Failed to sign message: ", err.Error())
    }

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *Falcon512Wallet) TransferPWR(to string, amount int, feePerByte int) (rpc.BroadcastResponse) {
    if len(to) != 42 {
        log.Fatal("Invalid address: ", to)
    }

    response := makeSurePublicKeyIsSet(feePerByte, w)
    if response != nil && !response.Success { return *response }

    var buffer []byte
    buffer, err := encode.FalconTransferTxBytes(amount, to, w.GetNonce(), w.Address, feePerByte)
    if err != nil {
        log.Fatal("Failed to get tx bytes: ", err.Error())
    }

    txn_bytes, err := w.SignTx(buffer)
    if err != nil {
        log.Fatal("Failed to sign message: ", err.Error())
    }

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func (w *Falcon512Wallet) SendVMData(vmId int, data []byte, feePerByte int) (rpc.BroadcastResponse) {
    response := makeSurePublicKeyIsSet(feePerByte, w)
    if response != nil && !response.Success { return *response }

    var buffer []byte
    buffer, err := encode.FalconVmDataBytes(vmId, data, w.GetNonce(), w.Address, feePerByte)
    if err != nil {
        log.Fatal("Failed to get tx bytes: ", err.Error())
    }

    txn_bytes, err := w.SignTx(buffer)
    if err != nil {
        log.Fatal("Failed to sign message: ", err.Error())
    }

	return w.rpc.BroadcastTransaction(txn_bytes)
}

func makeSurePublicKeyIsSet(feePerByte int, w *Falcon512Wallet) *rpc.BroadcastResponse {
    nonce := w.GetNonce()

    if nonce == 0 {
        tx := w.SetPublicKey(w.PublicKey, feePerByte)
        return &tx
    }

    return nil
}
