package pwrnosqldb

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// PowerKvError represents errors from PowerKv operations
type PowerKvError struct {
	Type    string
	Message string
}

func (e *PowerKvError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// NewPowerKvError creates a new PowerKvError
func NewPowerKvError(errorType, message string) *PowerKvError {
	return &PowerKvError{Type: errorType, Message: message}
}

// StoreDataRequest represents the JSON payload for storing data
type StoreDataRequest struct {
	ProjectID string `json:"projectId"`
	Secret    string `json:"secret"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}

// GetValueResponse represents the JSON response for getting data
type GetValueResponse struct {
	Value string `json:"value"`
}

// ErrorResponse represents error responses from the server
type ErrorResponse struct {
	Message string `json:"message"`
}

// PowerKv represents the basic PowerKv client
type PowerKv struct {
	client    *http.Client
	serverURL string
	projectID string
	secret    string
}

// NewPowerKv creates a new PowerKv instance
func NewPowerKv(projectID, secret string) (*PowerKv, error) {
	if strings.TrimSpace(projectID) == "" {
		return nil, NewPowerKvError("InvalidInput", "Project ID cannot be null or empty")
	}
	if strings.TrimSpace(secret) == "" {
		return nil, NewPowerKvError("InvalidInput", "Secret cannot be null or empty")
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	return &PowerKv{
		client:    client,
		serverURL: "https://pwrnosqlvida.pwrlabs.io/",
		projectID: projectID,
		secret:    secret,
	}, nil
}

// GetServerURL returns the server URL
func (p *PowerKv) GetServerURL() string {
	return p.serverURL
}

// GetProjectID returns the project ID
func (p *PowerKv) GetProjectID() string {
	return p.projectID
}

// toHexString converts bytes to hex string
func (p *PowerKv) toHexString(data []byte) string {
	return hex.EncodeToString(data)
}

// fromHexString converts hex string to bytes
func (p *PowerKv) fromHexString(hexString string) ([]byte, error) {
	// Handle both with and without 0x prefix
	if strings.HasPrefix(hexString, "0x") || strings.HasPrefix(hexString, "0X") {
		hexString = hexString[2:]
	}

	data, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, NewPowerKvError("HexDecodeError", fmt.Sprintf("Invalid hex: %v", err))
	}
	return data, nil
}

// toBytes converts various types to bytes
func (p *PowerKv) toBytes(data interface{}) ([]byte, error) {
	if data == nil {
		return nil, NewPowerKvError("InvalidInput", "Data cannot be nil")
	}

	switch v := data.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	case int:
		return []byte(strconv.Itoa(v)), nil
	case int32:
		return []byte(strconv.FormatInt(int64(v), 10)), nil
	case int64:
		return []byte(strconv.FormatInt(v, 10)), nil
	case float32:
		return []byte(strconv.FormatFloat(float64(v), 'f', -1, 32)), nil
	case float64:
		return []byte(strconv.FormatFloat(v, 'f', -1, 64)), nil
	default:
		return nil, NewPowerKvError("InvalidInput", "Data must be []byte, string, or number")
	}
}

// Put stores data with the given key
func (p *PowerKv) Put(key, data interface{}) (bool, error) {
	keyBytes, err := p.toBytes(key)
	if err != nil {
		return false, err
	}

	dataBytes, err := p.toBytes(data)
	if err != nil {
		return false, err
	}

	return p.PutBytes(keyBytes, dataBytes)
}

// PutBytes stores byte data with byte key
func (p *PowerKv) PutBytes(key, data []byte) (bool, error) {
	url := p.serverURL + "/storeData"

	payload := StoreDataRequest{
		ProjectID: p.projectID,
		Secret:    p.secret,
		Key:       p.toHexString(key),
		Value:     p.toHexString(data),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return false, NewPowerKvError("NetworkError", fmt.Sprintf("Failed to marshal JSON: %v", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return false, NewPowerKvError("NetworkError", fmt.Sprintf("Failed to create request: %v", err))
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return false, NewPowerKvError("NetworkError", fmt.Sprintf("Request failed: %v", err))
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, NewPowerKvError("NetworkError", fmt.Sprintf("Failed to read response: %v", err))
	}

	if resp.StatusCode == 200 {
		return true, nil
	}

	// Parse error message
	var errorResp ErrorResponse
	message := fmt.Sprintf("HTTP %d", resp.StatusCode)
	if err := json.Unmarshal(responseBody, &errorResp); err == nil && errorResp.Message != "" {
		message = errorResp.Message
	} else if len(responseBody) > 0 {
		message = fmt.Sprintf("HTTP %d — %s", resp.StatusCode, string(responseBody))
	}

	return false, NewPowerKvError("ServerError", fmt.Sprintf("storeData failed: %s", message))
}

// GetValue retrieves data for the given key
func (p *PowerKv) GetValue(key interface{}) ([]byte, error) {
	keyBytes, err := p.toBytes(key)
	if err != nil {
		return nil, err
	}

	return p.GetValueBytes(keyBytes)
}

// GetValueBytes retrieves data for the given byte key
func (p *PowerKv) GetValueBytes(key []byte) ([]byte, error) {
	keyHex := p.toHexString(key)
	baseURL := p.serverURL + "/getValue"

	params := url.Values{}
	params.Add("projectId", p.projectID)
	params.Add("key", keyHex)

	fullURL := baseURL + "?" + params.Encode()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, NewPowerKvError("NetworkError", fmt.Sprintf("Failed to create request: %v", err))
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, NewPowerKvError("NetworkError", fmt.Sprintf("Request failed: %v", err))
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewPowerKvError("NetworkError", fmt.Sprintf("Failed to read response: %v", err))
	}

	if resp.StatusCode == 200 {
		var response GetValueResponse
		if err := json.Unmarshal(responseBody, &response); err != nil {
			return nil, NewPowerKvError("ServerError", fmt.Sprintf("Unexpected response shape from /getValue: %s", string(responseBody)))
		}

		return p.fromHexString(response.Value)
	}

	// Parse error message
	var errorResp ErrorResponse
	message := fmt.Sprintf("HTTP %d", resp.StatusCode)
	if err := json.Unmarshal(responseBody, &errorResp); err == nil && errorResp.Message != "" {
		message = errorResp.Message
	} else if len(responseBody) > 0 {
		message = fmt.Sprintf("HTTP %d — %s", resp.StatusCode, string(responseBody))
	}

	return nil, NewPowerKvError("ServerError", fmt.Sprintf("getValue failed: %s", message))
}

// GetStringValue retrieves data as string
func (p *PowerKv) GetStringValue(key interface{}) (string, error) {
	data, err := p.GetValue(key)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// GetIntValue retrieves data as int
func (p *PowerKv) GetIntValue(key interface{}) (int, error) {
	data, err := p.GetValue(key)
	if err != nil {
		return 0, err
	}

	value, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, NewPowerKvError("ServerError", fmt.Sprintf("Invalid integer: %v", err))
	}
	return value, nil
}

// GetLongValue retrieves data as int64
func (p *PowerKv) GetLongValue(key interface{}) (int64, error) {
	data, err := p.GetValue(key)
	if err != nil {
		return 0, err
	}

	value, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return 0, NewPowerKvError("ServerError", fmt.Sprintf("Invalid long: %v", err))
	}
	return value, nil
}

// GetDoubleValue retrieves data as float64
func (p *PowerKv) GetDoubleValue(key interface{}) (float64, error) {
	data, err := p.GetValue(key)
	if err != nil {
		return 0, err
	}

	value, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return 0, NewPowerKvError("ServerError", fmt.Sprintf("Invalid double: %v", err))
	}
	return value, nil
}
