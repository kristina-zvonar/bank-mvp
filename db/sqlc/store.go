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
		
		if(arg.SourceAccountID.Int64 < arg.DestAccountID.Int64) {		
			result.SourceAccount, result.DestAccount, result.ExtSourceAccount, result.ExtDestAccount, err = addMoney(ctx, q, arg.SourceAccountID, arg.ExtSourceAccountID, arg.Amount.Neg(), arg.DestAccountID, arg.ExtDestAccountID, arg.Amount)
		} else {
			result.DestAccount, result.SourceAccount, result.ExtDestAccount, result.ExtSourceAccount, err = addMoney(ctx, q, arg.DestAccountID, arg.ExtDestAccountID, arg.Amount, arg.SourceAccountID, arg.ExtDestAccountID, arg.Amount.Neg())			
		}

		if err != nil {
			return err
		}
		return nil
	})

	return result, err	
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 sql.NullInt64,
	extAccountID1 sql.NullString,
	amount1 decimal.Decimal,
	accountID2 sql.NullInt64,
	extAccountID2 sql.NullString,
	amount2 decimal.Decimal,
) (account1 Account, account2 Account, extAccount1 string, extAccount2 string, err error) {
	if accountID1.Valid {
		account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID: accountID1.Int64,
			Amount: amount1,
		})
		if err != nil {
			return
		}
	} else if(len(extAccountID1.String) > 0) {
		extAccount1 = extAccountID1.String
	} else {
		err = errors.New("error: must provide either internal or external bank account ID")		
	}

	if accountID2.Valid {
		account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID: accountID2.Int64,
			Amount: amount2,
		})
		if err != nil {
			return
		}
	} else if(len(extAccountID2.String) > 0) {
		extAccount2 = extAccountID2.String
	} else {
		err = errors.New("error: must provide either internal or external bank account ID")		
	}

	return
}