package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/fanaticscripter/EggContractor/api"
	"github.com/fanaticscripter/EggContractor/util"
)

var _afxConfigCommand = &cobra.Command{
	Use:     "afx-config",
	Short:   "Explore /afx_config",
	Args:    cobra.NoArgs,
	PreRunE: subcommandPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		req := &api.ArtifactsConfigurationRequestPayload{
			ClientVersion: api.ClientVersion,
		}
		resp := &api.ArtifactsConfigurationResponsePayload{}
		err := api.Request("/ei_afx/config", req, resp)
		if err != nil {
			return err
		}
		config := resp.Config

		table := [][]string{
			{"Ship", "Type", "Duration", "Capacity", "Quality", "Quality Range"},
		}
		for _, m := range config.MissionParameters {
			table = append(table, []string{
				"----", "----", "--------", "--------", "-------", "-------------",
			})
			ship := m.Ship
			for i, d := range m.Durations {
				var shipNameField string
				if i == 0 {
					shipNameField = ship.Name()
				}
				table = append(table, []string{
					shipNameField,
					d.DurationType.Display(),
					util.FormatDurationWhole(util.DoubleToDuration(d.Seconds)),
					fmt.Sprintf("%d", d.Capacity),
					fmt.Sprintf("%.1f", d.Quality),
					fmt.Sprintf("%.1f - %.1f", d.MinQuality, d.MaxQuality),
				})
			}
		}
		util.PrintTable(table)

		fmt.Println()
		for _, a := range config.ArtifactParameters {
			tierNumber := a.Spec.TierNumber()
			if tierNumber == 1 || a.Spec.Rarity > 0 {
				continue
			}
			fmt.Printf("{api.ArtifactSpec_%s, %d}: {%f, %f, %d, %v},\n",
				a.Spec.Name.String(), tierNumber,
				a.CraftingPrice, a.CraftingPriceLow, a.CraftingPriceDomain, a.CraftingPriceCurve,
			)
		}

		return nil
	},
}

func init() {
	_rootCmd.AddCommand(_afxConfigCommand)
}
