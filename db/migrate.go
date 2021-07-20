package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	sqlite3driver "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"

	"github.com/fanaticscripter/EggContractor/api"
)

const _schemaVersion = 6

func runMigrations(dbPath string) (err error) {
	sql.Register("sqlite3_with_migration_funcs", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			if err := conn.RegisterFunc("contract_expiry_timestamp", contractExpiryTimestamp, true); err != nil {
				return err
			}
			if err := conn.RegisterFunc("contract_expiry_year", contractExpiryYear, true); err != nil {
				return err
			}
			if err := conn.RegisterFunc("contract_expiry_month", contractExpiryMonth, true); err != nil {
				return err
			}
			return nil
		},
	})
	db, err := sql.Open("sqlite3_with_migration_funcs", dbPath+"?_foreign_keys=on&_journal_mode=WAL")
	if err != nil {
		return errors.Wrapf(err, "failed to open SQLite3 database %#v for migrations", dbPath)
	}
	defer func() {
		closeErr := db.Close()
		if err == nil && closeErr != nil {
			err = errors.Wrap(closeErr, "error closing database after migrations")
		}
	}()

	var driver database.Driver
	driver, err = sqlite3driver.WithInstance(db, &sqlite3driver.Config{})
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
	return
}

// Used in migration #3 (3_add_contract_expiry_timestamp.up.sql)
func contractExpiryTimestamp(textId string, props []byte) float64 {
	contract := &api.ContractProperties{}
	if err := proto.Unmarshal(props, contract); err != nil {
		msg := fmt.Sprintf("cannot unmarshal props for contract %s", textId)
		log.Error(msg)
		panic(msg)
	}
	return contract.ExpiryTimestamp
}

// Used in migration #5 (5_add_contract_expiry_month.up.sql)
func contractExpiryYear(textId string, props []byte) int {
	contract := &api.ContractProperties{}
	if err := proto.Unmarshal(props, contract); err != nil {
		msg := fmt.Sprintf("cannot unmarshal props for contract %s", textId)
		log.Error(msg)
		panic(msg)
	}
	expiry := contract.ExpiryTime().In(time.UTC)
	return int(expiry.Year())
}

// Used in migration #5 (5_add_contract_expiry_month.up.sql)
func contractExpiryMonth(textId string, props []byte) int {
	contract := &api.ContractProperties{}
	if err := proto.Unmarshal(props, contract); err != nil {
		msg := fmt.Sprintf("cannot unmarshal props for contract %s", textId)
		log.Error(msg)
		panic(msg)
	}
	expiry := contract.ExpiryTime().In(time.UTC)
	return int(expiry.Month())
}
