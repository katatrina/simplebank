package api

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	
	"github.com/golang-jwt/jwt/v5"
	db "github.com/katatrina/simplebank/db/sqlc"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

type createAccountRequest struct {
	Currency string `json:"currency"`
}

func (server *Server) createAccount(ctx echo.Context) error {
	req := new(createAccountRequest)
	
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	
	authPayload := ctx.Get(authorizationPayloadKey).(*jwt.RegisteredClaims)
	
	arg := db.CreateAccountParams{
		Owner:    authPayload.Subject,
		Balance:  0,
		Currency: req.Currency,
	}
	
	account, err := server.store.CreateAccount(context.Background(), arg)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code.Name() {
			case "unique_violation", "foreign_key_violation":
				return ctx.JSON(http.StatusUnprocessableEntity, errorResponse(err))
			}
		}
		
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	
	return ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `param:"id"`
}

// getAccount returns an account of the authenticated user by the account ID.
func (server *Server) getAccount(ctx echo.Context) error {
	req := new(getAccountRequest)
	
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	
	account, err := server.store.GetAccount(context.Background(), req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ctx.JSON(http.StatusNotFound, errorResponse(err))
		}
		
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	
	authPayload := ctx.Get(authorizationPayloadKey).(*jwt.RegisteredClaims)
	if account.Owner != authPayload.Subject {
		err = errors.New("requested account does not belong to the authenticated user")
		return ctx.JSON(http.StatusForbidden, errorResponse(err))
	}
	
	return ctx.JSON(http.StatusOK, account)
}

type listAccountsRequest struct {
	PageID   int32 `query:"page_id"`
	PageSize int32 `query:"page_size"`
}

// listAccounts returns a list of all accounts of the authenticated user.
func (server *Server) listAccounts(ctx echo.Context) error {
	req := new(listAccountsRequest)
	
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	
	authPayload := ctx.Get(authorizationPayloadKey).(*jwt.RegisteredClaims)
	
	arg := db.ListAccountsByOwnerParams{
		Owner:  authPayload.Subject,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	
	accounts, err := server.store.ListAccountsByOwner(context.Background(), arg)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	
	return ctx.JSON(http.StatusOK, accounts)
}
