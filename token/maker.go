package token

import (
	"time"
	
	"github.com/golang-jwt/jwt/v5"
)

type Maker interface {
	// CreateToken creates a new token string for the specific username and duration.
	CreateToken(username string, duration time.Duration) (token string, payload *jwt.RegisteredClaims, err error)
	// VerifyToken verifies the token string and returns the payload if the token is valid.
	VerifyToken(tokenString string) (payload *jwt.RegisteredClaims, err error)
}
