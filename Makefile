GOPATH := $(shell go env GOPATH)/bin

export PATH := $(PATH):$(GOPATH)

BIN_DIR = bin
PB_DIR = gen/pb

all: clean generate tidy test build

build:
	go build -o $(BIN_DIR)/ -v ./...

test:
	go test -v ./...

clean:
	go clean
	rm -rf $(BIN_DIR)
	rm -rf $(PB_DIR)/*

tidy:
	go mod tidy
	go mod vendor
	go vet ./...

generate:
	protoc --go_out=$(PB_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PB_DIR) --go-grpc_opt=paths=source_relative \
		graph/*.proto
	protoc --go_out=$(PB_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PB_DIR) --go-grpc_opt=paths=source_relative \
		api/*.proto

postgres:
	docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres16 createdb --username=root --owner=root workflow

dropdb:
	docker exec -it postgres16 dropdb workflow

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/workflow?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/workflow?sslmode=disable" -verbose down

.PHONY: all build test clean tidy generate postgres createdb dropdb migrateup migratedown