package db

import (
	"bank-mvp/util"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomClient(t *testing.T) Client {
	arg := CreateClientParams {
		FirstName: util.RandomString(10),
		LastName: util.RandomString(10),
		CountryID: util.RandomInt(1, 229),
	}

	client, err := testQueries.CreateClient(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, client)

	require.Equal(t, arg.FirstName, client.FirstName)
	require.Equal(t, arg.LastName, client.LastName)
	require.Equal(t, arg.CountryID, client.CountryID)

	require.NotZero(t, client.ID)
	require.True(t, client.Active)
	require.NotZero(t, client.CreatedAt)

	return client
}

func TestGetClient(t *testing.T) {
	client1 := createRandomClient(t)	
	client2, err := testQueries.GetClient(context.Background(), client1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, client2)

	require.Equal(t, client1.ID, client2.ID)
	require.Equal(t, client1.FirstName, client2.FirstName)
	require.Equal(t, client1.LastName, client2.LastName)
	require.Equal(t, client1.CountryID, client2.CountryID)
	require.Equal(t, client1.Active, client2.Active)
}

