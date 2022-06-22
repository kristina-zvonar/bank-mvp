-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: ListAccounts :many
SELECT * FROM accounts
LIMIT $1
OFFSET $2;

-- name: CreateAccount :one
INSERT INTO accounts
(
    balance,
    currency,
    client_id
)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateAccount :many
UPDATE accounts
SET
    balance = $2,
    active = $3,
    locked = $4
WHERE id = $1
RETURNING *;