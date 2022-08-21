-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1 
FOR NO KEY UPDATE;

-- name: ListAccounts :many
SELECT * FROM accounts
WHERE client_id = $1
LIMIT $2
OFFSET $3;

-- name: CreateAccount :one
INSERT INTO accounts
(
    balance,
    currency,
    client_id
)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateAccount :one
UPDATE accounts
SET
    balance = $2,
    active = $3,
    locked = $4
WHERE id = $1
RETURNING *;

-- name: AddAccountBalance :one
UPDATE accounts
SET 
    balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;