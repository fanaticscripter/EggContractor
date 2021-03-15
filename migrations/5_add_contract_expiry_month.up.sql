-- At the current rate of leggacy contracts it is a good chance we will
-- eventually get two runs of the same contract in the same calendar year.
-- Adding expiry_month to future-proof such conflicts.
--
-- We have to rebuild the contract table since the UNIQUE CONTRAINT has to be
-- modified.

PRAGMA foreign_keys = off;

CREATE TABLE IF NOT EXISTS contract_new (
    id INTEGER PRIMARY KEY,
    -- 'text_id' is the contract ID used by the API, e.g. 'deep-foods-2020'.
    text_id TEXT NOT NULL,
    -- 'expiry_year' is the year of the expiry timestamp.
    expiry_year INTEGER NOT NULL,
    -- 'expiry_month' is the month of the expiry timestamp.
    expiry_month INTEGER NOT NULL,
    coop_allowed BOOLEAN NOT NULL,
    expiry_timestamp REAL NOT NULL,
    first_seen_timestamp REAL,
    -- 'props' is the protobuf-serialized version of contract properties as
    -- received from API.
    props BLOB NOT NULL,
    UNIQUE(text_id, expiry_year, expiry_month)
);
CREATE UNIQUE INDEX IF NOT EXISTS contract_text_id_expiry_year_month ON contract_new(text_id, expiry_year, expiry_month);

INSERT INTO contract_new(
    id,
    text_id,
    expiry_year,
    expiry_month,
    coop_allowed,
    expiry_timestamp,
    first_seen_timestamp,
    props
)
SELECT
    id,
    text_id,
    contract_expiry_year(text_id, props),
    contract_expiry_month(text_id, props),
    coop_allowed,
    expiry_timestamp,
    first_seen_timestamp,
    props
FROM "contract"
ORDER BY id;

DROP TABLE "contract";

ALTER TABLE contract_new RENAME TO "contract";

PRAGMA foreign_keys = on;
