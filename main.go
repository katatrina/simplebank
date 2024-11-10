package main

import (
	"database/sql"
	"os"
	
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	
	"github.com/katatrina/simplebank/api"
	db "github.com/katatrina/simplebank/db/sqlc"
	"github.com/katatrina/simplebank/util"
	
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig("./app.env")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config file")
	}
	
	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	
	connPool, err := sql.Open(config.DriverName, config.DataSourceName)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to validate db connection")
	}
	
	pingErr := connPool.Ping()
	if pingErr != nil {
		log.Fatal().Err(pingErr).Msg("failed to connect to db")
	}
	
	store := db.NewStore(connPool)
	
	runHTTPServer(config, store)
}

func runHTTPServer(config util.Config, store db.Store) {
	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create HTTP server")
	}
	
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start HTTP server")
	}
}

// func runGrpcServer(config util.Config, store db.Store) {
// 	server, err := gapi.NewServer(store, config)
// 	if err != nil {
// 		log.Fatalf("cannot create server <= %v", err)
// 	}
//
// 	grpcServer := grpc.NewServer()
// 	pb.RegisterSimpleBankServer(grpcServer, server)
// 	reflection.Register(grpcServer)
//
// 	listener, err := net.Listen("tcp", config.GRPCServerAddress)
// 	if err != nil {
// 		log.Fatalf("cannot create listener <= %v", err)
// 	}
//
// 	log.Printf("start gRPC server at %s", listener.Addr().String())
// 	err = grpcServer.Serve(listener)
// 	if err != nil {
// 		log.Fatalf("cannot start gRPC server <= %v", err)
// 	}
// }
