ALTER TABLE mysteryTraders ADD COLUMN requestedBoon INTEGER;
ALTER TABLE mysteryTraders ADD COLUMN destinationX REAL;
ALTER TABLE mysteryTraders ADD COLUMN destinationY REAL;
ALTER TABLE mysteryTraders ADD COLUMN rewardType TEXT;
ALTER TABLE players ADD COLUMN acquiredTechs TEXT;
ALTER TABLE shipDesigns ADD COLUMN mysteryTrader NUMERIC;
