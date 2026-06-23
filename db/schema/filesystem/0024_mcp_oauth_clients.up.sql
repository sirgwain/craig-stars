CREATE TABLE mcp_oauth_clients (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    client_id TEXT UNIQUE NOT NULL,
    client_name TEXT NOT NULL DEFAULT '',
    client_uri TEXT NOT NULL DEFAULT '',
    token_endpoint_auth_method TEXT NOT NULL DEFAULT 'none',
    scope TEXT NOT NULL DEFAULT '',
    client_id_issued_at INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE mcp_oauth_redirect_uris (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    client_id TEXT NOT NULL,
    redirect_uri TEXT NOT NULL,
    CONSTRAINT fk_mcp_oauth_redirect_client FOREIGN KEY (client_id) REFERENCES mcp_oauth_clients (client_id) ON DELETE CASCADE,
    UNIQUE (client_id, redirect_uri)
);

CREATE INDEX idx_mcp_oauth_redirect_uris_client_id ON mcp_oauth_redirect_uris (client_id);
