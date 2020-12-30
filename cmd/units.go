package cmd

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/fanaticscripter/EggContractor/util"
)

var _unitsCommand = &cobra.Command{
	Use:   "units",
	Short: "Print a table of units (order of magnitudes)",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		table := [][]string{
			{"Symb", "OoM", "Symb", "OoM", "Symb", "OoM"},
			{"----", "---", "----", "---", "----", "---"},
		}
		row := make([]string, 0, 6)
		for _, u := range util.Units {
			row = append(row, u.Symbol, strconv.Itoa(u.OoM))
			if len(row) == 6 {
				table = append(table, row)
				row = make([]string, 0, 6)
			}
		}
		util.PrintTable(table)
	},
}

func init() {
	_rootCmd.AddCommand(_unitsCommand)
}
