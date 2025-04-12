package rpc

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

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

func parseBroadcastResponse(responseStr string) (response BroadcastResponse) {
	err := json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		log.Fatalf("Error unmarshaling: %s", err)
	}
	return
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
