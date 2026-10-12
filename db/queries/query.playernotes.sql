--
-- PlayerNotes
--
-- name: GetPlayerNote :one
SELECT
    *
FROM
    player_notes
WHERE
    id = ?;

-- name: GetPlayerNotesForPlayer :many
SELECT
    *
FROM
    player_notes
WHERE
    game_id = ?
    AND player_num = ?
ORDER BY
    id;

-- name: GetPlayerNoteForTarget :one
SELECT
    *
FROM
    player_notes
WHERE
    game_id = ?
    AND player_num = ?
    AND target_type = ?
    AND target_num = ?
    AND target_player_num = ?;

-- name: CreatePlayerNote :execlastid
INSERT INTO
    player_notes (
        created_at,
        updated_at,
        game_id,
        player_num,
        target_type,
        target_num,
        target_player_num,
        x,
        y,
        text,
        show_on_map,
        tags
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
        ?
    );

-- name: UpdatePlayerNote :execrows
UPDATE player_notes
SET
    updated_at = CURRENT_TIMESTAMP,
    x = ?,
    y = ?,
    text = ?,
    show_on_map = ?,
    tags = ?
WHERE
    id = ?;

-- name: DeletePlayerNote :execrows
DELETE FROM player_notes
WHERE
    id = ?;

-- name: DeletePlayerNoteForTarget :execrows
DELETE FROM player_notes
WHERE
    game_id = ?
    AND player_num = ?
    AND target_type = ?
    AND target_num = ?
    AND target_player_num = ?;
