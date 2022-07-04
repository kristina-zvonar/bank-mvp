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
	existed := make(map[int]bool)
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
		require.True(t, diff1.GreaterThan(decimal.Zero))
		require.True(t, diff1.Mod(amount).Equal(decimal.Zero)) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount
				
		quotient := diff1.Div(amount)
		require.True(t, quotient.GreaterThanOrEqual(decimal.NewFromInt(1)) && quotient.LessThanOrEqual(decimal.NewFromInt(n)))
		require.NotContains(t, existed, quotient)
		existed[int(quotient.IntPart())] = true
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

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	// run n concurrent transfer transactions
	n := 10
	amount := decimal.NewFromInt(10)

	errs := make(chan error)
	
	for i := 0; i < n; i++ {
		sourceAccountID := acc1.ID
		destAccountID := acc2.ID

		if i % 2 == 1 {
			sourceAccountID = acc2.ID
			destAccountID = acc1.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransactionTxParams{
				SourceAccountID: sql.NullInt64{
					Int64: sourceAccountID,
					Valid: true,
				},
				DestAccountID: sql.NullInt64{
					Int64: destAccountID,
					Valid: true,
				},
				Amount: amount,
			})

			errs <- err
		}()
	}

	// check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check final balance
	updatedAcc1, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)
	require.True(t, acc1.Balance.Equal(updatedAcc1.Balance))
	
	updatedAcc2, err := testQueries.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)
	require.True(t, acc2.Balance.Equal(updatedAcc2.Balance))
}