package pkg

// Transaction is the model that we expose in our public interface.
// Notice we also have an internal TransactionModel model. Having separate internal data layer
// models and public interface model allows us to logically decouple, and easily configure which
// fields should be exposed and which fields should not.
type Transaction struct {
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	From             string `json:"from"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            string `json:"nonce"`
	To               string `json:"to"`
	TransactionIndex string `json:"transactionIndex"`
	Value            string `json:"value"`
	Type             string `json:"type"`
	ChainId          string `json:"chainId"`
}

// Parser is the interface required by the notifications service.
//
// Assumptions about the notifications service:
// 1. It allows users/clients to subscribe to incoming or outgoing transaction notifications.
// 2. It operates on a polling basis, and will notify users when it detects new transactions during the latest poll.
// 3. The polling interval is arbitrary. It can be as slow as fast as deemed necessary.
type Parser interface {

	// GetCurrentBlock returns the id of the last parsed block.
	GetCurrentBlock() int

	// Subscribe adds the given address to the list of addresses that the Parser should observe.
	// Returns true if added successfully, false otherwise.
	Subscribe(address string) bool

	// GetTransactions returns a list of inbound or outbound transactions for a subscribed address.
	// If the given address has not subscribed before, this will return an empty array.
	GetTransactions(address string) []Transaction
}
