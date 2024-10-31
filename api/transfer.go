package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	
	"github.com/golang-jwt/jwt/v5"
	db "github.com/katatrina/simplebank/db/sqlc"
	"github.com/labstack/echo/v4"
)

type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" validate:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" validate:"required,min=1"`
	Amount        int64  `json:"amount" validate:"required,gt=0"`
	Currency      string `json:"currency" validate:"required,currency"`
}

func (server *Server) createTransfer(ctx echo.Context) error {
	req := new(createTransferRequest)
	
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	
	authPayload := ctx.Get(authorizationPayloadKey).(*jwt.RegisteredClaims)
	
	fromAccount, valid, err := server.validAccount(ctx, req.FromAccountID, req.Currency)
	if !valid {
		return err
	}
	
	if fromAccount.Owner != authPayload.Subject {
		err = errors.New("from account does not belong to the authenticated user")
		return ctx.JSON(http.StatusForbidden, errorResponse(err))
	}
	
	_, valid, err = server.validAccount(ctx, req.ToAccountID, req.Currency)
	if !valid {
		return err
	}
	
	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	
	result, err := server.store.TransferTx(context.Background(), arg)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	
	return ctx.JSON(http.StatusOK, result)
}

// validAccount checks if an account exists and has desired currency.
// It has also handled the response for the caller.
func (server *Server) validAccount(ctx echo.Context, accountID int64, currency string) (db.Account, bool, error) {
	account, err := server.store.GetAccount(context.Background(), accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return account, false, ctx.JSON(http.StatusNotFound, errorResponse(err))
		}
		
		return account, false, ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	
	if account.Currency != currency {
		err = fmt.Errorf("account %d currency mismatch: %s vs %s", accountID, account.Currency, currency)
		return account, false, ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	
	return account, true, nil
}
