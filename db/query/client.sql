-- name: GetClient :one
SELECT * FROM clients
WHERE id = $1 LIMIT 1;

-- name: ListClients :many
SELECT * FROM clients
LIMIT $1
OFFSET $2;

-- name: CreateClient :one
INSERT INTO clients
(
    first_name,
    last_name,
    country_id,
    user_id
)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateClient :one
UPDATE clients
SET 
    first_name = $2,
    last_name = $3,
    country_id = $4,
    active = $5
WHERE id = $1
RETURNING *;