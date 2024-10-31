package api

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"
	
	"github.com/google/uuid"
	db "github.com/katatrina/simplebank/db/sqlc"
	"github.com/katatrina/simplebank/util"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

func (server *Server) createUser(ctx echo.Context) error {
	req := new(createUserRequest)
	
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	
	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}
	
	user, err := server.store.CreateUser(context.Background(), arg)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return ctx.JSON(http.StatusUnprocessableEntity, errorResponse(err))
			}
		}
		
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	
	return ctx.JSON(http.StatusOK, user)
}

type loginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginUserResponse struct {
	SessionID             uuid.UUID `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	User                  db.User   `json:"user"`
}

func (server *Server) loginUser(ctx echo.Context) error {
	req := new(loginUserRequest)
	
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	
	user, err := server.store.GetUser(context.Background(), req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ctx.JSON(http.StatusBadRequest, errorResponse(err))
		}
		
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	
	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	}
	
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	
	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.RefreshTokenDuration)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	
	sessionID, err := uuid.Parse(refreshPayload.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	
	session, err := server.store.CreateSession(context.Background(), db.CreateSessionParams{
		ID:           sessionID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request().UserAgent(),
		ClientIp:     ctx.RealIP(),
		ExpiresAt:    refreshPayload.ExpiresAt.Time,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	
	resp := loginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiresAt.Time,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiresAt.Time,
		User:                  user,
	}
	return ctx.JSON(http.StatusOK, resp)
}
