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

		// check accounts
		sourceAccount := result.SourceAccount
		require.NotEmpty(t, sourceAccount)
		require.Equal(t, acc1.ID, sourceAccount.ID)

		destAccount := result.DestAccount
		require.NotEmpty(t, destAccount)
		require.Equal(t, acc2.ID, destAccount.ID)

		// check accounts' balance
		diff1 := acc1.Balance.Sub(sourceAccount.Balance)
		diff2 := destAccount.Balance.Sub(acc2.Balance)

		require.Equal(t, diff1, diff2)
		require.Positive(t, diff1)
				
		quotient := diff1.Div(amount)
		require.True(t, quotient.GreaterThanOrEqual(decimal.NewFromInt(1)) && quotient.LessThanOrEqual(decimal.NewFromInt(n)))

	}

	// check final balance
	finalAmount := amount.Mul(decimal.NewFromInt(n))
	updatedAcc1, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)
	require.True(t, acc1.Balance.Sub(finalAmount).Equal(updatedAcc1.Balance))
	
	updatedAcc2, err := testQueries.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)
	require.True(t, acc2.Balance.Add(finalAmount).Equal(updatedAcc2.Balance))
}