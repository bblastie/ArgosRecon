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
ON CONFLICT (name) DO UPDATE SET
    updated_at = NOW()
RETURNING id, name, created_at, updated_at;

-- name: AllDomains :many
SELECT id, name, created_at, updated_at FROM domains;

-- name: LookupDomainByName :one
SELECT id, name, created_at, updated_at FROM domains WHERE name = $1;   

-- name: LookupDomainByID :one
SELECT id, name, created_at, updated_at FROM domains WHERE id = $1;

-- name: DeleteOneDomain :exec
DELETE FROM domains WHERE id = $1; 