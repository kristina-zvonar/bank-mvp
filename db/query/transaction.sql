-- name: GetTransaction :one
SELECT * FROM transactions
WHERE id = $1 LIMIT 1;

-- name: ListTransactions :many
SELECT * FROM transactions
LIMIT $1 
OFFSET $2;

-- name: CreateTransaction :one
INSERT INTO transactions
(
    amount,
    source_account_id,
    dest_account_id,
    ext_source_account_id,
    ext_dest_account_id,
    category,
    service_id
)
VALUES($1, $2, $3, $4, $5, $6, $7)
RETURNING *;
