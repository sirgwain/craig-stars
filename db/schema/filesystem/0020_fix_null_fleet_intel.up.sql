-- fix issue where fleet intels were returning null strings
UPDATE players
SET
  fleet_intels = 'null'
WHERE
  fleet_intels = '[null]';