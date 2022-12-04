package txengine

import (
	"binance_parser/config"
	"binance_parser/domain"
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

type Txengine struct {
	ctx    context.Context
	logger *log.Logger
	cnf    *config.Config
	stopCh chan bool

	parserUsecase domain.ParserUsecase

	newTransactions chan []domain.Transaction
}

func NewTxengine(ctx context.Context,
	logger *log.Logger,
	cnf *config.Config,
	parserUsecase domain.ParserUsecase) *Txengine {
	return &Txengine{
		ctx:             ctx,
		logger:          logger,
		cnf:             cnf,
		stopCh:          make(chan bool, 1),
		parserUsecase:   parserUsecase,
		newTransactions: make(chan []domain.Transaction),
	}
}

func (engine *Txengine) Run() {
	go engine.startTxSync()
	engine.startBlockSync()
}

// startBlockSync for copy the blocks and transactions from Ethereum
func (engine *Txengine) startBlockSync() {
	engine.logger.Writer().Write([]byte("txengine started\n"))
	for {
		select {
		case <-time.After(5 * time.Second):
			ethLatestBlockNumber, err := engine.parserUsecase.GetETHLatestBlockNumber()
			if err != nil {
				fmt.Errorf("txengine: request latest block number error: %s. Retrying in 5 seconds...\n", err.Error())
				time.Sleep(5 * time.Second)
				continue
			}
			// TODO fix
			// start check with the previous block, as there can be transactions since last check.
			// I assume that there can't be 2 or more filled blocks in less than 5 seconds
			// massive repeat caching tries because of this
			currBlockNumber := engine.parserUsecase.GetLatestBlockNumber() - 1
			fmt.Printf("latest eth block number: %d, latest processed block number: %d. Processing %d blocks...\n", ethLatestBlockNumber, currBlockNumber, ethLatestBlockNumber-currBlockNumber)
			// for each unprocessed block number
			for blockNumber := currBlockNumber; blockNumber <= ethLatestBlockNumber; blockNumber++ {
				block, err := engine.parserUsecase.GetETHBlockByNumber(blockNumber)
				if err != nil {
					fmt.Errorf("txengine: request block by number %d error: %s. Retrying in 5 seconds...\n", blockNumber, err.Error())
					time.Sleep(5 * time.Second)
					break
				}
				fmt.Printf("hash: %+v\ngasUsed: %d\ngasLimit: %d\n", block.Hash, block.GasUsed, block.GasLimit)
				if len(block.Transactions) > 0 {
					// TODO add only transactions for subscribed addresses
					engine.parserUsecase.AddTransactionsInBlock(block.Hash, block.Transactions)
					engine.newTransactions <- block.Transactions
				} else {
					fmt.Printf("no transactions in block hash: %s. Continue...", block.Hash)
				}
			}
			engine.parserUsecase.UpdateLatestBlockNumber(ethLatestBlockNumber)
		case <-engine.stopCh:
			break
		}
	}
}

func (engine *Txengine) removeAddressFromCollection(address string) error {
	if res := engine.parserUsecase.Unsubscribe(address); !res {
		engine.logger.Writer().Write([]byte("can't remove bad address from collection"))
		return errors.New("can't remove bad address from collection")
	}
	return nil
}
