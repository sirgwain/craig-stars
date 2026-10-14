--
-- PlayerMail
--
-- name: GetPlayerMail :one
SELECT
    *
FROM
    player_mail
WHERE
    id = ?;

-- the mail a player wrote, and the mail delivered to them before a year
-- name: GetPlayerMailForPlayer :many
SELECT
    *
FROM
    player_mail
WHERE
    game_id = sqlc.arg (game_id)
    AND (
        from_player_num = sqlc.arg (player_num)
        OR (
            year < sqlc.arg (year)
            AND to_player_num IN (0, sqlc.arg (player_num))
        )
    )
ORDER BY
    id;

-- name: GetPlayerMailForGameYear :many
SELECT
    *
FROM
    player_mail
WHERE
    game_id = ?
    AND year = ?
ORDER BY
    id;

-- name: CreatePlayerMail :execlastid
INSERT INTO
    player_mail (
        created_at,
        updated_at,
        game_id,
        year,
        from_player_num,
        to_player_num,
        text
    )
VALUES
    (
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP,
        ?,
        ?,
        ?,
        ?,
        ?
    );

-- name: UpdatePlayerMail :execrows
UPDATE player_mail
SET
    updated_at = CURRENT_TIMESTAMP,
    to_player_num = ?,
    text = ?
WHERE
    id = ?;

-- name: DeletePlayerMail :execrows
DELETE FROM player_mail
WHERE
    id = ?;
