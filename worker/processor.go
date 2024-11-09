package worker

import (
	"context"
	
	"github.com/hibiken/asynq"
	db "github.com/katatrina/simplebank/db/sqlc"
)

/*
 This file contains code that will pick up the tasks from the Redis queue and process them.
*/

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(
		cxt context.Context,
		task *asynq.Task,
	) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) TaskProcessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{}, // Use predefined default configurations.
	)
	
	return &RedisTaskProcessor{
		server: server,
		store:  store,
	}
}

// Start registers the task handlers for the mux, attaches the mux to the asynq server, and starts the server.
func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	
	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)
	
	return processor.server.Start(mux)
}
