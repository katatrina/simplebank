package main

import (
	"database/sql"
	"log"
	
	"github.com/hibiken/asynq"
	"github.com/katatrina/simplebank/api"
	db "github.com/katatrina/simplebank/db/sqlc"
	"github.com/katatrina/simplebank/util"
	"github.com/katatrina/simplebank/worker"
	
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
	
	redisOpt := asynq.RedisClientOpt{Addr: config.RedisServerAddress}
	
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	go runTaskProcessor(redisOpt, store)
	runHTTPServer(config, store, taskDistributor)
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store)
	log.Println("start task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatalf("failed to start task processor <= %v", err)
	}
}

func runHTTPServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := api.NewServer(store, config, taskDistributor)
	if err != nil {
		log.Fatalf("cannot create HTTP server <= %v", err)
	}
	
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatalf("cannot start HTTP server <= %v", err)
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
