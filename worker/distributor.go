package worker

import (
	"context"
	"fmt"
	
	"github.com/hibiken/asynq"
)

/*
This file will contain the codes to create tasks and distributes them to the workers via Redis queue.
*/

type TaskDistributor interface {
	DistributeTaskSendVerifyEmail(
		ctx context.Context,
		payload *PayloadSendVerifyEmail,
		opts ...asynq.Option,
	) error
}

type RedisTaskDistributor struct {
	client *asynq.Client // client sends tasks to redis queue.
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) (TaskDistributor, error) {
	client := asynq.NewClient(redisOpt)
	err := client.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis server: %w", err)
	}
	
	return &RedisTaskDistributor{
		client: client,
	}, nil
}
