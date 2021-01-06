package db

import (
	"database/sql"
	"os"
	"path/filepath"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/fanaticscripter/EggContractor/config"
)

const _schemaVersion = 2

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

		_db, err = sql.Open("sqlite3", dbPath+"?_foreign_keys=on&_journal_mode=WAL")
		if err != nil {
			err = errors.Wrapf(err, "failed to open SQLite3 database %#v", dbPath)
			return
		}
		var driver database.Driver
		driver, err = sqlite3.WithInstance(_db, &sqlite3.Config{})
		if err != nil {
			err = errors.Wrap(err, "failed to initialize migrations driver")
			return
		}
		var m *migrate.Migrate
		m, err = migrate.NewWithDatabaseInstance("file://migrations", "sqlite3", driver)
		if err != nil {
			err = errors.Wrap(err, "failed to initialize schema migrator")
			return
		}
		err = m.Migrate(_schemaVersion)
		if err != nil && err != migrate.ErrNoChange {
			err = errors.Wrapf(err, "failed to migrate schemas in database %#v", dbPath)
			return
		}
		err = nil
	})
	return err
}
