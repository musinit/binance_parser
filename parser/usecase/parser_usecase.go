package usecase

import (
	"binance_parser/config"
	"binance_parser/domain"
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

func (p *ParserUsecase) GetTransactions(address string) []domain.Transaction {
	return p.parserRepository.GetTransactions(address)
}

func (p *ParserUsecase) GetTransactionsByAddress(address string, startblock int32) ([]domain.Transaction, error) {
	return p.parserRepository.GetTransactionsByAddress(address, startblock)
}

func (p *ParserUsecase) GetAddressOffset(address string) int32 {
	return p.parserRepository.GetAddressOffset(address)
}

func (p *ParserUsecase) Unsubscribe(address string) bool {
	return p.parserRepository.Unsubscribe(address)
}

func (p *ParserUsecase) GetAllAddresses() []string {
	return p.parserRepository.GetAllAddresses()
}

func (p *ParserUsecase) AddTransactions(address string, txs []domain.Transaction) error {
	return p.parserRepository.AddTransactions(address, txs)
}

func (p *ParserUsecase) UpdateOffset(address string, offset int32) error {
	return p.parserRepository.UpdateOffset(address, offset)
}

func (p *ParserUsecase) UpdateLatestBlockNumber(val int32) {
	p.parserRepository.UpdateLatestBlockNumber(val)
}

func (p *ParserUsecase) GetLatestBlockNumber() int32 {
	return p.parserRepository.GetLatestBlockNumber()
}
