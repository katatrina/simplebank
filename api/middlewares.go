package api

import (
	"errors"
	"net/http"
	"strings"
	
	"github.com/labstack/echo/v4"
)

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "Bearer"
	authorizationPayloadKey = "authPayload"
)

// authMiddleware requires the client to provide a valid access token.
func (server *Server) authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		authorizationHeader := ctx.Request().Header.Get(authorizationHeaderKey)
		if authorizationHeader == "" {
			err := errors.New("authorization header is not provided")
			return ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		}
		
		fields := strings.Split(authorizationHeader, " ")
		if len(fields) != 2 {
			err := errors.New("invalid authorization header format")
			return ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		}
		
		authorizationHeaderType := fields[0]
		if authorizationHeaderType != authorizationTypeBearer {
			err := errors.New("unsupported authorization header type")
			return ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		}
		
		payload, err := server.tokenMaker.VerifyToken(fields[1])
		if err != nil {
			return ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		}
		
		ctx.Set(authorizationPayloadKey, payload)
		
		return next(ctx)
	}
}
