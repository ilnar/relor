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
		api/*.proto

.PHONY: all build test clean tidy generate