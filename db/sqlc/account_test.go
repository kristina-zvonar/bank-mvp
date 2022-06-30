package db

import (
	"bank-mvp/util"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		ClientID: createRandomClient(t).ID,
		Balance:  util.RandomDecimal(1, 50),
		Currency: util.RandomString(3),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.ClientID, account.ClientID)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	return account
}