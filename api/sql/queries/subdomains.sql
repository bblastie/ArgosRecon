-- name: InsertSubdomain :one
INSERT INTO subdomains (
    id,
    name,
    domain_id,
    created_at,
    updated_at
)
VALUES (
    gen_random_uuid(), $1, $2, NOW(), NOW()
)

ON CONFLICT (name) DO UPDATE SET
    updated_at = NOW()
RETURNING id, name, domain_id, created_at, updated_at;

-- name: LookupSubdomainByName :one
SELECT 
    s.id,
    s.name,
    s.domain_id,
    d.name as domain_name, 
    s.created_at,
    s.updated_at
FROM subdomains s 
JOIN domains d ON s.domain_id = d.id WHERE s.name = $1;

-- name: LookupSubdomainByID :one
SELECT * FROM subdomains WHERE id = $1;

-- name: DeleteSubdomainByID :exec
DELETE from subdomains WHERE id = $1; 