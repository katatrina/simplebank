package worker

import (
	"context"
	
	"github.com/hibiken/asynq"
	db "github.com/katatrina/simplebank/db/sqlc"
	"github.com/katatrina/simplebank/mail"
	"github.com/rs/zerolog/log"
)

/*
 This file contains code that will pick up the tasks from the Redis queue and process them.
*/

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(
		cxt context.Context,
		task *asynq.Task,
	) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store         // store is used to access the database.
	mailer mail.EmailSender // mailer is used to send emails.
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, mailer mail.EmailSender) TaskProcessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error().Err(err).Str("type", task.Type()).
					Bytes("payload", task.Payload()).Msg("process task failed")
			}),
			Logger: NewLogger(),
		},
	)
	
	return &RedisTaskProcessor{
		server: server,
		store:  store,
		mailer: mailer,
	}
}

// Start registers the task handlers for the mux, attaches the mux to the asynq server, and starts the server.
func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	
	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)
	
	return processor.server.Start(mux)
}
