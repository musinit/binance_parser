package domain

type CloudfareRPSRequest struct {
	JSONRPS string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
	ID      int      `json:"id"`
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
