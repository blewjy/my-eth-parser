package main

import (
	"eth-parser/pkg"
	"fmt"
	"time"
)

const ADDRESS = "0x95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5"

// main
//
// This program here just shows the operation of the entire Parser interface.
// It will subscribe to the given ADDRESS, and then query its transactions every 2 seconds.
// If there are new transactions, it will be printed to stdout.
//
// In our problem statement, the notification system that is hooking up to this Parser can also
// have a similar operation -- it can poll the Parser periodically, and then send notifications
// to users when it detects new transactions in adjacent polls.
func main() {
	parser := pkg.NewMyEthParser()
	parser.GetCurrentBlock()
	parser.Subscribe(ADDRESS)

	prevLength := 0
	for {
		transactions := parser.GetTransactions(ADDRESS)
		if len(transactions) != prevLength {
			fmt.Println("Got some new transactions...")
			for i, tx := range transactions {
				fmt.Printf("\tTransaction #%d: [From: %s, To: %s, Value: %s, Gas Price: %s]\n", i+1, tx.From, tx.To, tx.Value, tx.GasPrice)
			}
			prevLength = len(transactions)
		}
		time.Sleep(2 * time.Second)
	}
}
