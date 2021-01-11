package cmd

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/fanaticscripter/EggContractor/config"
	"github.com/fanaticscripter/EggContractor/coop"
	"github.com/fanaticscripter/EggContractor/db"
)

type sortBy struct {
	name string
	by   coop.By
}

func (by sortBy) String() string {
	return by.name
}

func (by *sortBy) Set(s string) error {
	switch strings.ToLower(s) {
	case "eggs_laied":
		fallthrough
	case "contribution":
		fallthrough
	case "total":
		fallthrough
	case "laied":
		by.by = coop.ByEggsLaid

	case "laying_rate":
		fallthrough
	case "rate":
		by.by = coop.ByLayingRate

	case "earning_bonus":
		fallthrough
	case "eb":
		by.by = coop.ByEarningBonus

	default:
		return fmt.Errorf("unrecognized sorting criterion %#v", s)
	}
	by.name = s
	return nil
}

func (by sortBy) Type() string {
	return "criterion"
}

var (
	_sortBy  sortBy
	_verbose bool
	_debug   bool
	_cfgFile string
	_rootCmd = &cobra.Command{
		Use: config.ProgName,
	}
)

func init() {
	cobra.OnInitialize(cobraHouseKeeping, configureLoggingLevel)

	_sortBy = sortBy{
		name: "eggs_laid",
		by:   coop.ByEggsLaid,
	}
	_rootCmd.PersistentFlags().VarP(&_sortBy, "sort", "s",
		"sort coop members by one of the following criteria: 'eggs_laid' (aliases: 'contribution', 'total', 'laid'), 'laying_rate' (alias: 'rate'), or 'earning_bonus' (alias: 'eb')")
	_rootCmd.PersistentFlags().BoolVarP(&_verbose, "verbose", "v", false, "enable verbose logging")
	_rootCmd.PersistentFlags().BoolVar(&_debug, "debug", false, "enable debug logging")
	_rootCmd.PersistentFlags().StringVar(&_cfgFile, "config", "", "config file, could also be set through env var EGGCONTRACTOR_CONFIG_FILE (default ~/.config/EggContractor/config.toml)")
}

func Execute() {
	if err := _rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func cobraHouseKeeping() {
	// Do not repeat error and print usage when a subcommand's RunE returns an
	// error. See https://github.com/spf13/cobra/issues/340.
	_rootCmd.SilenceUsage = true
	_rootCmd.SilenceErrors = true
}

func configureLoggingLevel() {
	log.SetLevel(log.WarnLevel)
	if _verbose {
		log.SetLevel(log.InfoLevel)
	}
	if _debug {
		log.SetLevel(log.DebugLevel)
	}
}

// The common prerun function used by most but not all subcommands.
func subcommandPreRunE(cmd *cobra.Command, args []string) error {
	if err := initConfig(); err != nil {
		return err
	}
	if err := db.InitDB(_config.Database); err != nil {
		return err
	}
	return nil
}
