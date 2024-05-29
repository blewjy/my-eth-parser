package internal

// TransactionModel is the model for the internal dao layer.
type TransactionModel struct {
	BlockHash            string `json:"blockHash"`
	BlockNumber          string `json:"blockNumber"`
	From                 string `json:"from"`
	Gas                  string `json:"gas"`
	GasPrice             string `json:"gasPrice"`
	MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas"`
	MaxFeePerGas         string `json:"maxFeePerGas"`
	MaxFeePerBlobGas     string `json:"maxFeePerBlobGas"`
	Hash                 string `json:"hash"`
	Input                string `json:"input"`
	Nonce                string `json:"nonce"`
	To                   string `json:"to"`
	TransactionIndex     string `json:"transactionIndex"`
	Value                string `json:"value"`
	Type                 string `json:"type"`
	AccessList           []struct {
		Address     string   `json:"address"`
		StorageKeys []string `json:"storageKeys"`
	} `json:"accessList"`
	ChainId             string   `json:"chainId"`
	V                   string   `json:"v"`
	R                   string   `json:"r"`
	S                   string   `json:"s"`
	BlobVersionedHashes []string `json:"blobVersionedHashes"`
}
