-- +migrate Up
CREATE TABLE IF NOT EXISTS app.session (
    id uuid NOT NULL,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NOT NULL,

    user_id UUID NOT NULL REFERENCES app.user(id) ON DELETE CASCADE,  -- Link to the user

    refresh_token VARCHAR(512) NOT NULL UNIQUE,           -- The refresh token (store as a hash)
    issued_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),         -- When the token was issued
    expires_at TIMESTAMPTZ NOT NULL,                      -- When the token expires
    revoked BOOLEAN NOT NULL DEFAULT FALSE,               -- If the token has been revoked

    user_agent VARCHAR(255),                              -- Optional: Track the device/browser

    CONSTRAINT unique_token UNIQUE (token),                -- Ensure tokens are unique
    CONSTRAINT refresh_token_pkey PRIMARY KEY (id)
);

-- +migrate Down
DROP TABLE IF EXISTS app.session;