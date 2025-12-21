-- computing spec at runtime now
ALTER TABLE minefields
DROP COLUMN spec;

ALTER TABLE players
DROP COLUMN spec;

ALTER TABLE fleets
ADD COLUMN report_age INTEGER NOT NULL DEFAULT 0;

ALTER TABLE ship_designs
ADD COLUMN report_age INTEGER NOT NULL DEFAULT 0;

ALTER TABLE planets
ADD COLUMN report_age INTEGER NOT NULL DEFAULT 0;

ALTER TABLE mineral_packets
ADD COLUMN report_age INTEGER NOT NULL DEFAULT 0;

ALTER TABLE salvages
ADD COLUMN report_age INTEGER NOT NULL DEFAULT 0;

ALTER TABLE wormholes
ADD COLUMN report_age INTEGER NOT NULL DEFAULT 0;

ALTER TABLE mystery_traders
ADD COLUMN report_age INTEGER NOT NULL DEFAULT 0;

ALTER TABLE minefields
ADD COLUMN report_age INTEGER NOT NULL DEFAULT 0;

UPDATE PLAYERS
SET fleet_intels = (
  SELECT json_group_array(
           json_set(
             -- remove scanRange and scanRangePen from root
             json_remove(
               json_remove(f.value, '$.scanRange'),
               '$.scanRangePen'
             ),
             -- add them into the spec object
             '$.spec.scanRange', json_extract(f.value, '$.scanRange'),
             '$.spec.scanRangePen', json_extract(f.value, '$.scanRangePen')
           )
         )
  FROM json_each(fleet_intels) AS f
);
