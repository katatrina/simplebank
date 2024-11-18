package api

import (
	"context"
	"time"
	
	"github.com/gin-gonic/gin"
)

type healthCheckResponse struct {
	Status    string            `json:"status"`
	Failures  map[string]string `json:"failures,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
}

func (server *Server) healthCheck(ctx *gin.Context) {
	failures := make(map[string]string)
	
	// Check database connection
	dbErr := server.store.Ping(ctx)
	if dbErr != nil {
		failures["postgres"] = dbErr.Error()
	}
	
	// Check Redis connection
	redisErr := server.taskDistributor.Ping()
	if redisErr != nil {
		failures["redis"] = redisErr.Error()
	}
	
	// Check SMTP connection
	smtpErr := server.mailer.Ping(context.Background())
	if smtpErr != nil {
		failures["gmail-smtp"] = smtpErr.Error()
	}
	
	status := "OK"
	if dbErr != nil {
		status = "Unavailable"
	} else if redisErr != nil || smtpErr != nil {
		status = "Partially Available"
	}
	
	resp := healthCheckResponse{
		Status:    status,
		Failures:  failures,
		Timestamp: time.Now(),
	}
	ctx.JSON(200, resp)
}
