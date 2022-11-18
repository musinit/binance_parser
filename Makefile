BINARY=binance_parser

build:
	go build -o ${BINARY} app/main.go

run: build
	ENV=dev ./${BINARY}

run_env:
	docker-compose -f docker-compose-env.yml up -d

stop_env:
	docker-compose -f docker-compose-env.yml down