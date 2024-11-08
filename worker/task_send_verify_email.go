package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	
	"github.com/hibiken/asynq"
)

const (
	TaskSendVerifyEmail = "email:verify"
)

// PayloadSendVerifyEmail contain all data of the task that we want to store in Redis.
type PayloadSendVerifyEmail struct {
	Username string
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}
	
	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	_, err = distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}
	
	// TODO: Log successful enqueue message
	
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(
	ctx context.Context,
	task *asynq.Task,
) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}
	
	_, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user does not exist: %w", asynq.SkipRetry)
		}
		
		return fmt.Errorf("failed to get user: %w", err)
	}
	
	// TODO: send email to user
	
	return nil
}
