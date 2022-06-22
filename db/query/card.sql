-- name: GetCard :one
SELECT * FROM cards
WHERE 
    number = $1 AND valid_through = $2 and cvc = $3 and active = true
LIMIT 1;

-- name: ListCards :many
SELECT * FROM cards
LIMIT $1
OFFSET $2;

-- name: CreateCard :one
INSERT INTO cards
(
    number,
    valid_through,
    cvc,
    account_id
)
VALUES($1, $2, $3, $4)
RETURNING *;

-- name: UpdateCard :one
UPDATE cards
SET 
    active = $2
WHERE id = $1
RETURNING *;