-- +migrate Up
CREATE SCHEMA IF NOT EXISTS app;

CREATE TABLE IF NOT EXISTS app.user (
    id uuid NOT NULL,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NOT NULL,

    username varchar(255) NOT NULL,
    email varchar(255),

    CONSTRAINT user_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS app.auth_providers (
    id uuid NOT NULL,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NOT NULL,

    user_id UUID NOT NULL REFERENCES app.user(id) ON DELETE CASCADE,
    method VARCHAR(50) NOT NULL,
    provider_user_id VARCHAR(255),
    password_hash VARCHAR(255),

    CONSTRAINT unique_provider_user UNIQUE (method, provider_user_id)
);

-- +migrate Down
DROP SCHEMA app cascade;