package usecase_test

import (
	"binance_parser/config"
	parserRepository "binance_parser/parser/repository"
	parserUsecase "binance_parser/parser/usecase"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestParserUsecase_GetTransactions(t *testing.T) {
	cnf := config.Config{
		CloudfareConfig: config.CloudfareCnf{
			API: "https://cloudflare-eth.com",
		},
	}
	logger := log.New(log.Writer(), "binance_parser", 0)
	parserRepo := parserRepository.NewMemParserRepository(&cnf, logger)
	parserUcase := parserUsecase.NewParserUsecase(&cnf, parserRepo, logger)

	txs := parserUcase.GetTransactions("test")

	assert.True(t, len(txs) > 0)
}
