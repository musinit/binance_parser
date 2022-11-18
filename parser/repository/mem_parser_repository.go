package repository

import (
	"binance_parser/config"
	"binance_parser/domain"
	"binance_parser/utils"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/mitchellh/mapstructure"
	"log"
	"net/http"
)

type MemParserRepository struct {
	cnf      *config.Config
	logger   *log.Logger
	redisCli redis.UniversalClient

	transactions      map[string][]domain.Transaction
	addrsOffsets      map[string]int32
	latestBlockNumber int32
}

func NewMemParserRepository(
	config *config.Config,
	logger *log.Logger,
	redisCli redis.UniversalClient) domain.ParserRepository {
	return &MemParserRepository{
		cnf:          config,
		logger:       logger,
		redisCli:     redisCli,
		transactions: map[string][]domain.Transaction{},
		addrsOffsets: map[string]int32{},
	}
}

func (p *MemParserRepository) GetCurrentBlock() int32 {
	return p.latestBlockNumber
}

func (p *MemParserRepository) Subscribe(address string) bool {
	if _, ok := p.transactions[address]; ok {
		return false
	}
	p.transactions[address] = []domain.Transaction{}
	// Start tracking for the new ones from the latest block number
	p.addrsOffsets[address] = p.latestBlockNumber
	return true
}

func (p *MemParserRepository) GetTransactions(address string) []domain.Transaction {
	if _, ok := p.transactions[address]; !ok {
		p.logger.Writer().Write([]byte("no transactions for such address"))
		return []domain.Transaction{}
	} else {
		return p.transactions[address]
	}
}

func (p *MemParserRepository) GetAllAddresses() []string {
	result := make([]string, 0)
	for key, _ := range p.transactions {
		result = append(result, key)
	}
	return result
	//ctx := context.Background()
	//iter := p.redisCli.Scan(ctx, 0, fmt.Sprintf("%s*", p.cnf.RedisConfig.Prefix), 0).Iterator()
	//result := make([]string, 0)
	//for iter.Next(ctx) {
	//	result = append(result, iter.Val())
	//}
	//if err := iter.Err(); err != nil {
	//	panic(err)
	//}
	//
	//return result
}

func (p *MemParserRepository) GetAddressOffset(address string) int32 {
	if _, ok := p.transactions[address]; !ok {
		p.logger.Writer().Write([]byte("no offset for such address"))
		return 0
	} else {
		return p.addrsOffsets[address]
	}
}

func (p *MemParserRepository) Unsubscribe(address string) bool {
	if _, ok := p.transactions[address]; !ok {
		p.logger.Writer().Write([]byte("no subscription for such address"))
		return false
	} else {
		// TODO transaction here
		delete(p.transactions, address)
		delete(p.addrsOffsets, address)
		return true
	}
}

// The most buggy part, not tested
func (p *MemParserRepository) GetTransactionsByAddress(address string, startblock int32) ([]domain.Transaction, error) {
	ctx := context.Background()
	url := fmt.Sprintf("%s/eth/mainnet/api?module=account&action=txlist&address=%s&start_block=%d&sort=desc&offset=100&page=0",
		p.cnf.BlockscautAPI, address, startblock)
	req, err := http.NewRequest("GET", url, nil)
	req = req.WithContext(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("BlockscautUsecase err: %s", err.Error())
		return nil, err
	}
	var res domain.EtherscansResponse
	if err = utils.ReadRequestBody(resp.Body, &res); err != nil {
		return nil, err
	}
	if res.Status != "1" {
		if res.Message == "No transactions found" {
			return nil, nil
		} else {
			p.logger.Writer().Write([]byte("BlockscautUsecase bad response, err: %s"))
			return nil, fmt.Errorf("BlockscautUsecase bad response")
		}
	}
	var result []domain.BlockscautTxResponse
	if err := mapstructure.Decode(res.Result, &result); err != nil {
		fmt.Printf("BlockscautUsecase decode err: %s", err.Error())
		return nil, err
	}
	transactions := domain.ToBlockskautTransactionDTOs(result, address)
	return transactions, nil
}

func (p *MemParserRepository) AddTransactions(address string, txs []domain.Transaction) error {
	if _, ok := p.transactions[address]; !ok {
		p.logger.Writer().Write([]byte("no subscription for such address"))
		return errors.New("no such address with subscription")
	} else {
		p.transactions[address] = append(p.transactions[address], txs...)
		return nil
	}
}

func (p *MemParserRepository) UpdateOffset(address string, offset int32) error {
	if _, ok := p.addrsOffsets[address]; !ok {
		p.logger.Writer().Write([]byte("no offset for such address"))
		return errors.New("no such address with offset")
	} else {
		p.addrsOffsets[address] = offset
		return nil
	}
}

func (p *MemParserRepository) GetLatestBlockNumber() int32 {
	return p.latestBlockNumber
}

func (p *MemParserRepository) UpdateLatestBlockNumber(val int32) {
	p.latestBlockNumber = val
}

func (p *MemParserRepository) Close() {
	p.redisCli.Close()
}
