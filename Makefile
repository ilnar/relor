BIN_DIR = bin

all: test build

build:
	go build -o $(BIN_DIR)/ -v ./...

test:
	go test -v ./...

clean:
	go clean
	rm -rf $(BIN_DIR)

tidy:
	go mod tidy
	go mod vendor
	go vet ./...

.PHONY: all build test clean tidy