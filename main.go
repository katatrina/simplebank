package main

import (
	"context"
	"os"
	
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/katatrina/simplebank/mail"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	
	"github.com/hibiken/asynq"
	"github.com/katatrina/simplebank/api"
	db "github.com/katatrina/simplebank/db/sqlc"
	"github.com/katatrina/simplebank/util"
	
	"github.com/katatrina/simplebank/worker"
)

func main() {
	config, err := util.LoadConfig("./app.env")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config file")
	}
	
	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	
	connPool, err := pgxpool.New(context.Background(), config.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to validate db connection")
	}
	
	store := db.NewStore(connPool)
	
	mailer, err := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	if err != nil {
		log.Err(err).Msg("failed to establish our email client")
	}
	
	redisOpt := asynq.RedisClientOpt{Addr: config.RedisServerAddress}
	
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	
	go runTaskProcessor(redisOpt, store, mailer)
	runHTTPServer(config, store, taskDistributor, mailer)
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, mailer mail.EmailSender) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)
	log.Info().Msg("start task processor")
	
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}
}

func runHTTPServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor, mailer mail.EmailSender) {
	server, err := api.NewServer(store, config, taskDistributor, mailer)
	
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create HTTP server")
	}
	
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start HTTP server")
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
