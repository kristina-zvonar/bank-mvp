// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package db

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

type Account struct {
	ID        int64           `json:"id"`
	Balance   decimal.Decimal `json:"balance"`
	Currency  string          `json:"currency"`
	Active    bool            `json:"active"`
	Locked    bool            `json:"locked"`
	CreatedAt time.Time       `json:"created_at"`
	ClientID  int64           `json:"client_id"`
}

type Bank struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CountryID int64  `json:"country_id"`
}

type Card struct {
	ID           int64     `json:"id"`
	Number       string    `json:"number"`
	ValidThrough time.Time `json:"valid_through"`
	Cvc          string    `json:"cvc"`
	Active       bool      `json:"active"`
	AccountID    int64     `json:"account_id"`
}

type Client struct {
	ID        int64     `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	Active    bool      `json:"active"`
	CountryID int64     `json:"country_id"`
	UserID    int64     `json:"user_id"`
}

type Country struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Service struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type int32  `json:"type"`
}

type Transaction struct {
	ID                 int64           `json:"id"`
	Amount             decimal.Decimal `json:"amount"`
	SourceAccountID    sql.NullInt64   `json:"source_account_id"`
	DestAccountID      sql.NullInt64   `json:"dest_account_id"`
	ExtSourceAccountID sql.NullString  `json:"ext_source_account_id"`
	ExtDestAccountID   sql.NullString  `json:"ext_dest_account_id"`
	Category           int32           `json:"category"`
	ServiceID          sql.NullInt64   `json:"service_id"`
}

type User struct {
	ID                int64     `json:"id"`
	Username          string    `json:"username"`
	Password          string    `json:"password"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}
