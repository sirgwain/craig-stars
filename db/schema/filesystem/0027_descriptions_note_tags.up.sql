ALTER TABLE ship_designs
ADD description TEXT NOT NULL DEFAULT '';

ALTER TABLE races
ADD description TEXT NOT NULL DEFAULT '';

ALTER TABLE player_notes
ADD tags TEXT NOT NULL DEFAULT '[]';

-- fleet nums are reused, so a note on a fleet goes away when the fleet does
CREATE TRIGGER trg_fleets_delete_player_notes AFTER DELETE ON fleets WHEN OLD.starbase = 0 BEGIN
DELETE FROM player_notes
WHERE
  game_id = OLD.game_id
  AND target_type = 'Fleet'
  AND target_player_num = OLD.player_num
  AND target_num = OLD.num;

END;

-- a fleet given to another player, or renumbered, isn't the fleet the note was written for anymore
CREATE TRIGGER trg_fleets_update_player_notes AFTER
UPDATE OF player_num,
num ON fleets WHEN OLD.starbase = 0
AND (
  OLD.player_num <> NEW.player_num
  OR OLD.num <> NEW.num
) BEGIN
DELETE FROM player_notes
WHERE
  game_id = OLD.game_id
  AND target_type = 'Fleet'
  AND target_player_num = OLD.player_num
  AND target_num = OLD.num;

END;
