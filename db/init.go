package db

import (
	"database/sql"
	"os"
	"path/filepath"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/fanaticscripter/EggContractor/config"
)

var (
	_db         *sql.DB
	_initDBOnce sync.Once
)

const _schema = `CREATE TABLE IF NOT EXISTS contract (
	id INTEGER PRIMARY KEY,
	-- 'text_id' is the contract ID used by the API, e.g. 'deep-foods-2020'.
	text_id TEXT NOT NULL,
	-- 'expiry_year' is the year of the expiry timestamp used for future-proofing,
	-- since I'm not exactly sure if the text_id's are recycled in Leggacy mode.
	-- The expiry timestamp is used because that's the only static timestamp
	-- the API returns for each contract.
	expiry_year INTEGER NOT NULL,
	coop_allowed BOOLEAN NOT NULL,
	-- 'props' is the protobuf-serialized version of contract properties as
	-- received from API.
	props BLOB NOT NULL,
	UNIQUE(text_id, expiry_year)
);
CREATE UNIQUE INDEX IF NOT EXISTS contract_text_id_expiry_year ON contract(text_id, expiry_year);

CREATE TABLE IF NOT EXISTS coop (
	id INTEGER PRIMARY KEY,
	contract_id INTEGER NOT NULL REFERENCES contract(id)
		ON UPDATE CASCADE ON DELETE CASCADE,
	code TEXT NOT NULL,
	UNIQUE(contract_id, code)
);
CREATE INDEX IF NOT EXISTS coop_contract_id ON coop(contract_id);
CREATE UNIQUE INDEX IF NOT EXISTS coop_contract_id_code ON coop(contract_id, code);

CREATE TABLE IF NOT EXISTS coop_status (
	id INTEGER PRIMARY KEY,
	coop_id INTEGER NOT NULL REFERENCES coop(id)
		ON UPDATE CASCADE ON DELETE CASCADE,
	refresh_id INTEGER REFERENCES refresh(id)
	    ON UPDATE CASCADE ON DELETE SET NULL,
	-- 'timestamp' is the timestamp at which the status was retrieved.
	timestamp REAL NOT NULL,
	-- 'status' is the protobuf-serialized version of coop status as received
	-- from API.
	status BLOB NOT NULL
);
CREATE INDEX IF NOT EXISTS coop_status_coop_id ON coop_status(coop_id);
CREATE INDEX IF NOT EXISTS coop_status_refresh_id ON coop_status(refresh_id);

CREATE TABLE IF NOT EXISTS solo_status(
	id INTEGER PRIMARY KEY,
	contract_id INTEGER NOT NULL REFERENCES contract(id)
		ON UPDATE CASCADE ON DELETE CASCADE,
	refresh_id INTEGER REFERENCES refresh(id)
	    ON UPDATE CASCADE ON DELETE SET NULL,
	-- 'timestamp' is the timestamp at which the status was retrieved.
	timestamp REAL NOT NULL,
	-- 'status' is the protobuf-serialized version of the solo contract status.
	status BLOB NOT NULL
);
CREATE INDEX IF NOT EXISTS solo_status_contract_id ON solo_status(contract_id);
CREATE INDEX IF NOT EXISTS solo_status_refresh_id ON solo_status(refresh_id);

CREATE TABLE IF NOT EXISTS refresh (
	id INTEGER PRIMARY KEY,
	timestamp REAL NOT NULL
);

-- All fields are marked not null in this table to ease coding.
-- Use zero values when field is unknown (due to not being able to match coop to contract).
CREATE TABLE IF NOT EXISTS peeked (
	id INTEGER PRIMARY KEY,
	contract_id TEXT NOT NULL,
	code TEXT NOT NULL,
	last_peeked REAL NOT NULL,
	has_completed BOOLEAN NOT NULL,
	openings INTEGER NOT NULL,
	eggs_laid REAL NOT NULL,
	eggs_per_hour REAL NOT NULL,
	required_eggs_per_hour REAL NOT NULL,
	-- 'time_left' is measured in seconds.
	time_left REAL NOT NULL,
	-- 'max_eb_percentage' is the EB percentage of the highest EB member.
	max_eb_percentage REAL NOT NULL,
	-- 'mean_eb_percentage' is the geometric mean EB percentage of all current members.
	mean_eb_percentage REAL NOT NULL,
	UNIQUE(contract_id, code)
);
CREATE INDEX IF NOT EXISTS peeked_contract_id_code ON peeked(contract_id, code);

CREATE TABLE IF NOT EXISTS event (
	id TEXT PRIMARY KEY,
	event_type TEXT NOT NULL,
	multiplier REAL NOT NULL,
	message TEXT NOT NULL,
	first_seen_timestamp REAL NOT NULL,
	last_seen_timestamp REAL NOT NULL,
	expiry_timestamp REAL NOT NULL
);`

func InitDB(conf *config.Config) error {
	var err error
	_initDBOnce.Do(func() {
		dbPath := conf.Database.Path
		log.Debugf("database path: %s", dbPath)

		parentDir := filepath.Dir(dbPath)
		err = os.MkdirAll(parentDir, 0o755)
		if err != nil {
			err = errors.Wrapf(err, "failed to create parent directory %#v for database", parentDir)
			return
		}

		_db, err = sql.Open("sqlite3", dbPath+"?_foreign_keys=on&_journal_mode=WAL")
		if err != nil {
			err = errors.Wrapf(err, "failed to open SQLite3 database %#v", dbPath)
			return
		}
		_, err = _db.Exec(_schema)
		if err != nil {
			err = errors.Wrapf(err, "failed to initialize schema in database %#v", dbPath)
			return
		}
		err = nil
	})
	return err
}
