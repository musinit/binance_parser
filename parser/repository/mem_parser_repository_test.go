package repository_test

import (
	"binance_parser/config"
	"binance_parser/domain"
	redisParserRepository "binance_parser/parser/repository"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestMemParserRepository_GetCurrentBlock(t *testing.T) {
	rparserRepository := getMemParserRepository()

	currentBlock := rparserRepository.GetCurrentBlock()

	assert.True(t, currentBlock > 0)
}

func TestMemParserRepository_GetAllTransactions(t *testing.T) {
	rparserRepository := getMemParserRepository()

	addresses := rparserRepository.GetAllAddresses()

	assert.True(t, len(addresses) > 0)
}

func Test_GetETHLatestBlockNumber(t *testing.T) {
	rparserRepository := getMemParserRepository()

	data, err := rparserRepository.GetETHLatestBlockNumber()
	assert.Nil(t, err)
	assert.True(t, data != 0)
}

func Test_GetETHTransactionsByBlockNumber(t *testing.T) {
	rparserRepository := getMemParserRepository()

	data, err := rparserRepository.GetETHLatestBlockNumber()
	assert.Nil(t, err)
	assert.True(t, data != 0)

	block, err := rparserRepository.GetETHBlockByNumber(1207)
	assert.Nil(t, err)
	assert.True(t, len(block.Transactions) > 0)
}

func getMemParserRepository() domain.ParserRepository {
	cnf := config.Config{
		CloudfareConfig: config.CloudfareCnf{
			API: "https://cloudflare-eth.com",
		},
		RedisConfig: config.RedisCnf{
			Addr:   "localhost:6379",
			Prefix: "bp_",
		},
		CloudflareAPI: "https://cloudflare-eth.com",
	}
	logger := log.New(log.Writer(), "binance_parser", 0)
	return redisParserRepository.NewMemParserRepository(&cnf, logger)
}
