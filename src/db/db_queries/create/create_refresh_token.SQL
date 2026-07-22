CREATE TABLE refresh_tokens(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token_hash TEXT NOT NULL UNIQUE,
    user_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    FOREIGN KEY (user_id) REFERENCES users(id)
)