package api

import (
	"database/sql"
	"errors"
	"time"
	
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewAccessToken(ctx *fiber.Ctx) error {
	req := new(renewAccessTokenRequest)
	
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}
	
	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(errorResponse(err))
	}
	
	sessionID, err := uuid.Parse(refreshPayload.ID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}
	
	session, err := server.store.GetSession(ctx.Context(), sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ctx.Status(fiber.StatusNotFound).JSON(errorResponse(err))
		}
		
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}
	
	if session.IsBlocked {
		err = errors.New("blocked session")
		return ctx.Status(fiber.StatusUnauthorized).JSON(errorResponse(err))
	}
	
	if session.Username != refreshPayload.Subject {
		err = errors.New("mismatch session user")
		return ctx.Status(fiber.StatusUnauthorized).JSON(errorResponse(err))
	}
	
	if session.RefreshToken != req.RefreshToken {
		err = errors.New("mismatch refresh token")
		return ctx.Status(fiber.StatusUnauthorized).JSON(errorResponse(err))
	}
	
	if time.Now().After(refreshPayload.ExpiresAt.Time) {
		err = errors.New("expired session")
		return ctx.Status(fiber.StatusUnauthorized).JSON(errorResponse(err))
	}
	
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(session.Username, server.config.AccessTokenDuration)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}
	
	res := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiresAt.Time,
	}
	return ctx.Status(fiber.StatusOK).JSON(res)
}
