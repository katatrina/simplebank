package api

import (
	"database/sql"
	"errors"
	"fmt"
	
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	db "github.com/katatrina/simplebank/db/sqlc"
)

type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" validate:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" validate:"required,min=1"`
	Amount        int64  `json:"amount" validate:"required,gt=0"`
	Currency      string `json:"currency" validate:"required,currency"`
}

func (server *Server) createTransfer(ctx *fiber.Ctx) error {
	req := new(createTransferRequest)
	
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}
	
	authPayload := ctx.UserContext().Value(authorizationPayloadKey).(*jwt.RegisteredClaims)
	
	fromAccount, err := server.validAccount(ctx, req.FromAccountID, req.Currency)
	if err != nil {
		return err
	}
	
	if fromAccount.Owner != authPayload.Subject {
		err = errors.New("from account does not belong to the authenticated user")
		return ctx.Status(fiber.StatusForbidden).JSON(errorResponse(err))
	}
	
	_, err = server.validAccount(ctx, req.ToAccountID, req.Currency)
	if err != nil {
		return err
	}
	
	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	
	result, err := server.store.TransferTx(ctx.Context(), arg)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}
	
	return ctx.Status(fiber.StatusOK).JSON(result)
}

// validAccount checks if an account exists and has desired currency.
// It has also handled the response for the caller.
func (server *Server) validAccount(ctx *fiber.Ctx, accountID int64, currency string) (db.Account, error) {
	account, err := server.store.GetAccount(ctx.Context(), accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return account, ctx.Status(fiber.StatusNotFound).JSON(errorResponse(err))
		}
		
		return account, ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}
	
	if account.Currency != currency {
		err = fmt.Errorf("account %d currency mismatch: %s vs %s", accountID, account.Currency, currency)
		return account, ctx.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}
	
	return account, nil
}
