package usecase

import (
	"binance_parser/config"
	"binance_parser/domain"
	"fmt"
	"log"
)

type ParserUsecase struct {
	cnf              *config.Config
	parserRepository domain.ParserRepository
	logger           *log.Logger
}

func NewParserUsecase(config *config.Config,
	parserRep domain.ParserRepository,
	logger *log.Logger) *ParserUsecase {
	return &ParserUsecase{
		cnf:              config,
		parserRepository: parserRep,
		logger:           logger,
	}
}

func (p *ParserUsecase) GetCurrentBlock() int32 {
	return p.parserRepository.GetCurrentBlock()
}

func (p *ParserUsecase) Subscribe(address string) bool {
	return p.parserRepository.Subscribe(address)
}

func (p *ParserUsecase) GetTransactions(address string) []*domain.Transaction {
	return p.parserRepository.GetTransactions(address)
}

func (p *ParserUsecase) Unsubscribe(address string) bool {
	return p.parserRepository.Unsubscribe(address)
}

func (p *ParserUsecase) GetAllAddresses() map[string]struct{} {
	return p.parserRepository.GetAllAddresses()
}

func (p *ParserUsecase) CacheTransaction(address string, tx *domain.Transaction) error {
	if !p.parserRepository.IsUserAddressInTxhash(address, tx.Hash) {
		p.parserRepository.AddUserAddressTxHash(address, tx.Hash)
		return p.parserRepository.CacheTransaction(address, tx)
	} else {
		fmt.Printf("ERROR trying to cache tx %s for address %s twice", tx.Hash, address)
		return nil
	}
}

func (p *ParserUsecase) AddTransactionsInBlock(blockHash string, txs []domain.Transaction) {
	p.parserRepository.AddTransactionsInBlock(blockHash, txs)
}

func (p *ParserUsecase) UpdateLatestBlockNumber(val int32) {
	p.parserRepository.UpdateLatestBlockNumber(val)
}

func (p *ParserUsecase) GetLatestBlockNumber() int32 {
	return p.parserRepository.GetLatestBlockNumber()
}

func (p *ParserUsecase) GetTxNum() int {
	return p.parserRepository.GetTxNum()
}

func (p *ParserUsecase) GetETHLatestBlockNumber() (int32, error) {
	return p.parserRepository.GetETHLatestBlockNumber()
}

func (p *ParserUsecase) GetETHBlockByNumber(blockNumber int32) (*domain.Block, error) {
	return p.parserRepository.GetETHBlockByNumber(blockNumber)
}
