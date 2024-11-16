package worker

import (
	"context"
	"encoding/json"
	"fmt"
	
	db "github.com/katatrina/simplebank/db/sqlc"
	"github.com/katatrina/simplebank/util"
	"github.com/rs/zerolog/log"
	
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
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}
	
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("task enqueued")
	
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
	
	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		// if errors.Is(err, db.ErrRecordNotFound) {
		// 	return fmt.Errorf("user does not exist: %w", asynq.SkipRetry)
		// }
		
		return fmt.Errorf("failed to get user: %w", err)
	}
	
	verifyEmail, err := processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: util.RandomString(32),
	})
	if err != nil {
		return fmt.Errorf("failed to create verify email: %w", err)
	}
	
	subject := "Welcome to SimpleBank"
	verifyURL := fmt.Sprintf("http://localhost:8080/v1/users/verify_email?email_id=%d&secret_code=%s", verifyEmail.ID, verifyEmail.SecretCode)
	body := fmt.Sprintf(`Hello %s,<br/>
Thank you for registering with us!<br/>
Please <a href="%s">click here</a> to verify your email address.<br/>`, user.FullName, verifyURL)
	to := []string{user.Email}
	err = processor.mailer.SendEmail(subject, body, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}
	
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("email", user.Email).Msg("task processed")
	
	return nil
}
