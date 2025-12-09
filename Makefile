LOCAL_BIN:=$(CURDIR)/bin

include .env.example
export

install-deps:
	set GOBIN=$(LOCAL_BIN) && go install github.com/pressly/goose/v3/cmd/goose@latest

install-golangci-lint:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.7.2

lint:
	golangci-lint run --config .golangci.yml

lint-fix:
	golangci-lint run --config .golangci.yml --fix

local-migration-status:
	goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} status -v

local-migration-up:
	goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} up -v

local-migration-down:
	goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} down -v
