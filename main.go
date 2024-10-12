package main

import (
	"database/sql"
	"log"
	
	"github.com/katatrina/simplebank/api"
	db "github.com/katatrina/simplebank/db/sqlc"
	"github.com/katatrina/simplebank/util"
	
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig("./app.env")
	if err != nil {
		log.Fatalf("cannot load config file: %v", err)
	}
	
	connPool, err := sql.Open(config.DriverName, config.DataSourceName)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}
	
	pingErr := connPool.Ping()
	if pingErr != nil {
		log.Fatalf("cannot connect to db: %v", pingErr)
	}
	
	store := db.NewStore(connPool)
	
	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatalf("cannot create server <= %v", err)
	}
	
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatalf("cannot start server: %v", err)
	}
}
