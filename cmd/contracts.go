package cmd

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/fanaticscripter/EggContractor/db"
	"github.com/fanaticscripter/EggContractor/util"
)

var _contractsCommand = &cobra.Command{
	Use:     "contracts",
	Short:   "Print a list of current and past contracts",
	Args:    cobra.NoArgs,
	PreRunE: subcommandPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		contracts, err := db.GetContracts()
		if err != nil {
			return err
		}

		if len(contracts) == 0 {
			log.Warn("no contracts found in database, try using the refresh subcommand to populate the contract table")
			return nil
		}

		table := [][]string{
			{"ID", "Name", "Egg", "E. Goal", "Time", "Size", "Token", "Expiry"},
			{"--", "----", "---", "-------", "----", "----", "-----", "------"},
		}

		for _, c := range contracts {
			if c.Id == "first-contract" {
				continue
			}
			table = append(table, []string{
				c.Id, c.Name, c.EggType.Display(),
				util.NumfmtWhole(c.UltimateGoal(true)),
				util.FormatDurationWhole(c.Duration()),
				fmt.Sprintf("%d", c.MaxCoopSize),
				fmt.Sprintf("%.0fm", c.TokenInterval),
				util.FormatDate(c.ExpiryTime()),
			})
		}
		util.PrintTable(table)

		return nil
	},
}

func init() {
	_rootCmd.AddCommand(_contractsCommand)
}
