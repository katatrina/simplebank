package main

import (
	"database/sql"
	"os"
	
	"github.com/katatrina/simplebank/mail"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	
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
	
	redisOpt := asynq.RedisClientOpt{Addr: config.RedisServerAddress}
	
	taskDistributor, err := worker.NewRedisTaskDistributor(redisOpt)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create task distributor")
	}
	
	go runTaskProcessor(config, redisOpt, store)
	runHTTPServer(config, store, taskDistributor)
}

func runTaskProcessor(config util.Config, redisOpt asynq.RedisClientOpt, store db.Store) {
	mailer := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)
	log.Info().Msg("start task processor")
	
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}
}

func runHTTPServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := api.NewServer(store, config, taskDistributor)
	
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
