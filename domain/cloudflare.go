package domain

import (
	"fmt"
	"strconv"
)

type CloudflareRequest struct {
	ID      int32  `json:"id"`
	JSONRPS string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
}

type CloudflareResponse struct {
	ID      int32  `json:"id"`
	JSONRPS string `json:"jsonrpc"`
	Result  string `json:"method"`
}

type CloudfareRPSBasicResponse struct {
	JSONRPS string `json:"jsonrpc"`
	ID      int    `json:"id"`
}

type CloudflareRPSStringResponse struct {
	CloudfareRPSBasicResponse
	Result string `json:"result"`
}

type CloudflareRPSTransactionsResponse struct {
	CloudfareRPSBasicResponse
	Result CloudflareBlockResponse `json:"result"`
}

type CloudflareBlockResponse struct {
	TransactionsRoot string        `json:"transactionsRoot"`
	Transactions     []Transaction `json:"transactions"`
	Timestamp        string        `json:"timestamp"`
	Hash             string        `json:"hash"`
	GasLimit         string        `json:"gasLimit"`
	GasUsed          string        `json:"gasUsed"`
}

func (cr *CloudflareBlockResponse) ToBlockDTO() (*Block, error) {
	if cr.GasLimit == "" || cr.GasUsed == "" || cr.Timestamp == "" {
		return nil, fmt.Errorf("some of fields are missing")
	}
	gl, err := strconv.ParseInt(cr.GasLimit[2:], 16, 32)
	if err != nil {
		return nil, err
	}
	gu, err := strconv.ParseInt(cr.GasUsed[2:], 16, 32)
	if err != nil {
		return nil, err
	}
	timestamp, err := strconv.ParseInt(cr.Timestamp[2:], 16, 32)
	if err != nil {
		return nil, err
	}
	return &Block{
		TransactionsRoot: cr.TransactionsRoot,
		Transactions:     cr.Transactions,
		Timestamp:        int32(timestamp),
		Hash:             cr.Hash,
		GasLimit:         int32(gl),
		GasUsed:          int32(gu),
	}, nil
}
