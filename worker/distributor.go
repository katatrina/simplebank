package worker

import (
	"context"
	
	"github.com/hibiken/asynq"
)

/*
This file will contain the codes to create tasks and distributes them to the Redis queue.
*/

type TaskDistributor interface {
	DistributeTaskSendVerifyEmail(
		ctx context.Context,
		payload *PayloadSendVerifyEmail,
		opts ...asynq.Option,
	) error
	Ping() error
}

type RedisTaskDistributor struct {
	client *asynq.Client // client sends tasks to redis queue.
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt)
	
	return &RedisTaskDistributor{
		client: client,
	}
}

// Ping checks if the Redis connection is alive.
func (distributor *RedisTaskDistributor) Ping() error {
	return distributor.client.Ping()
}
