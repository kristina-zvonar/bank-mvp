package db

import (
	"bank-mvp/util"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams {
		Username: util.RandomString(10),
		Password: hashedPassword,
		Email: util.RandomEmail(10),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Password, user.Password)
	require.Equal(t, arg.Email, user.Email)
	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())

	return user
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Password, user2.Password)
	require.Equal(t, user1.Email, user2.Email)
}

func TestGetUserByUsername(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUserByUsername(context.Background(), user1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Password, user2.Password)
	require.Equal(t, user1.Email, user2.Email)
}

func TestCreateUser(t *testing.T) {
	arg := CreateUserParams {
		Username: util.RandomString(10),
		Password: util.RandomString(10), // temp
		Email: util.RandomEmail(10),
	}

	result, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, arg.Username, result.Username)
	require.Equal(t, arg.Password, result.Password)
	require.Equal(t, arg.Email, result.Email)
}