package repository_test

import (
	"binance_parser/config"
	"binance_parser/domain"
	redisParserRepository "binance_parser/parser/repository"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestRedisParserRepository_GetCurrentBlock(t *testing.T) {
	rparserRepository := getRedisParserRepository()

	currentBlock := rparserRepository.GetCurrentBlock()

	assert.True(t, currentBlock > 0)
}

func TestRedisParserRepository_GetAllTransactions(t *testing.T) {
	rparserRepository := getRedisParserRepository()

	addresses := rparserRepository.GetAllAddresses()

	assert.True(t, len(addresses) > 0)
}

func TestRedisParserRepository_GetTransactionsByAddress(t *testing.T) {
	rparserRepository := getRedisParserRepository()

	testAddress := "0xcd0fBe49Ac5e009858DFd1c5F7330907C710fe96"
	txs, err := rparserRepository.GetTransactionsByAddress(testAddress, 0)

	assert.Nil(t, err)
	assert.True(t, len(txs) > 0)
}

func getRedisParserRepository() domain.ParserRepository {
	cnf := config.Config{
		CloudfareEHTConfig: config.CloudfareETHCnf{
			API: "https://cloudflare-eth.com",
		},
		RedisConfig: config.RedisCnf{
			Addr:   "localhost:6379",
			Prefix: "bp_",
		},
		BlockscautAPI: "https://blockscout.com",
	}
	redisClient := redis.NewUniversalClient(
		&redis.UniversalOptions{
			Addrs:    []string{cnf.RedisConfig.Addr},
			DB:       0,
			ReadOnly: false,
		},
	)
	logger := log.New(log.Writer(), "binance_parser", 0)
	return redisParserRepository.NewRedisParserRepository(&cnf, logger, redisClient)
}
