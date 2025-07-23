include .env
export $(shell sed 's/=.*//' .env)

migration_run:
	tern migrate --migrations ./internal/store/pgstore/migrations --config ./internal/store/pgstore/migrations/tern.conf

migration_create:
	cd ./internal/store/pgstore/migrations/ && tern new  $(name)

sqlc_generate:
	sqlc generate -f ./internal/store/pgstore/sqlc.yml

start:
	air --build.cmd "go build -o ./bin/api ./cmd/api" --build.bin "./bin/api" 

