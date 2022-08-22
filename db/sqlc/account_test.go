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

func TestAddAccountBalance(t *testing.T) {
	account1 := createRandomAccount(t)
	arg := AddAccountBalanceParams{
		Amount: util.RandomDecimal(1, 100),
		ID: account1.ID,
	}

	account2, err := testQueries.AddAccountBalance(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.ClientID, account2.ClientID)
	require.Equal(t, account1.Active, account2.Active)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.Balance.Add(arg.Amount), account2.Balance)
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.ClientID, account2.ClientID)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.Active, account2.Active)
	require.Equal(t, account1.Locked, account2.Locked)
}

func TestListAccount(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	arg := ListAccountsParams{
		ClientID: lastAccount.ClientID,
		Limit: 5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.ClientID, account.ClientID)
	}
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID: account1.ID,
		Balance: account1.Balance,
		Active: false,
		Locked: true,
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.ClientID, account2.ClientID)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, arg.Active, account2.Active)
	require.Equal(t, arg.Locked, account2.Locked)
}