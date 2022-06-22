-- name: GetService :one
SELECT * FROM services
WHERE id = $1 LIMIT 1;

-- name: ListServices :many
SELECT * FROM services
LIMIT $1
OFFSET $2;

-- name: CreateService :one
INSERT INTO services
(
    name, 
    type
)
VALUES($1, $2)
RETURNING *;

-- name: UpdateService :one
UPDATE services
SET
    name = $2,
    type = $3
WHERE id = $1
RETURNING *;

-- name: DeleteService :exec
DELETE FROM services
WHERE id = $1;