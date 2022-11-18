package domain

import "github.com/google/uuid"

type ParserUsecase interface {
	// last parsed block
	GetCurrentBlock() int32

	// add address to observer
	Subscribe(address string) bool

	// list of inbound or outbound transactions for an address
	GetTransactions(address string) []Transaction

	GetAllAddresses() []string
	GetTransactionsByAddress(address string, startblock int32) ([]Transaction, error)
	GetAddressOffset(address string) int32
	Unsubscribe(address string) bool
	AddTransactions(address string, transactions []Transaction) error
	UpdateOffset(address string, offset int32) error
	UpdateLatestBlockNumber(val int32)
	GetLatestBlockNumber() int32
}

type ParserRepository interface {
	GetCurrentBlock() int32
	Subscribe(address string) bool
	Unsubscribe(address string) bool
	GetTransactions(address string) []Transaction
	GetAllAddresses() []string
	GetTransactionsByAddress(address string, startblock int32) ([]Transaction, error)
	GetAddressOffset(address string) int32
	AddTransactions(address string, transactions []Transaction) error
	UpdateOffset(address string, offset int32) error
	UpdateLatestBlockNumber(val int32)
	GetLatestBlockNumber() int32
	Close()
}

type Transaction struct {
	ID               *uuid.UUID `json:"id"`
	BlockNumber      int32      `json:"blockNumber"`
	TransactionIndex int32      `json:"transactionIndex"`
	TimeStamp        int64      `json:"timeStamp"`
	Hash             string     `json:"hash"`
	BlockHash        string     `json:"blockHash"`
	From             string     `json:"from"`
	To               string     `json:"to"`
	Value            float32    `json:"value"`
	IsError          bool       `json:"isError"`
}

type AddressRequest struct {
	Address string `json:"address"`
}
