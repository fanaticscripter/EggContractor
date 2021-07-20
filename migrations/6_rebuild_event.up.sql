-- Event id is not unique, so we need to change the primary key.

PRAGMA foreign_keys = off;

CREATE TABLE IF NOT EXISTS event_new (
    rowid INTEGER PRIMARY KEY,
    id TEXT NOT NULL,
    event_type TEXT NOT NULL,
    multiplier REAL NOT NULL,
    message TEXT NOT NULL,
    first_seen_timestamp REAL NOT NULL,
    last_seen_timestamp REAL NOT NULL,
    expiry_timestamp REAL NOT NULL
);

INSERT INTO event_new(
    id,
    event_type,
    multiplier,
    message,
    first_seen_timestamp,
    last_seen_timestamp,
    expiry_timestamp
)
SELECT
    id,
    event_type,
    multiplier,
    message,
    first_seen_timestamp,
    last_seen_timestamp,
    expiry_timestamp
FROM "event"
ORDER BY first_seen_timestamp;

DROP TABLE "event";

ALTER TABLE event_new RENAME TO "event";

PRAGMA foreign_keys = on;
