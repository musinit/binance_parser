package delivery

import (
	"binance_parser/domain"
	"binance_parser/utils"
	"errors"
	"log"
	"net/http"
)

type ParserHTTPHandler struct {
	parserUsecase domain.ParserUsecase
	logger        *log.Logger
}

func NewParserHTTPHandler(router *utils.Router,
	parserUsecase domain.ParserUsecase,
	logger *log.Logger) {
	handler := &ParserHTTPHandler{
		logger:        logger,
		parserUsecase: parserUsecase,
	}
	router.Mux.Get("/parser/blocks/current", handler.GetCurrentBlockNumber)
	router.Mux.Get("/parser/overview", handler.GetOverview)
	router.Mux.Post("/parser/transactions", handler.GetTransactions)
	router.Mux.Post("/parser/subscription", handler.Subscribe)
	router.Mux.Get("/parser/address", handler.GetAllAddresses)
}

func (h *ParserHTTPHandler) GetCurrentBlockNumber(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	currentBlock := h.parserUsecase.GetCurrentBlock()

	utils.RenderSuccess(ctx, w, currentBlock)
}

func (h *ParserHTTPHandler) GetOverview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	currentBlock := h.parserUsecase.GetCurrentBlock()
	txnum := h.parserUsecase.GetTxNum()

	parserOverview := domain.ParserOverview{
		LatestBlockNumber: currentBlock,
		TxNum:             txnum,
	}

	utils.RenderSuccess(ctx, w, parserOverview)
}

func (h *ParserHTTPHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error
	var req domain.AddressRequest
	if err = utils.ReadRequestBody(r.Body, &req); err != nil {
		utils.RenderFail(ctx, w, err)
		return
	}
	transactions := h.parserUsecase.GetTransactions(req.Address)
	utils.RenderSuccess(ctx, w, transactions)
}

func (h *ParserHTTPHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	var req domain.AddressRequest
	if err = utils.ReadRequestBody(r.Body, &req); err != nil {
		utils.RenderFail(ctx, w, err)
		return
	}

	if val := h.parserUsecase.Subscribe(req.Address); !val {
		utils.RenderFail(ctx, w, errors.New("some error"))
		return
	}
	utils.RenderSuccess(ctx, w, "done")
}

func (h *ParserHTTPHandler) GetAllAddresses(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	addresses := h.parserUsecase.GetAllAddresses()
	utils.RenderSuccess(ctx, w, addresses)
}
