-- +migrate Up
CREATE TABLE IF NOT EXISTS user_auths
(
    id SERIAL NOT NULL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    access_token TEXT NOT NULL,
    refresh_token TEXT,
    id_token TEXT NOT NULL,
    access_token_expired_at TIMESTAMPTZ,
    refresh_token_expired_at TIMESTAMPTZ,
    id_token_expired_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ
);

-- +migrate Down
DROP TABLE IF EXISTS user_auths;
