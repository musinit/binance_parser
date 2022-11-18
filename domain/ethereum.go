package domain

import (
	"strconv"
	"strings"
)

type CloudfareRPSRequest struct {
	JSONRPS string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
	ID      int      `json:"id"`
}

type CloudfareRPSResponse struct {
	JSONRPS string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  string `json:"result"`
}

type EtherscansResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

type BlockscautBasicResponse struct {
	BlockNumber      string `json:"blockNumber"`
	Hash             string `json:"hash"`
	Nonce            string `json:"nonce"`
	BlockHash        string `json:"blockHash"`
	TransactionIndex string `json:"transactionIndex"`
	From             string `json:"from"`
	To               string `json:"to"`
	Value            string `json:"value"`
	IsError          string `json:"isError"`
}

type BlockscautTxResponse struct {
	BlockscautBasicResponse `mapstructure:",squash"`
	TimeStamp               interface{} `mapstructure:",omitempty"`
}

func ToBlockskautTransactionDTOs(models []BlockscautTxResponse, selfAddress string) []Transaction {
	result := make([]Transaction, 0)
	for _, val := range models {
		r := val
		value, _ := strconv.ParseInt(val.Value, 10, 64)
		val := r.ToBlockscautTransactionDTO(selfAddress)
		val.Value = float32(value)
		result = append(result, val)
	}
	return result
}

func (t *BlockscautTxResponse) ToBlockscautTransactionDTO(selfAddress string) Transaction {
	blockNumber, _ := strconv.Atoi(t.BlockNumber)
	transactionIndex, _ := strconv.Atoi(t.TransactionIndex)
	value, _ := strconv.ParseFloat(t.Value, 32)
	isError := false
	if t.IsError != "0" {
		isError = true
	}
	to := strings.ToLower(t.To)

	tt := int64(0)
	switch t.TimeStamp.(type) {
	case int:
		tt = t.TimeStamp.(int64)
	}
	return Transaction{
		ID:               nil,
		BlockNumber:      int32(blockNumber),
		TransactionIndex: int32(transactionIndex),
		TimeStamp:        tt,
		Hash:             t.Hash,
		BlockHash:        t.BlockHash,
		From:             t.From,
		To:               to,
		Value:            float32(value),
		IsError:          isError,
	}
}
