-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
    id SERIAL NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone_number VARCHAR(25) NOT NULL,
    encrypted_password VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT FALSE,
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL,
    created_by INT,
    updated_at TIMESTAMPTZ NOT NULL,
    updated_by INT,
    deleted_at TIMESTAMPTZ,
    deleted_by INT
);

-- +migrate Down
DROP TABLE IF EXISTS users;
