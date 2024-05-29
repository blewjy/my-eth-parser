package pkg

import (
	"eth-parser/internal"
	"fmt"
	"strings"
	"time"
)

// MyEthParser is my implementation of the Parser interface
type MyEthParser struct {
	dao internal.TransactionDao
	gw  internal.EthGateway
}

// NewMyEthParser creates a new Parser.
// By default, we use the internal InMemoryDao and CloudflareEthGateway.
// This can be easily replaced with other implementations where needed.
func NewMyEthParser() Parser {
	parser := &MyEthParser{
		dao: internal.NewInMemoryDao(),
		gw:  internal.NewClouflareEthGateway(),
	}

	go parser.update()

	return parser
}

// update is the goroutine worker function that will update transactions
// for subscribed addresses when new blocks are detected
func (p *MyEthParser) update() {
	defer func() {
		if r := recover(); r != nil {
			// some panic has occurred
			// this is not good, it means that our update goroutine has stopped
			// this shouldn't happen, but if it does, we probably want to emit some alarm in our monitoring system
			fmt.Printf("recover panic: %v\n", r)
		}
	}()

	// if some error occurs during the update, we want to retry earlier
	var earlyRetry bool

	for {
		// get the latest block on the chain
		latestBlockNumber, err := p.gw.GetMostRecentBlockNumber()
		if err != nil {
			fmt.Println(err)
			// if encounter error, then we sleep for 1 second and retry
			time.Sleep(1 * time.Second)
			continue
		}

		// check against the current latest block already synced
		lastSyncedBlockNumber := p.dao.GetLatestBlock()

		// reset this variable
		earlyRetry = false

		// get the list of subscribed addresses
		subscribedAddresses := p.dao.GetAllAddresses()

		// if equal, then do nothing.
		// else not equal, then we want to sync our storage up to the latest block
		if latestBlockNumber != lastSyncedBlockNumber {

			// if the last synced block is 0, means we never synced before
			// so just start from the latest
			if lastSyncedBlockNumber == 0 {
				lastSyncedBlockNumber = latestBlockNumber
			}

			for blockNumber := lastSyncedBlockNumber; blockNumber <= latestBlockNumber; blockNumber++ {
				// fetch transactions by block number
				transactions, err := p.gw.GetTransactionsByBlockNumber(blockNumber)
				if err != nil {
					// if there is some error, we should not indicate that we have synced this block
					// we should also break out of this for loop so that we don't end up with missing blocks in between
					// and, we also want to retry early instead of waiting for 12 seconds
					fmt.Printf("error occurred while querying blockchain: %v\n", err)
					earlyRetry = true
					break
				}

				// for each transaction, if the to_addr/from_addr belongs to some address
				// in our subscription, then save it
				for _, tx := range transactions {

					// technically possible to have the same From and To, but let the dao layer handle it
					if _, from := subscribedAddresses[tx.From]; from {
						p.dao.SaveTransaction(tx.From, tx)
					}
					if _, to := subscribedAddresses[tx.To]; to {
						p.dao.SaveTransaction(tx.To, tx)
					}
				}

				// once all done, then update the latest block in dao
				p.dao.UpdateLatestBlock(blockNumber)
			}
		}

		// on average, Ethereum is designed to have a block interval of 12 seconds
		// so, we query the blockchain every 12 seconds
		// https://ethereum.stackexchange.com/questions/49460/ethereum-blockchain-node-block-terminology
		//
		// but if there is an error, we retry earlier
		if earlyRetry {
			time.Sleep(1 * time.Second)
		} else {
			time.Sleep(12 * time.Second)
		}
	}

}

// GetCurrentBlock will return the latest block that has been synced.
func (p *MyEthParser) GetCurrentBlock() int {
	return p.dao.GetLatestBlock()
}

// Subscribe will add the address given to the list of observed addresses.
func (p *MyEthParser) Subscribe(address string) bool {
	p.dao.CreateAddress(strings.ToLower(address))
	return true
}

// GetTransactions will return all the transactions belonging to the given address,
// starting from the time that it was added to the subscriptions.
func (p *MyEthParser) GetTransactions(address string) []Transaction {
	var transactions []Transaction
	for _, tx := range p.dao.GetTransactions(strings.ToLower(address)) {
		transactions = append(transactions, Transaction{
			BlockHash:        tx.BlockHash,
			BlockNumber:      tx.BlockNumber,
			From:             tx.From,
			Gas:              tx.Gas,
			GasPrice:         tx.GasPrice,
			Hash:             tx.Hash,
			Input:            tx.Input,
			Nonce:            tx.Nonce,
			To:               tx.To,
			TransactionIndex: tx.TransactionIndex,
			Value:            tx.Value,
			Type:             tx.Type,
			ChainId:          tx.ChainId,
		})
	}
	return transactions
}
