package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	db "github.com/katatrina/simplebank/db/sqlc"
)

type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTransferRequest
	
	// Binding the JSON body
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	
	authPayload := ctx.MustGet(authorizationPayloadKey).(*jwt.RegisteredClaims)
	
	fromAccount, valid := server.validAccount(ctx, req.FromAccountID, req.Currency)
	if !valid {
		return
	}
	
	// Verify account ownership
	if fromAccount.Owner != authPayload.Subject {
		err := errors.New("from account does not belong to the authenticated user")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	
	// Check the to-account validity
	_, valid = server.validAccount(ctx, req.ToAccountID, req.Currency)
	if !valid {
		return
	}
	
	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	
	// Perform the transfer
	result, err := server.store.TransferTx(context.Background(), arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	ctx.JSON(http.StatusOK, result)
}

// validAccount checks if an account exists and has the desired currency.
// It also handles any errors in the response for the caller.
func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(context.Background(), accountID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}
		
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}
	
	// Check currency match
	if account.Currency != currency {
		err = fmt.Errorf("account %d currency mismatch: %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}
	
	return account, true
}
