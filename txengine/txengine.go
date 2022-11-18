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
}

func NewTxengine(ctx context.Context,
	logger *log.Logger,
	cnf *config.Config,
	parserUsecase domain.ParserUsecase) *Txengine {
	return &Txengine{
		ctx:           ctx,
		logger:        logger,
		cnf:           cnf,
		stopCh:        make(chan bool, 1),
		parserUsecase: parserUsecase,
	}
}

func (engine *Txengine) Run() {
	engine.start()
}

func (engine *Txengine) start() {
	engine.logger.Writer().Write([]byte("txengine started\n"))
	for {
		select {
		case <-time.After(5 * time.Second):
			addresses := engine.parserUsecase.GetAllAddresses()
			for _, address := range addresses {
				fmt.Printf("processing %s address...\n", address)
				startIdx := engine.parserUsecase.GetAddressOffset(address)
				fmt.Printf("offset: %d\n", startIdx)
				// TODO startIdx+1 rewrite to check by startIdx + transactionIdx inside block
				txs, err := engine.parserUsecase.GetTransactionsByAddress(address, startIdx+1)
				fmt.Printf("transactions:: %+v\n", txs)
				if err != nil {
					engine.removeAddressFromCollection(address)
				}
				if len(txs) > 0 {
					// TODO transaction here
					if err := engine.parserUsecase.AddTransactions(address, txs); err != nil {
						// TODO retry here
						engine.removeAddressFromCollection(address)
					}
					if err := engine.parserUsecase.UpdateOffset(address, txs[0].BlockNumber); err != nil {
						engine.removeAddressFromCollection(address)
					}
					latestBlockNumber := engine.parserUsecase.GetLatestBlockNumber()
					actualLatestBlockNumber := txs[0].BlockNumber
					if actualLatestBlockNumber > latestBlockNumber {
						engine.parserUsecase.UpdateLatestBlockNumber(actualLatestBlockNumber)
					}
					fmt.Printf("%s successfully processed\n", address)
				} else {
					fmt.Printf("%s no transactions for adding\n", address)
				}

			}
			fmt.Println("processing")
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
