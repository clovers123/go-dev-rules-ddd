CREATE SCHEMA IF NOT EXISTS {{.Schema}};

CREATE TABLE IF NOT EXISTS "{{.Schema}}"."{{.Table}}" (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    owned_by UUID NOT NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    create_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    update_time TIMESTAMP WITHOUT TIME ZONE,
    delete_time TIMESTAMP WITHOUT TIME ZONE,
    created_by UUID,
    updated_by UUID,
    deleted_by UUID,
    feature JSONB NOT NULL DEFAULT '{}'::jsonb
);

CREATE INDEX IF NOT EXISTS idx_{{.Table}}_owned_by_deleted
    ON "{{.Schema}}"."{{.Table}}" (owned_by, is_deleted);
