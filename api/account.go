package api

import (
	"database/sql"
	"errors"
	
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	db "github.com/katatrina/simplebank/db/sqlc"
	"github.com/lib/pq"
)

type createAccountRequest struct {
	Currency string `json:"currency" validate:"required,currency"`
}

func (server *Server) createAccount(ctx *fiber.Ctx) error {
	req := new(createAccountRequest)
	
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}
	
	authPayload := ctx.UserContext().Value(authorizationPayloadKey).(*jwt.RegisteredClaims)
	
	arg := db.CreateAccountParams{
		Owner:    authPayload.Subject,
		Balance:  0,
		Currency: req.Currency,
	}
	
	account, err := server.store.CreateAccount(ctx.Context(), arg)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code.Name() {
			case "unique_violation", "foreign_key_violation":
				return ctx.Status(fiber.StatusUnprocessableEntity).JSON(errorResponse(err))
			}
		}
		
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}
	
	return ctx.Status(fiber.StatusOK).JSON(account)
}

type getAccountRequest struct {
	ID int64 `params:"id" validate:"required,min=1"`
}

// getAccount returns an account of the authenticated user by the account ID.
func (server *Server) getAccount(ctx *fiber.Ctx) error {
	req := new(getAccountRequest)
	
	if err := ctx.ParamsParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}
	
	account, err := server.store.GetAccount(ctx.Context(), req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ctx.Status(fiber.StatusNotFound).JSON(errorResponse(err))
		}
		
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}
	
	authPayload := ctx.UserContext().Value(authorizationPayloadKey).(*jwt.RegisteredClaims)
	if account.Owner != authPayload.Subject {
		err = errors.New("requested account does not belong to the authenticated user")
		return ctx.Status(fiber.StatusForbidden).JSON(errorResponse(err))
	}
	
	return ctx.Status(fiber.StatusOK).JSON(account)
}

type listAccountsRequest struct {
	PageID   int32 `query:"page_id" validate:"required,min=1"`
	PageSize int32 `query:"page_size" validate:"required,min=5,max=10"`
}

// listAccounts returns a list of all accounts of the authenticated user.
func (server *Server) listAccounts(ctx *fiber.Ctx) error {
	req := new(listAccountsRequest)
	
	if err := ctx.QueryParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}
	
	authPayload := ctx.UserContext().Value(authorizationPayloadKey).(*jwt.RegisteredClaims)
	
	arg := db.ListAccountsByOwnerParams{
		Owner:  authPayload.Subject,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	
	accounts, err := server.store.ListAccountsByOwner(ctx.Context(), arg)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}
	
	return ctx.Status(fiber.StatusOK).JSON(accounts)
}
