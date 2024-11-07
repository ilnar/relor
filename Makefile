GOPATH := $(shell go env GOPATH)/bin

export PATH := $(PATH):$(GOPATH)

BIN_DIR = bin
PB_DIR = gen/pb
SQLC_DIR = gen/sqlc

SQLC_VERSION = 1.27.0
MIGRATE_VERSION = v4.18.1

DB = postgres16

all: clean generate tidy test build 

build:
	go build -o $(BIN_DIR)/ -v ./...

cov:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

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
	docker run -v ./db/migration:/migration --network host migrate/migrate:${MIGRATE_VERSION} \
		-path=/migration/ -database "postgresql://root:secret@localhost:5432/workflow?sslmode=disable" -verbose up
	docker stop $(DB)

migratedown:
	docker start $(DB)
	sleep 2
	docker run -v ./db/migration:/migration --network host migrate/migrate:${MIGRATE_VERSION} \
		-path=/migration/ -database "postgresql://root:secret@localhost:5432/workflow?sslmode=disable" -verbose down
	docker stop $(DB)

test:
	go test -v ./...

clean:
	go clean
	rm -rf $(BIN_DIR)
	rm -rf $(PB_DIR)/*
	rm -rf $(SQLC_DIR)
	#TODO remove all docker containers

tidy:
	go mod tidy
	go mod vendor
	go vet ./...

generate: sqlc proto

sqlc:
	docker run --rm -v .:/src -w /src/db sqlc/sqlc:${SQLC_VERSION} generate

proto:
	docker build -f Dockerfile.protoc -t protoc-tool .
	docker run  --rm -v .:/src -w /src protoc-tool \
		--go_out=$(PB_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PB_DIR) --go-grpc_opt=paths=source_relative \
		graph/*.proto
	docker run  --rm -v .:/src -w /src protoc-tool \
		--go_out=$(PB_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PB_DIR) --go-grpc_opt=paths=source_relative \
		api/*.proto
	docker run  --rm -v .:/src -w /src protoc-tool \
		--go_out=$(PB_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PB_DIR) --go-grpc_opt=paths=source_relative \
		config/*.proto

.PHONY: all build test clean cov tidy generate \
	initpg startpg cleanpg migrateup migratedown sqlc proto