-- +migrate Up
CREATE TABLE IF NOT EXISTS app_settings (
    id SERIAL NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    value TEXT NOT NULL,
    slug VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    created_by INT,
    updated_at TIMESTAMPTZ NOT NULL,
    updated_by INT,
    deleted_at TIMESTAMPTZ,
    deleted_by INT
);

-- +migrate Down
DROP TABLE IF EXISTS app_settings;
