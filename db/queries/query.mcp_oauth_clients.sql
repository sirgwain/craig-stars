-- name: CreateMCPOAuthClient :one
INSERT INTO
    mcp_oauth_clients (
        client_id,
        client_name,
        client_uri,
        token_endpoint_auth_method,
        scope,
        client_id_issued_at
    )
VALUES
    (?, ?, ?, ?, ?, ?) RETURNING *;

-- name: CreateMCPOAuthRedirectURI :exec
INSERT INTO
    mcp_oauth_redirect_uris (client_id, redirect_uri)
VALUES
    (?, ?);

-- name: GetMCPOAuthClient :one
SELECT
    *
FROM
    mcp_oauth_clients
WHERE
    client_id = ?;

-- name: GetMCPOAuthRedirectURI :one
SELECT
    *
FROM
    mcp_oauth_redirect_uris
WHERE
    client_id = ?
    AND redirect_uri = ?;
