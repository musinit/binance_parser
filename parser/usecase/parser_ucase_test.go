package usecase_test

import (
	"binance_parser/config"
	parserRepository "binance_parser/parser/repository"
	parserUsecase "binance_parser/parser/usecase"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestParserUsecase_GetTransactions(t *testing.T) {
	cnf := config.Config{
		CloudfareEHTConfig: config.CloudfareETHCnf{
			API: "https://cloudflare-eth.com",
		},
	}
	logger := log.New(log.Writer(), "binance_parser", 0)
	redisClient := redis.NewUniversalClient(
		&redis.UniversalOptions{
			Addrs:    []string{cnf.RedisConfig.Addr},
			DB:       0,
			ReadOnly: false,
		},
	)
	parserRepo := parserRepository.NewMemParserRepository(&cnf, logger, redisClient)
	parserUcase := parserUsecase.NewParserUsecase(&cnf, parserRepo, logger)

	txs := parserUcase.GetTransactions("test")

	assert.True(t, len(txs) > 0)
}
