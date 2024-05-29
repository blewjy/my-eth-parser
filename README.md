# Goal
Implement Ethereum blockchain parser that will allow to query transactions for subscribed
addresses.

# Problem
Users not able to receive push notifications for incoming/outgoing transactions. By
Implementing Parser interface we would be able to hook this up to notifications service to
notify about any incoming/outgoing transactions.

# Limitations
- Use Go Language 
- Avoid usage of external libraries 
- Use Ethereum JSONRPC to interact with Ethereum Blockchain 
- Use memory storage for storing any data (should be easily extendable to support any
storage in the future)

Expose public interface for external usage either via code or command line or rest api that will
include supported list of operations defined in the Parser interface

``` golang
type Parser interface {
    // last parsed block
    GetCurrentBlock() int
    // add address to observer
    Subscribe(address string) bool
    // list of inbound or outbound transactions for an address
    GetTransactions(address string) []Transaction
}
```
