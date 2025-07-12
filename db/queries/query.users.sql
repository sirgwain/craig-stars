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
    AND gameId = ?;

-- name: GetGetGuestUserForGame :one
SELECT
    *
FROM
    users
WHERE
    role = 'guest'
    AND gameId = ?
    AND playerNum = ?;

-- name: GetUsersForGame :many
SELECT
    *
FROM
    users
WHERE
    gameId = ?;

-- name: CreateUser :one
INSERT INTO
    users (
        createdAt,
        updatedAt,
        username,
        gameId,
        playerNum,
        password,
        email,
        role,
        banned,
        verified,
        lastLogin,
        discordId,
        discordAvatar,
        discordWebhookUrl
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
    updatedAt = CURRENT_TIMESTAMP,
    username = ?,
    gameId = ?,
    playerNum = ?,
    password = ?,
    email = ?,
    role = ?,
    banned = ?,
    verified = ?,
    lastLogin = ?,
    discordId = ?,
    discordAvatar = ?,
    discordWebhookUrl = ?
WHERE
    id = ?;

-- name: UpdateUserSettings :exec
UPDATE users
SET
    updatedAt = CURRENT_TIMESTAMP,
    discordWebhookUrl = ?
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
    AND gameId = ?
    AND playerNum = ?;

-- name: DeleteGameGuestUsers :exec
DELETE FROM users
WHERE
    role = 'guest'
    AND gameId = ?;
