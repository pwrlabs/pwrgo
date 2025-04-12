package rpc

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func SetRpcNodeUrl(url string) *RPC {
	return &RPC{
		RpcEndpoint: url,
	}
}

func (r *RPC) GetRpcNodeUrl() string {
	return r.RpcEndpoint
}

func (r *RPC) BroadcastTransaction(txn_bytes []byte) BroadcastResponse {
	var transferTx = hexutil.Encode(txn_bytes)
	var transferTxn = `{"txn":"` + transferTx[2:] + `"}`
	result := post(r.GetRpcNodeUrl()+"/broadcast/", transferTxn)

	hash := crypto.Keccak256Hash(txn_bytes)
	txResponse := parseBroadcastResponse(result)

	if txResponse.Message == "Txn broadcast to validator nodes" {
		txResponse.Success = true
		txResponse.Hash = hash.Hex()
	} else {
		txResponse.Success = false
		txResponse.Error = txResponse.Message
	}

	return txResponse
}
