DB_URL=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable

postgres:
	docker run --name postgres16 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

createdb:
	docker exec -it postgres16 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres16 dropdb --owner=root simeple_bank

migrate-up:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up

migrate-up-1:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up 1

migrate-down:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down

migrate-down-1:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down 1

MIGRATION_VERSION ?= 1

migrate-force:
	migrate -path db/migrations -database "$(DB_URL)" -verbose force "$(MIGRATION_VERSION)"

test:
	go test -v -cover ./...

sqlc:
	sqlc generate

db-docs:
	dbdocs build doc/db.dbml

db-schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

server:
	go run main.go
