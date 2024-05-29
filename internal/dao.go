package internal

import (
	"sync"
)

// TransactionDao is the data access object interface that allows our Parser
// to interact with our storage layer.
type TransactionDao interface {
	GetLatestBlock() int
	UpdateLatestBlock(block int)
	CreateAddress(address string)
	GetAllAddresses() map[string]int
	GetTransactions(address string) []TransactionModel
	SaveTransaction(address string, transaction TransactionModel)
}

// InMemoryDao is the struct that will handle storage in-memory, by using a simple
// map and a RWMutex.
//
// In this situation, an RWMutex is required to ensure safe concurrent operations on
// the map.
type InMemoryDao struct {
	latestBlock int

	mutex sync.RWMutex
	m     map[string][]TransactionModel
}

// NewInMemoryDao creates a new dao struct and initializing the required.
func NewInMemoryDao() TransactionDao {
	return &InMemoryDao{
		m: make(map[string][]TransactionModel),
	}
}

// GetLatestBlock will return the latest block sync to our data layer.
func (dao *InMemoryDao) GetLatestBlock() int {
	return dao.latestBlock
}

// UpdateLatestBlock will update the latest block number synced. This method should
// be called after all the transactions to this block has been successfully saved.
func (dao *InMemoryDao) UpdateLatestBlock(block int) {
	dao.latestBlock = block
}

// CreateAddress will create a new address to be tracked.
func (dao *InMemoryDao) CreateAddress(address string) {
	dao.mutex.Lock()
	defer dao.mutex.Unlock()

	if _, ok := dao.m[address]; !ok {
		dao.m[address] = []TransactionModel{}
	}
}

// GetAllAddresses will return all the subscribed addresses.
func (dao *InMemoryDao) GetAllAddresses() map[string]int {
	dao.mutex.RLock()
	defer dao.mutex.RUnlock()

	result := map[string]int{}
	for key := range dao.m {
		result[key] = 1
	}
	return result
}

// SaveTransaction will save the given TransactionModel against the given address in the data storage.
// If the address has not yet been created, we will create it anyway.
func (dao *InMemoryDao) SaveTransaction(address string, transaction TransactionModel) {
	dao.mutex.Lock()
	defer dao.mutex.Unlock()

	if _, ok := dao.m[address]; !ok {
		dao.m[address] = []TransactionModel{}
	}

	dao.m[address] = append(dao.m[address], transaction)
}

// GetTransactions will return all the transactions belonging to the given address.
func (dao *InMemoryDao) GetTransactions(address string) []TransactionModel {
	dao.mutex.RLock()
	defer dao.mutex.RUnlock()

	if _, ok := dao.m[address]; !ok {
		return nil
	}

	return dao.m[address]
}
