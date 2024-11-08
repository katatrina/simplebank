package api

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	db "github.com/katatrina/simplebank/db/sqlc"
	"github.com/lib/pq"
)

type createAccountRequest struct {
	Currency string `json:"currency"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	
	// Binding the request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	
	// Extracting JWT claims
	authPayload := ctx.MustGet(authorizationPayloadKey).(*jwt.RegisteredClaims)
	
	arg := db.CreateAccountParams{
		Owner:    authPayload.Subject,
		Balance:  0,
		Currency: req.Currency,
	}
	
	// Creating an account
	account, err := server.store.CreateAccount(context.Background(), arg)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code.Name() {
			case "unique_violation", "foreign_key_violation":
				ctx.JSON(http.StatusUnprocessableEntity, errorResponse(err))
				return
			}
		}
		
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	
	// Binding URI parameter
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	
	account, err := server.store.GetAccount(context.Background(), req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	// Verifying account ownership
	authPayload := ctx.MustGet(authorizationPayloadKey).(*jwt.RegisteredClaims)
	if account.Owner != authPayload.Subject {
		err = errors.New("requested account does not belong to the authenticated user")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	
	ctx.JSON(http.StatusOK, account)
}

type listAccountsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountsRequest
	
	// Binding query parameters
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	
	authPayload := ctx.MustGet(authorizationPayloadKey).(*jwt.RegisteredClaims)
	
	arg := db.ListAccountsByOwnerParams{
		Owner:  authPayload.Subject,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	
	// Listing accounts by owner
	accounts, err := server.store.ListAccountsByOwner(context.Background(), arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	ctx.JSON(http.StatusOK, accounts)
}
