package api

import (
	db "bank-mvp/db/sqlc"
	"bank-mvp/token"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type createTransactionRequest struct {
	SourceAccountID int64 `json:"source_account_id"`
	DestAccountID   int64 `json:"dest_account_id"`
	Amount          decimal.Decimal `json:"amount" binding:"required,gt=0"`
	Currency 		string `json:"currency" binding:"currency"`
}

func (server *Server) createTransaction(ctx *gin.Context) {
	var req createTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	sourceAccount, valid := server.validAccount(ctx, req.SourceAccountID, req.Currency)
	if req.SourceAccountID != 0 && !valid {
		return
	}

	_, valid = server.validAccount(ctx, req.DestAccountID, req.Currency)
	if req.DestAccountID != 0 && !valid {
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	clientID, err := server.getClientID(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if req.SourceAccountID != 0 && (sourceAccount.ClientID != clientID) {
		err := errors.New("from account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	// TODO: find a way better way to do this
	arg := db.TransactionTxParams {
		SourceAccountID: sql.NullInt64{
			Int64: req.SourceAccountID,
			Valid: true,
		},
		DestAccountID: sql.NullInt64{
			Int64: req.DestAccountID,
			Valid: true,
		},
		Amount: req.Amount,
	}

	transaction, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse((err)))
		return
	}

	ctx.JSON(http.StatusOK, transaction)
}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return account, false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err = fmt.Errorf("currency mismatch for account ID %d: %s vs %s", accountID, account.Currency, currency)
		return account, false
	}

	return account, true
}