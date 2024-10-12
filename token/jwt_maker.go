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

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("cannot create token: %w", err)
	}
	
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    tokenIssuer,
			Subject:   username,
			Audience:  jwt.ClaimStrings{"client"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        tokenID.String(),
		},
	)
	
	signedToken, err := unsignedToken.SignedString([]byte(maker.secretKey))
	return signedToken, err
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
	
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, fmt.Errorf("unknown claims type, cannot proceed")
	}
	
	return claims, nil
}
