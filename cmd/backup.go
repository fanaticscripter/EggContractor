package cmd

import (
	"fmt"
	"time"

	"github.com/fanaticscripter/EggContractor/db"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var _backupCommand = &cobra.Command{
	Use:   "backup",
	Short: "Backup database",
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := initConfig(); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		dbPath := _config.Database.Path
		backupPath := dbPath + fmt.Sprintf(".%d.bak", start.Unix())
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Minute)
		err := db.BackupDatabaseWithContext(ctx, dbPath, backupPath)
		elapsed := time.Since(start).Seconds()
		if err != nil {
			return errors.Wrapf(err, "failed to back up %#v to %#v (%.1fs elapsed)", dbPath, backupPath, elapsed)
		}
		fmt.Printf("Backed up %#v to %#v in %.1fs.\n", dbPath, backupPath, elapsed)
		err = db.GzipFileRemoveOriginal(backupPath)
		if err != nil {
			return errors.Wrapf(err, "failed to gzip backup")
		}
		fmt.Println("Gzipped backup.")
		return nil
	},
}

func init() {
	_rootCmd.AddCommand(_backupCommand)
}
