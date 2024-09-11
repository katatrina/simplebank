package main

import (
	"database/sql"
	"log"

	"github.com/katatrina/simplebank/api"
	db "github.com/katatrina/simplebank/db/sqlc"

	_ "github.com/lib/pq"
)

const (
	driverName     = "postgres"
	dataSourceName = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"

	serverAddress = "0.0.0.0:8080"
)

func main() {
	connPool, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	pingErr := connPool.Ping()
	if pingErr != nil {
		log.Fatalf("cannot connect to db: %v", pingErr)
	}

	store := db.NewStore(connPool)

	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatalf("cannot start server: %v", err)
	}
}
