package db

import (
	"compress/gzip"
	"context"
	"database/sql"
	"io"
	"os"
	"time"

	"github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

func BackupDatabase(dbPath, backupPath string) error {
	return BackupDatabaseWithContext(context.Background(), dbPath, backupPath)
}

func BackupDatabaseWithContext(ctx context.Context, dbPath, backupPath string) (err error) {
	conns := make([]*sqlite3.SQLiteConn, 0, 2)
	sql.Register("sqlite3_record_conns", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			conns = append(conns, conn)
			return nil
		},
	})

	srcDb, err := sql.Open("sqlite3_record_conns", dbPath+"?_foreign_keys=on&_journal_mode=WAL")
	if err != nil {
		return errors.Wrapf(err, "error opening source database %#v", dbPath)
	}
	defer func() { srcDb.Close() }()
	// Invoke connect hook; source SQLiteConn stored in conns[0].
	err = srcDb.Ping()
	if err != nil {
		return errors.Wrapf(err, "error connecting to source database %#v", dbPath)
	}
	dstDb, err := sql.Open("sqlite3_record_conns", backupPath)
	if err != nil {
		return errors.Wrapf(err, "error opening backup database %#v", dbPath)
	}
	defer func() { dstDb.Close() }()
	// Invoke connect hook; dest SQLiteConn stored in conns[1].
	err = dstDb.Ping()
	if err != nil {
		return errors.Wrapf(err, "error connecting to backup database %#v", backupPath)
	}

	backup, err := conns[1].Backup("main", conns[0], "main")
	if err != nil {
		return errors.Wrapf(err, "error initiating backup")
	}
	defer func() {
		finishErr := backup.Finish()
		if err == nil && finishErr != nil {
			err = finishErr
		}
	}()

	// https://sqlite.org/backup.html
	for {
		select {
		case <-ctx.Done():
			err = errors.Wrapf(ctx.Err(), "backup cancelled with %d/%d pages remaining",
				backup.Remaining(), backup.PageCount())
			return
		default:
		}
		var done bool
		done, err = backup.Step(16)
		if err != nil {
			err = errors.Wrap(err, "during backup step")
			return
		}
		if done {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	return nil
}

func GzipFileRemoveOriginal(path string) error {
	err := GzipFile(path)
	if err != nil {
		return err
	}
	err = os.Remove(path)
	if err != nil {
		return errors.Wrapf(err, "failed to remove %#v", path)
	}
	return nil
}

func GzipFile(path string) error {
	dest := path + ".gz"
	fin, err := os.Open(path)
	if err != nil {
		return errors.Wrapf(err, "error opening %#v", path)
	}
	defer fin.Close()
	fout, err := os.Create(dest)
	if err != nil {
		return errors.Wrapf(err, "error creating %#v", dest)
	}
	defer fout.Close()

	zw := gzip.NewWriter(fout)
	buf := make([]byte, 65536)
LoopReadChunks:
	for {
		n, err := fin.Read(buf)
		switch {
		case err == io.EOF:
			break LoopReadChunks
		case err != nil:
			return errors.Wrapf(err, "error reading %#v", path)
		}
		_, err = zw.Write(buf[:n])
		if err != nil {
			return errors.Wrapf(err, "error writing to %#v", dest)
		}
	}
	err = zw.Close()
	if err != nil {
		return errors.Wrapf(err, "error writing gzip footer to %#v", dest)
	}
	return nil
}
