-- +migrate Up
CREATE TABLE IF NOT EXISTS app.password_reset (
    id uuid NOT NULL,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NOT NULL,

    user_id UUID NOT NULL REFERENCES app.user(id) ON DELETE CASCADE,  -- Links to user
    
    expires_at timestamptz NOT NULL,  -- Expiry time (e.g., 15 minutes from request)
    reset_token VARCHAR(512) NOT NULL UNIQUE,  -- Secure token for password reset
    used BOOLEAN NOT NULL DEFAULT FALSE,  -- Marks if the token has been used

    CONSTRAINT password_reset_pkey PRIMARY KEY (id)
);

-- +migrate Down
DROP TABLE IF EXISTS app.password_reset;