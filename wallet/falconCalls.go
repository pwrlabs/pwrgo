package wallet

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

var falconApiUrl = "http://localhost:3000"


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

func generateKeypair(seed string) (publicKey string, secretKey string) {
	var response = get(falconApiUrl + "/generateKeypair?seed=" + hex.EncodeToString([]byte(seed)))

	var jsonResponse map[string]interface{}
	json.Unmarshal([]byte(response), &jsonResponse)

	publicKey = jsonResponse["data"].(map[string]interface{})["publicKey"].(string)
	secretKey = jsonResponse["data"].(map[string]interface{})["secretKey"].(string)

	return publicKey, secretKey
}

func generateRandomKeypair(wordCount int) (publicKey string, secretKey string, seedPhrase string) {
	var response = get(falconApiUrl + "/generateRandomKeypair?wordCount=" + fmt.Sprintf("%d", wordCount))
	
	var jsonResponse map[string]interface{}
	json.Unmarshal([]byte(response), &jsonResponse)

	publicKey = jsonResponse["data"].(map[string]interface{})["publicKey"].(string)
	secretKey = jsonResponse["data"].(map[string]interface{})["secretKey"].(string)
	seedPhrase = jsonResponse["data"].(map[string]interface{})["seedPhrase"].(string)

	return publicKey, secretKey, seedPhrase
}
