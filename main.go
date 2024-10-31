package main

import (
	"database/sql"
	"log"
	"net"
	
	"github.com/katatrina/simplebank/api"
	db "github.com/katatrina/simplebank/db/sqlc"
	"github.com/katatrina/simplebank/gapi"
	"github.com/katatrina/simplebank/pb"
	"github.com/katatrina/simplebank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	
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
	
	// runGrpcServer(config, store)
	runHTTPServer(config, store)
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(store, config)
	if err != nil {
		log.Fatalf("cannot create server <= %v", err)
	}
	
	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)
	
	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatalf("cannot create listener <= %v", err)
	}
	
	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("cannot start gRPC server <= %v", err)
	}
}

func runHTTPServer(config util.Config, store db.Store) {
	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatalf("cannot create HTTP server <= %v", err)
	}
	
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatalf("cannot start HTTP server <= %v", err)
	}
}
