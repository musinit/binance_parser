package utils

import (
	"binance_parser/domain/model"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func Render(w http.ResponseWriter, httpCode int, v interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	json.NewEncoder(w).Encode(v)
}

func RenderFail(ctx context.Context, w http.ResponseWriter, err error) {
	traceID := ctx.Value(TraceIDCtxKey)
	Render(w, GetStatusCode(err), model.Response{
		Error:      err.Error(),
		Success:    false,
		TrackingId: fmt.Sprint(traceID),
	})
}

func RenderSuccess(ctx context.Context, w http.ResponseWriter, v interface{}) {
	traceID := ctx.Value(TraceIDCtxKey)
	Render(w, http.StatusOK, model.Response{
		Payload:    v,
		Success:    true,
		TrackingId: fmt.Sprint(traceID),
	})
}

type traceIDContextKey string

const (
	TraceIDCtxKey traceIDContextKey = "traceIDContextKey"
)

func GetStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	return http.StatusInternalServerError
}
