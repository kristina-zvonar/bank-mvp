// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package db

import (
	"context"
)

type Querier interface {
	AddAccountBalance(ctx context.Context, arg AddAccountBalanceParams) (Account, error)
	CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error)
	CreateCard(ctx context.Context, arg CreateCardParams) (Card, error)
	CreateClient(ctx context.Context, arg CreateClientParams) (Client, error)
	CreateService(ctx context.Context, arg CreateServiceParams) (Service, error)
	CreateTransaction(ctx context.Context, arg CreateTransactionParams) (Transaction, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteService(ctx context.Context, id int64) error
	GetAccount(ctx context.Context, id int64) (Account, error)
	GetAccountForUpdate(ctx context.Context, id int64) (Account, error)
	GetCard(ctx context.Context, arg GetCardParams) (Card, error)
	GetClient(ctx context.Context, id int64) (Client, error)
	GetService(ctx context.Context, id int64) (Service, error)
	GetTransaction(ctx context.Context, id int64) (Transaction, error)
	GetUser(ctx context.Context, id int64) (User, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)
	ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error)
	ListCards(ctx context.Context, arg ListCardsParams) ([]Card, error)
	ListClients(ctx context.Context, arg ListClientsParams) ([]Client, error)
	ListServices(ctx context.Context, arg ListServicesParams) ([]Service, error)
	ListTransactions(ctx context.Context, arg ListTransactionsParams) ([]Transaction, error)
	UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error)
	UpdateCard(ctx context.Context, arg UpdateCardParams) (Card, error)
	UpdateClient(ctx context.Context, arg UpdateClientParams) (Client, error)
	UpdateService(ctx context.Context, arg UpdateServiceParams) (Service, error)
}

var _ Querier = (*Queries)(nil)
