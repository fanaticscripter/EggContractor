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

		err = runMigrations(dbPath)
		if err != nil {
			err = errors.Wrapf(err, "error occurred during schema migrations")
			return
		}

		_db, err = sql.Open("sqlite3", dbPath+"?_foreign_keys=on&_journal_mode=WAL")
		if err != nil {
			err = errors.Wrapf(err, "failed to open SQLite3 database %#v", dbPath)
			return
		}
		err = nil
	})
	return err
}
