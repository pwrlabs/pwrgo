package rpc

func SetRpcNodeUrl(url string) {
	rpcEndpoint = url
}

func GetRpcNodeUrl() string {
    return rpcEndpoint
}
