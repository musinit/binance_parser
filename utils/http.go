package utils

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
)

func ReadRequestBody(b io.ReadCloser, t interface{}) error {
	if reflect.ValueOf(t).Kind() != reflect.Ptr {
		return errors.New("ERR")
	}
	body, err := ioutil.ReadAll(b)
	if err != nil {
		return errors.New("ERR")
	}
	if err = json.Unmarshal(body, &t); err != nil {
		log.Printf(err.Error())
		return err
	}
	return nil
}

func TraceID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID, _ := uuid.NewRandom()
		ctx := context.WithValue(r.Context(), TraceIDCtxKey, traceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
