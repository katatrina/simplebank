package token

import (
	"time"
)

type Maker interface {
	// CreateToken creates a new token string for the specific username and duration.
	CreateToken(username string, role string, duration time.Duration) (token string, payload *Payload, err error)
	// VerifyToken verifies the token string and returns the payload if the token is valid.
	VerifyToken(tokenString string) (payload *Payload, err error)
}
