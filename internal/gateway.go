package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const (
	GATEWAY_URL     = "https://cloudflare-eth.com"
	GATEWAY_TIMEOUT = 2 * time.Second
)

// EthGateway is the gateway interface that allows our Parser to interact with the gateway APIs.
// In this case, these are the Ethereum JSON RPC APIs that we will be calling.
type EthGateway interface {
	GetMostRecentBlockNumber() (int, error)
	GetTransactionsByBlockNumber(block int) ([]TransactionModel, error)
}

// CloudflareEthGateway is the implementation for the EthGateway that uses the Cloudflare endpoint.
type CloudflareEthGateway struct {
	client *http.Client
}

// NewClouflareEthGateway will create a new EthGateway.
func NewClouflareEthGateway() EthGateway {
	return &CloudflareEthGateway{
		client: &http.Client{
			Timeout: GATEWAY_TIMEOUT,
		},
	}
}

// GetMostRecentBlockNumber will call the "eth_blockNumber" API to retrieve the latest block number in the blockchain.
func (gw *CloudflareEthGateway) GetMostRecentBlockNumber() (int, error) {
	respBytes, err := gw.post(map[string]interface{}{
		// just use time.Now() to keep it simple
		// in a proper implementation, I may choose to generate something more meaningful
		// allows us to potentially keep track of the requests made
		"id":      time.Now().Unix(),
		"jsonrpc": "2.0",
		"method":  "eth_blockNumber",
		"params":  []interface{}{},
	})
	if err != nil {
		return 0, err
	}

	// Parse the response body
	var resp struct {
		BlockNumberHex string `json:"result"`
	}
	if err = json.Unmarshal(respBytes, &resp); err != nil {
		fmt.Println("Error parsing JSON response:", err)
		return 0, err
	}

	// parse the hex string
	blockNumber, err := strconv.ParseInt(resp.BlockNumberHex[2:], 16, 64)
	if err != nil {
		return 0, errors.New("failed to parse the result")
	}

	return int(blockNumber), nil
}

// GetTransactionsByBlockNumber will take a block number, and call the "eth_getBlockByNumber" API to
// retrieve all the transactions belonging to this block number.
func (gw *CloudflareEthGateway) GetTransactionsByBlockNumber(block int) ([]TransactionModel, error) {
	blockNumberHex := fmt.Sprintf("0x%s", strconv.FormatInt(int64(block), 16))
	fmt.Println("[GetTransactionsByBlockNumber]", blockNumberHex, block)
	respBytes, err := gw.post(map[string]interface{}{
		"id":      time.Now().Unix(),
		"jsonrpc": "2.0",
		"method":  "eth_getBlockByNumber",
		"params": []interface{}{
			blockNumberHex,
			true,
		},
	})
	if err != nil {
		return nil, err
	}

	// Parse the response body
	var resp struct {
		Result struct {
			Transactions []TransactionModel `json:"transactions"`
		} `json:"result"`
	}
	if err = json.Unmarshal(respBytes, &resp); err != nil {
		fmt.Println("Error parsing JSON response:", err)
		return nil, err
	}
	return resp.Result.Transactions, nil
}

// post is a general helper method to fire HTTP POST requests to the GATEWAY_URL.
func (gw *CloudflareEthGateway) post(data map[string]interface{}) ([]byte, error) {
	v, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", GATEWAY_URL, bytes.NewBuffer(v))
	req.Header.Set("Content-Type", "application/json")
	resp, err := gw.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
