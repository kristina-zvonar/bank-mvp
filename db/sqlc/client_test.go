package db

import (
	"bank-mvp/util"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomClient(t *testing.T) Client {
	user := createRandomUser(t)

	arg := CreateClientParams {
		FirstName: util.RandomString(10),
		LastName: util.RandomString(10),
		CountryID: util.RandomInt(1, 228),
		UserID: user.ID,
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

func TestListClients(t *testing.T) {
	clientCnt := util.RandomInt(1, 20)	
	for i := 1; i < int(clientCnt); i++ {
		createRandomClient(t)
	}

	limit := util.RandomInt(1, 20)
	arg := ListClientsParams {
		Limit: int32(limit),
		Offset: 0,
	}

	clients, err := testQueries.ListClients(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, clients)
	if(clientCnt > limit) {
		require.True(t, len(clients) == int(limit))
	} else {
		require.True(t, len(clients) >= int(clientCnt))
	}

	for _, client := range clients {
		require.NotEmpty(t, client)		
	}
}

func TestCreateClient(t *testing.T) {
	createRandomClient(t)
}

func TestUpdateClient(t *testing.T) {
	client1 := createRandomClient(t)

	arg := UpdateClientParams {
		ID: client1.ID,
		FirstName: util.RandomString(10),
		LastName: util.RandomString(10),
		CountryID: util.RandomInt(1, 228),
		Active: false,
	}

	client2, err := testQueries.UpdateClient(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, client2)

	require.Equal(t, client1.ID, client2.ID)
	require.Equal(t, arg.FirstName, client2.FirstName)
	require.Equal(t, arg.LastName, client2.LastName)
	require.Equal(t, arg.CountryID, client2.CountryID)
	require.False(t, client2.Active)
}
