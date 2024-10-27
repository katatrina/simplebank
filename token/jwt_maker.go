package token

import (
	"fmt"
	"time"
	
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	minSecretKeyLength int    = 32
	tokenIssuer        string = "Simplebank"
)

type JWTMaker struct {
	secretKey string
}

// NewJWTMaker creates a new JWTMaker with the specific secret key.
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeyLength {
		return nil, fmt.Errorf("invalid secret key size: must be at least %d characters", minSecretKeyLength)
	}
	
	return &JWTMaker{secretKey}, nil
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, *jwt.RegisteredClaims, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return "", nil, fmt.Errorf("cannot create token: %w", err)
	}
	
	payload := &jwt.RegisteredClaims{
		ID:        tokenID.String(),
		Issuer:    tokenIssuer,
		Subject:   username,
		Audience:  jwt.ClaimStrings{"client"},
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
	}
	
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	
	signedToken, err := unsignedToken.SignedString([]byte(maker.secretKey))
	return signedToken, payload, err
}

func (maker *JWTMaker) VerifyToken(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		
		return []byte(maker.secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	
	payload, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, fmt.Errorf("unknown claims type, cannot proceed")
	}
	
	return payload, nil
}
