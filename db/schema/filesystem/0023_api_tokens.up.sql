CREATE TABLE api_tokens (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL DEFAULT '',
    token_prefix TEXT NOT NULL DEFAULT '',
    token_hash TEXT UNIQUE NOT NULL,
    scope TEXT NOT NULL DEFAULT '',
    expires_at TIMESTAMP NOT NULL,
    last_used_at TIMESTAMP,
    revoked_at TIMESTAMP,
    CONSTRAINT fk_users_api_tokens FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX idx_api_tokens_token_hash ON api_tokens (token_hash);
CREATE INDEX idx_api_tokens_user_id ON api_tokens (user_id);
