include .env
export $(shell sed 's/=.*//' .env)

tern_migrate:
	tern migrate --migrations ./internal/store/pgstore/migrations --config ./internal/store/pgstore/migrations/tern.conf

sqlc_generate:
	sqlc generate -f ./internal/store/pgstore/sqlc.yml

start:
	air --build.cmd "go build -o ./bin/api ./cmd/api" --build.bin "./bin/api" 

