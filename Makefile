include .env
export $(shell sed 's/=.*//' .env)

migration_run:
	tern migrate --migrations ./pgstore/migrations --config ./pgstore/migrations/tern.conf

migration_create:
	cd ./pgstore/migrations/ && tern new  $(name)

sqlc_generate:
	sqlc generate -f ./pgstore/sqlc.yml

start:
	podman compose up -d
	air

