-- Create the version table for tracking updates to data in the database --
CREATE TABLE IF NOT EXISTS versions (
    id INTEGER PRIMARY KEY,
    createdAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
    updatedAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
    current NUMERIC
);
INSERT INTO versions (current)
VALUES (0);