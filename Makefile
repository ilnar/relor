GOPATH := $(shell go env GOPATH)/bin

export PATH := $(PATH):$(GOPATH)

BIN_DIR = bin
PB_DIR = gen/pb
SQLC_DIR = gen/sqlc

DB = postgres16

all: clean generate tidy test build 

build:
	go build -o $(BIN_DIR)/ -v ./...

initpg: cleanpg
	docker run --name $(DB) -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine
	sleep 5
	docker exec -it postgres16 createdb --username=root --owner=root workflow
	docker stop $(DB)

startpg: initpg migrateup
	docker start $(DB)

cleanpg:
	docker rm -fv $(DB)

migrateup:
	docker start $(DB)
	sleep 2
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/workflow?sslmode=disable" -verbose up
	docker stop $(DB)

migratedown:
	docker start $(DB)
	sleep 2
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/workflow?sslmode=disable" -verbose down
	docker stop $(DB)

test:
	go test -v ./...

clean:
	go clean
	rm -rf $(BIN_DIR)
	rm -rf $(PB_DIR)/*
	rm -rf $(SQLC_DIR)

tidy:
	go mod tidy
	go mod vendor
	go vet ./...

generate: sqlc
	protoc --go_out=$(PB_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PB_DIR) --go-grpc_opt=paths=source_relative \
		graph/*.proto
	protoc --go_out=$(PB_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PB_DIR) --go-grpc_opt=paths=source_relative \
		api/*.proto

sqlc:
	docker run --rm -v $(shell pwd):/src -w /src/db sqlc/sqlc:1.27.0 generate

.PHONY: all build test clean tidy generate \
	initpg startpg cleanpg migrateup migratedown sqlc