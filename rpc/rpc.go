package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

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

// BroadcastTransaction broadcasts a transaction to the network via the configured RPC node.
// This method serializes the provided transaction as a hex string and sends
// it to the RPC node for broadcasting. Upon successful broadcast, the transaction
// hash is returned. In case of any issues during broadcasting, appropriate errors
// are captured in the Response.
func (r *RPC) BroadcastTransaction(transaction []byte) BroadcastResponse {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 3 * time.Second, // 3 seconds timeout
	}

	// Prepare request payload
	hexString := hexutil.Encode(transaction)[2:] // Remove "0x" prefix
	payload := map[string]string{
		"transaction": hexString, // Keep both fields as in Java
		"txn":         hexString,
	}
	
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return BroadcastResponse{
			Success: false,
			Hash:    "",
			Error:   fmt.Sprintf("Failed to marshal JSON: %v", err),
		}
	}

	// Create request
	req, err := http.NewRequest("POST", r.GetRpcNodeUrl()+"/broadcast", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return BroadcastResponse{
			Success: false,
			Hash:    "",
			Error:   fmt.Sprintf("Failed to create request: %v", err),
		}
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return BroadcastResponse{
			Success: false,
			Hash:    "",
			Error:   fmt.Sprintf("Failed to execute request: %v", err),
		}
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return BroadcastResponse{
			Success: false,
			Hash:    "",
			Error:   fmt.Sprintf("Failed to read response body: %v", err),
		}
	}

	// Calculate transaction hash
	hash := crypto.Keccak256Hash(transaction)
	
	// Handle response based on status code
	if resp.StatusCode == 200 {
		return BroadcastResponse{
			Message: "Txn broadcast to validator nodes",
			Success: true,
			Hash:    hash.Hex(),
			Error:   "",
		}
	} else if resp.StatusCode == 400 {
		var responseObj map[string]interface{}
		if err := json.Unmarshal(body, &responseObj); err != nil {
			return BroadcastResponse{
				Success: false,
				Hash:    "",
				Error:   fmt.Sprintf("Failed to parse error response: %v", err),
			}
		}
		
		message, _ := responseObj["message"].(string)
		fmt.Printf("broadcast response: %s\n", string(body))
		
		return BroadcastResponse{
			Message: message,
			Success: false,
			Hash:    "",
			Error:   message,
		}
	} else {
		errorMsg := fmt.Sprintf("Failed with HTTP error code: %d %s", resp.StatusCode, string(body))
		fmt.Println(errorMsg)
		return BroadcastResponse{
			Success: false,
			Hash:    "",
			Error:   errorMsg,
		}
	}
}