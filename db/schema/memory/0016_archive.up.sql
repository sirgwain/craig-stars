ALTER TABLE players
ADD archived NUMERIC NOT NULL default 0;

ALTER TABLE games
ADD archived NUMERIC NOT NULL default 0;
