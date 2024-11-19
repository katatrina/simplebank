package api

import (
	"errors"
	"net/http"
	"strings"
	
	"github.com/gin-gonic/gin"
	"github.com/katatrina/simplebank/token"
)

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "Bearer"
	authorizationPayloadKey = "authPayload"
)

// authMiddleware requires the client to provide a valid access token.
func authMiddleware(tokenMaker token.Maker, accessibleRoles []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if authorizationHeader == "" {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		
		fields := strings.Split(authorizationHeader, " ")
		if len(fields) != 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		
		authorizationHeaderType := fields[0]
		if authorizationHeaderType != authorizationTypeBearer {
			err := errors.New("unsupported authorization header type")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		
		payload, err := tokenMaker.VerifyToken(fields[1])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		
		if !hasPermissions(payload.Role, accessibleRoles) {
			err = errors.New("permission denied")
			ctx.AbortWithStatusJSON(http.StatusForbidden, errorResponse(err))
			return
		}
		
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

func hasPermissions(userRole string, accessibleRoles []string) bool {
	for _, role := range accessibleRoles {
		if userRole == role {
			return true
		}
	}
	
	return false
}
