-- name: GetUsers :many
SELECT
    *
FROM
    users;

-- name: GetUser :one
SELECT
    *
FROM
    users
WHERE
    id = ?;

-- name: GetUserByUsername :one
SELECT
    *
FROM
    users
WHERE
    username = ?;

-- name: GetGuestUser :one
SELECT
    *
FROM
    users
WHERE
    role = 'guest'
    AND password = ?;

-- name: GetGuestUsersForGame :many
SELECT
    *
FROM
    users
WHERE
    role = 'guest'
    AND game_id = ?;

-- name: GetGetGuestUserForGame :one
SELECT
    *
FROM
    users
WHERE
    role = 'guest'
    AND game_id = ?
    AND player_num = ?;

-- name: GetUsersForGame :many
SELECT
    *
FROM
    users
WHERE
    id IN (
        SELECT
            p.user_id
        FROM
            players p
        WHERE
            p.game_id = ?
            AND p.user_id != 0 -- skip AI slots
    );

-- name: CreateUser :one
INSERT INTO
    users (
        created_at,
        updated_at,
        username,
        game_id,
        player_num,
        password,
        email,
        role,
        banned,
        verified,
        last_login,
        discord_id,
        discord_avatar,
        discord_webhook_url
    )
VALUES
    (
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?
    ) RETURNING *;

-- name: UpdateUser :exec
UPDATE users
SET
    updated_at = CURRENT_TIMESTAMP,
    username = ?,
    game_id = ?,
    player_num = ?,
    password = ?,
    email = ?,
    role = ?,
    banned = ?,
    verified = ?,
    last_login = ?,
    discord_id = ?,
    discord_avatar = ?,
    discord_webhook_url = ?
WHERE
    id = ?;

-- name: UpdateUserSettings :exec
UPDATE users
SET
    updated_at = CURRENT_TIMESTAMP,
    discord_webhook_url = ?
WHERE
    id = ?;

-- name: DeleteUser :exec
DELETE FROM users
WHERE
    id = ?;

-- name: DeleteGameGuestUser :exec
DELETE FROM users
WHERE
    role = 'guest'
    AND game_id = ?
    AND player_num = ?;

-- name: DeleteGameGuestUsers :exec
DELETE FROM users
WHERE
    role = 'guest'
    AND game_id = ?;