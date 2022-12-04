package repository

import (
	"binance_parser/config"
	"binance_parser/domain"
	"binance_parser/utils"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// let's start by default not from the beginning of blockchain,
// but from some middle point
var startBlockNumber = int32(16110556)

type MemParserRepository struct {
	cnf    *config.Config
	logger *log.Logger

	// transactions per block
	transactions map[string][]domain.Transaction
	// transactions per user address
	userTransactions  map[string][]*domain.Transaction
	userAddressTxHash map[string]struct{}
	latestBlockNumber int32
}

func NewMemParserRepository(
	config *config.Config,
	logger *log.Logger) domain.ParserRepository {
	return &MemParserRepository{
		cnf:               config,
		logger:            logger,
		transactions:      map[string][]domain.Transaction{},
		userTransactions:  map[string][]*domain.Transaction{},
		userAddressTxHash: map[string]struct{}{},
		latestBlockNumber: startBlockNumber,
	}
}

func (p *MemParserRepository) GetCurrentBlock() int32 {
	return p.latestBlockNumber
}

func (p *MemParserRepository) Subscribe(address string) bool {
	if _, ok := p.userTransactions[address]; ok {
		return false
	}
	p.userTransactions[address] = []*domain.Transaction{}
	return true
}

func (p *MemParserRepository) GetTransactions(address string) []*domain.Transaction {
	if _, ok := p.userTransactions[address]; !ok {
		p.logger.Writer().Write([]byte("no subscription for this address"))
		return []*domain.Transaction{}
	} else {
		return p.userTransactions[address]
	}
}

func (p *MemParserRepository) GetAllAddresses() map[string]struct{} {
	var result = map[string]struct{}{}
	for key, _ := range p.userTransactions {
		result[key] = struct{}{}
	}
	return result
}

func (p *MemParserRepository) Unsubscribe(address string) bool {
	if _, ok := p.userTransactions[address]; !ok {
		p.logger.Writer().Write([]byte("no subscription for such address"))
		return false
	} else {
		// TODO transaction here
		delete(p.userTransactions, address)
		return true
	}
}

func (p *MemParserRepository) CacheTransaction(address string, tx *domain.Transaction) error {
	if _, ok := p.userTransactions[address]; !ok {
		p.logger.Writer().Write([]byte("no subscription for such address"))
		return errors.New("no such address with subscription")
	} else {
		p.userTransactions[address] = append(p.userTransactions[address], tx)
		return nil
	}
}

func (p *MemParserRepository) AddTransactionsInBlock(blockHash string, txs []domain.Transaction) {
	p.transactions[blockHash] = append(p.transactions[blockHash], txs...)
}

func (p *MemParserRepository) AddUserAddressTxHash(address string, txHash string) {
	p.userAddressTxHash[domain.TxcacheKey(address, txHash)] = struct{}{}
}

func (p *MemParserRepository) IsUserAddressInTxhash(address string, txHash string) bool {
	_, ok := p.userAddressTxHash[domain.TxcacheKey(address, txHash)]
	return ok
}

func (p *MemParserRepository) GetLatestBlockNumber() int32 {
	return p.latestBlockNumber
}

func (p *MemParserRepository) UpdateLatestBlockNumber(val int32) {
	p.latestBlockNumber = val
}

func (p *MemParserRepository) GetETHLatestBlockNumber() (int32, error) {
	ctx := context.Background()
	body := domain.CloudflareRequest{
		JSONRPS: "2.0",
		Method:  "eth_blockNumber",
		Params:  nil,
		ID:      1,
	}
	bJSON, err := json.Marshal(body)
	if err != nil {
		return -1, err
	}
	req, err := http.NewRequest("POST", p.cnf.CloudflareAPI, bytes.NewBuffer(bJSON))
	req = req.WithContext(ctx)
	if err != nil {
		return -1, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Cloudflare err: %s", err.Error())
		return -1, err
	}
	var res domain.CloudflareRPSStringResponse
	if err = utils.ReadRequestBody(resp.Body, &res); err != nil {
		return -1, err
	}
	result, err := strconv.ParseInt(res.Result[2:], 16, 32)
	if err != nil {
		return -1, err
	}
	return int32(result), nil
}

func (p *MemParserRepository) GetETHBlockByNumber(blockNumber int32) (*domain.Block, error) {
	blockNumberHex := fmt.Sprintf("0x%x", blockNumber)
	ctx := context.Background()
	body := domain.CloudflareRequest{
		JSONRPS: "2.0",
		Method:  "eth_getBlockByNumber",
		Params: []interface{}{
			blockNumberHex,
			true,
		},
		ID: 1,
	}
	bJSON, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", p.cnf.CloudflareAPI, bytes.NewBuffer(bJSON))
	req = req.WithContext(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Cloudflare err: %s", err.Error())
		return nil, err
	}
	var res domain.CloudflareRPSTransactionsResponse
	if err = utils.ReadRequestBody(resp.Body, &res); err != nil {
		return nil, err
	}

	return res.Result.ToBlockDTO()
}

func (p *MemParserRepository) Close() {
}
