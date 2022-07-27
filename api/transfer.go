package api

import (
	db "bank-mvp/db/sqlc"
	"database/sql"
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

	if req.SourceAccountID != 0 && !server.validAccount(ctx, req.SourceAccountID, req.Currency) {
		return
	}

	if req.DestAccountID != 0 && !server.validAccount(ctx, req.DestAccountID, req.Currency) {
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

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if account.Currency != currency {
		err = fmt.Errorf("currency mismatch for account ID %d: %s vs %s", accountID, account.Currency, currency)
		return false
	}

	return true
}