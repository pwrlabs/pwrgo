package rpc

import (
	"encoding/json"
    "log"
    "net/http"
    "io"
)

var rpcEndpoint = "https://pwrrpc.pwrlabs.io"

func get(url string) (response string) {
   	resp, err := http.Get(url)
   	if err != nil {
   	   log.Fatalln(err)
   	}

   	body, err := io.ReadAll(resp.Body)
	if err != nil {
   	   log.Fatalln(err)
	}
   	response = string(body)
   	return
}

func parseRPCResponse(responseStr string) (response RPCResponse) {
	err := json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		log.Fatalf("Error unmarshaling: %s", err)
	}

	if len(response.Message) != 0 {
		log.Printf("Error: %s", response.Message)
	}
	return response
}
