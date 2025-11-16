LOCAL_BIN:=$(CURDIR)/bin

include .env
export

install-deps:
	set GOBIN=$(LOCAL_BIN) && go install github.com/pressly/goose/v3/cmd/goose@latest

local-migration-status:
	goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} status -v

local-migration-up:
	goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} up -v

local-migration-down:
	goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} down -v
