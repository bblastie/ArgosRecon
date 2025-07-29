-- +goose Up
CREATE TABLE subdomains(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE,
    domain_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_domain
    FOREIGN KEY (domain_id)
    REFERENCES domains(id)
    ON DELETE CASCADE 
);

-- +goose Down
DROP TABLES subdomains;