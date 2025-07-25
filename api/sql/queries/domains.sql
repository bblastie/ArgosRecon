-- name: InsertDomain :one
INSERT INTO domains (
    id, 
    name,
    created_at,
    updated_at
)
VALUES (
    gen_random_uuid(), $1, NOW(), NOW()
)
RETURNING *;