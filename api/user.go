package api

import (
	"database/sql"
	"errors"
	
	"github.com/gofiber/fiber/v2"
	db "github.com/katatrina/simplebank/db/sqlc"
	"github.com/katatrina/simplebank/util"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username string `json:"username" validate:"required,alphanum"` // alphanum: only letters and numbers
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

func (server *Server) createUser(ctx *fiber.Ctx) error {
	req := new(createUserRequest)
	
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}
	
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}
	
	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}
	
	user, err := server.store.CreateUser(ctx.Context(), arg)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return ctx.Status(fiber.StatusUnprocessableEntity).JSON(errorResponse(err))
			}
		}
		
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}
	
	return ctx.Status(fiber.StatusOK).JSON(user)
}

type loginUserRequest struct {
	Username string `json:"username" validate:"required,alphanum"`
	Password string `json:"password" validate:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken string  `json:"access_token"`
	User        db.User `json:"user"`
}

func (server *Server) loginUser(ctx *fiber.Ctx) error {
	req := new(loginUserRequest)
	
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}
	
	user, err := server.store.GetUser(ctx.Context(), req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
		}
		
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}
	
	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(errorResponse(err))
	}
	
	accessToken, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}
	
	res := loginUserResponse{
		AccessToken: accessToken,
		User:        user,
	}
	return ctx.Status(fiber.StatusOK).JSON(res)
}
