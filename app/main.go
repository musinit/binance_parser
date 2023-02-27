package main

import (
	"binance_parser/config"
	"binance_parser/domain"
	"binance_parser/txengine"
	"binance_parser/utils"
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"os"
	"time"

	parserHTTPHandler "binance_parser/parser/delivery"
	parserRepository "binance_parser/parser/repository"
	parserUsecase "binance_parser/parser/usecase"
)

func main() {
	// TODO read from config file
	cnf := config.Config{
		Port: "9090",
		CloudfareConfig: config.CloudfareCnf{
			API: "https://cloudflare-eth.com",
		},
		RedisConfig: config.RedisCnf{
			Addr:   "localhost:6379",
			Prefix: "bp_",
		},
		CloudflareAPI: "https://cloudflare-eth.com",
	}
	r := utils.NewRouter()
	r.UseMiddlewares(
		utils.TraceID,
	)
	logger := log.New(log.Writer(), "binance_parser", 0)
	ctx := context.Background()

	parserRepository := parserRepository.NewMemParserRepository(&cnf, logger)
	defer parserRepository.Close()

	SetDomain(&cnf, r, parserRepository, logger)

	env := os.Getenv("ENV")
	if env == "prod" {
		s := &http.Server{
			Addr:           fmt.Sprint(":https"),
			Handler:        chi.ServerBaseContext(ctx, r.Mux),
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   5 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		logger.Writer().Write([]byte(fmt.Sprintf("starting service on https\n")))
		log.Fatal(s.ListenAndServeTLS(fmt.Sprintf("%sverdandi.uno.pem", "/app/"), fmt.Sprintf("%sverdandi.uno.key", "/app/")))
	} else {
		logger.Writer().Write([]byte(fmt.Sprintf("starting service on port: %s\n", cnf.Port)))
		s := &http.Server{
			Addr:           fmt.Sprint(fmt.Sprintf(":%s", cnf.Port)),
			Handler:        chi.ServerBaseContext(ctx, r.Mux),
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   5 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		log.Fatal(s.ListenAndServe())
	}

}

func SetDomain(
	cnf *config.Config,
	r *utils.Router,
	parserRedisRepository domain.ParserRepository,
	logger *log.Logger,
) {

	parserUsecase := parserUsecase.NewParserUsecase(cnf, parserRedisRepository, logger)
	parserHTTPHandler.NewParserHTTPHandler(r, parserUsecase, logger)

	txengine := txengine.NewTxengine(
		context.Background(),
		logger,
		cnf,
		parserUsecase,
	)
	go txengine.Run()
}
