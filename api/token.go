package api

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"
	
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewAccessToken(ctx echo.Context) error {
	req := new(renewAccessTokenRequest)
	
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	
	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	}
	
	sessionID, err := uuid.Parse(refreshPayload.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	
	session, err := server.store.GetSession(context.Background(), sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ctx.JSON(http.StatusNotFound, errorResponse(err))
		}
		
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	
	if session.IsBlocked {
		err = errors.New("blocked session")
		return ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	}
	
	if session.Username != refreshPayload.Subject {
		err = errors.New("mismatch session user")
		return ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	}
	
	if session.RefreshToken != req.RefreshToken {
		err = errors.New("mismatch refresh token")
		return ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	}
	
	if time.Now().After(refreshPayload.ExpiresAt.Time) {
		err = errors.New("expired session")
		return ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	}
	
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(session.Username, server.config.AccessTokenDuration)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	
	res := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiresAt.Time,
	}
	return ctx.JSON(http.StatusOK, res)
}
