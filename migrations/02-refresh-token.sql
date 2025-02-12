-- +migrate Up
CREATE TABLE IF NOT EXISTS app.session (
    id uuid NOT NULL,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NOT NULL,

    user_id UUID NOT NULL REFERENCES app.user(id) ON DELETE CASCADE,  

    refresh_token VARCHAR(512) NOT NULL UNIQUE,           
    issued_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),         
    expires_at TIMESTAMPTZ NOT NULL,                      
    revoked BOOLEAN NOT NULL DEFAULT FALSE,               

    user_agent VARCHAR(255),                              

    CONSTRAINT unique_token UNIQUE (refresh_token),                
    CONSTRAINT refresh_token_pkey PRIMARY KEY (id)
);

-- +migrate Down
DROP TABLE IF EXISTS app.session;