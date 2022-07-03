package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
)

// Store provides all functions to execute DB queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store {
		db: db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb err: %v", err, rbErr)
		}

		return err
	}

	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transaction DB TX
type TransactionTxParams struct {
	SourceAccountID    sql.NullInt64  `json:"source_account_id"`
	DestAccountID      sql.NullInt64  `json:"dest_account_id"`
	ExtSourceAccountID sql.NullString `json:"ext_source_account_id"`
	ExtDestAccountID   sql.NullString `json:"ext_dest_account_id"`
	Amount             decimal.Decimal`json:"amount"`
}

// TransactionTxResult is the result of the bank transaction
type TransactionTxResult struct {
	Transaction Transaction `json:"transaction"`
	SourceAccount Account `json:"src_account"`
	DestAccount Account `json:"dest_account"`
	ExtSourceAccount string `json:"ext_source_acc"`
	ExtDestAccount string `json:"ext_dest_acc"`
}

var txKey = struct{}{}

// TransferTx performs a money transfer from one account to the other
// It creates a transfer record and update accounts' balance within a single database transaction
func (store *Store) TransferTx(ctx context.Context, arg TransactionTxParams) (TransactionTxResult, error) {
	var result TransactionTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transaction, err = q.CreateTransaction(ctx, CreateTransactionParams{
			SourceAccountID: arg.SourceAccountID,
			DestAccountID: arg.DestAccountID,
			ExtSourceAccountID: arg.ExtSourceAccountID,
			ExtDestAccountID: arg.ExtDestAccountID,
			Amount: arg.Amount,
		})
		
		// get and update account's or accounts' balance
		if arg.SourceAccountID.Valid {
			sourceAccount, err := q.GetAccountForUpdate(ctx, arg.SourceAccountID.Int64)
			if err != nil {
				return err
			}

			result.SourceAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID: sourceAccount.ID,
				Balance: sourceAccount.Balance.Sub(arg.Amount),
			})
			if err != nil {
				return err
			}
		} else if(len(arg.ExtDestAccountID.String) > 0) {
			result.ExtSourceAccount = arg.ExtSourceAccountID.String
		} else {
			return errors.New("error: must provide either internal or external bank account ID")
		}
		
		if arg.DestAccountID.Valid {
			destAccount, err := q.GetAccountForUpdate(ctx, arg.DestAccountID.Int64)
			if err != nil {
				return err
			}

			result.DestAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID: destAccount.ID,
				Balance: destAccount.Balance.Add(arg.Amount),
			})
			if err != nil {
				return err
			}
		} else if(len(arg.ExtDestAccountID.String) > 0) {
			result.ExtDestAccount = arg.ExtDestAccountID.String
		} else {
			return errors.New("error: must provide either internal or external bank account ID")
		}

		if err != nil {
			return err
		}
		return nil
	})

	return result, err	
}