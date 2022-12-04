package domain

import (
	"fmt"
	"github.com/google/uuid"
)

type ParserUsecase interface {
	// last parsed block
	GetCurrentBlock() int32

	// add address to observer
	Subscribe(address string) bool

	// list of inbound or outbound transactions for an address
	GetTransactions(address string) []*Transaction

	GetAllAddresses() map[string]struct{}
	Unsubscribe(address string) bool
	UpdateLatestBlockNumber(val int32)
	GetLatestBlockNumber() int32

	GetETHLatestBlockNumber() (int32, error)
	GetETHBlockByNumber(blockNumber int32) (*Block, error)
	AddTransactionsInBlock(blockHash string, txs []Transaction)
	CacheTransaction(address string, transaction *Transaction) error
}

type ParserRepository interface {
	GetCurrentBlock() int32
	Subscribe(address string) bool
	Unsubscribe(address string) bool
	GetTransactions(address string) []*Transaction
	GetAllAddresses() map[string]struct{}
	UpdateLatestBlockNumber(val int32)
	GetLatestBlockNumber() int32

	GetETHLatestBlockNumber() (int32, error)
	GetETHBlockByNumber(blockNumber int32) (*Block, error)
	AddTransactionsInBlock(blockHash string, txs []Transaction)
	CacheTransaction(address string, transaction *Transaction) error
	AddUserAddressTxHash(address string, txHash string)
	IsUserAddressInTxhash(address string, txHash string) bool
	Close()
}

type Transaction struct {
	ID               *uuid.UUID `json:"id"`
	BlockNumber      string     `json:"blockNumber"`
	TransactionIndex string     `json:"transactionIndex"`
	TimeStamp        int64      `json:"timeStamp"`
	Hash             string     `json:"hash"`
	BlockHash        string     `json:"blockHash"`
	From             string     `json:"from"`
	To               string     `json:"to"`
	Value            string     `json:"value"`
}

type Block struct {
	TransactionsRoot string        `json:"transactionsRoot"`
	Transactions     []Transaction `json:"transactions"`
	Timestamp        int32         `json:"timestamp"`
	Hash             string        `json:"hash"`
	GasLimit         int32         `json:"gasLimit"`
	GasUsed          int32         `json:"gasUsed"`
}

type AddressRequest struct {
	Address string `json:"address"`
}

var TxcacheKey = func(address, txhash string) string {
	return fmt.Sprintf("%s_%s", address, txhash)
}
