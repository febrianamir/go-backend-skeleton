-- +migrate Up
CREATE TABLE IF NOT EXISTS user_verifications (
    id SERIAL NOT NULL PRIMARY KEY,
    type VARCHAR(255),
    user_id INT NOT NULL,
    code VARCHAR(255) NOT NULL,
    expired_at TIMESTAMPTZ,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ
);

-- +migrate Down
DROP TABLE IF EXISTS user_verifications;
