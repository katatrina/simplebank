package api

import (
	"context"
	"errors"
	"strings"
	
	"github.com/gofiber/fiber/v2"
	"github.com/katatrina/simplebank/token"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authPayload"
)

// authMiddleware requires the client to provide a valid access token.
func authMiddleware(tokenMaker token.Maker) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authorizationHeader := ctx.Get(authorizationHeaderKey)
		if authorizationHeader == "" {
			err := errors.New("authorization header is not provided")
			return ctx.Status(fiber.StatusUnauthorized).JSON(errorResponse(err))
		}
		
		fields := strings.Split(authorizationHeader, " ")
		if len(fields) != 2 {
			err := errors.New("invalid authorization header format")
			return ctx.Status(fiber.StatusUnauthorized).JSON(errorResponse(err))
		}
		
		authorizationHeaderType := strings.ToLower(fields[0])
		if authorizationHeaderType != authorizationTypeBearer {
			err := errors.New("unsupported authorization header type")
			return ctx.Status(fiber.StatusUnauthorized).JSON(errorResponse(err))
		}
		
		payload, err := tokenMaker.VerifyToken(fields[1])
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(errorResponse(err))
		}
		
		userCtx := context.WithValue(context.Background(), authorizationPayloadKey, payload)
		ctx.SetUserContext(userCtx)
		
		return ctx.Next()
	}
}
