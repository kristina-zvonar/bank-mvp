package db

import (
	"bank-mvp/util"
	"context"
	"database/sql"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	// run n concurrent transfer transactions
	n := util.RandomInt(1, 10)
	amount := decimal.NewFromInt(10)

	errs := make(chan error)
	results := make(chan TransactionTxResult)

	for i := int64(0); i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransactionTxParams{
				SourceAccountID: sql.NullInt64{
					Int64: acc1.ID,
					Valid: true,
				},
				DestAccountID: sql.NullInt64{
					Int64: acc2.ID,
					Valid: true,
				},
				Amount: amount,
			})

			errs <- err
			results <- result
		}()
	}

	// check results
	for i := int64(0); i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <- results
		require.NotEmpty(t, result)

		// check bank transaction
		transaction := result.Transaction
		require.NotEmpty(t, transaction)
		require.Equal(t, acc1.ID, transaction.SourceAccountID.Int64)
		require.Equal(t, acc2.ID, transaction.DestAccountID.Int64)
		require.Equal(t, amount, transaction.Amount)
		require.NotZero(t, transaction.ID)

		_, err = store.GetTransaction(context.Background(), transaction.ID)
		require.NoError(t, err)

		// TODO check accounts'  balance
	}
}