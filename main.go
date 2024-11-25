package main

import (
	"context"
	"crypto/tls"
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
	// Load configurations
	config, err := util.LoadConfig("./app.env")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config file")
	}
	
	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	
	log.Info().Msg("configurations loaded successfully 👍")
	
	// Create connection pool
	connPool, err := pgxpool.New(context.Background(), config.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to validate db connection string")
	}
	
	pingErr := connPool.Ping(context.Background())
	if pingErr != nil {
		log.Fatal().Err(pingErr).Msg("failed to connect to db")
	}
	log.Info().Msg("connected to db 👍")
	
	store := db.NewStore(connPool)
	
	mailer, err := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	if err != nil {
		log.Err(err).Msg("failed to establish our email client")
	}
	
	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisServerAddress,
	}
	if config.Environment == "production" {
		redisOpt.Password = config.RedisServerPassword
		redisOpt.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}
	
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	
	go runTaskProcessor(redisOpt, store, mailer)
	runHTTPServer(config, store, taskDistributor, mailer)
}

// runTaskProcessor creates a new task processor and starts it.
func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, mailer mail.EmailSender) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)
	
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}
	
	log.Info().Msg("task processor started 👍")
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
