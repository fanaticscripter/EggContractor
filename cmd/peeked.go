package cmd

import (
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/fanaticscripter/EggContractor/db"
	"github.com/fanaticscripter/EggContractor/util"
)

const _peekedThreshold = 7 * 24 * time.Hour

var _peekedCommand = &cobra.Command{
	Use:   "peeked [<contract-id>]",
	Short: "Print list of recently peeked coops",
	Long: `Print list of recently peeked coops

Recently currently means "within one week". If a contract-id is specified,
only coops for that contract is listed.

Note that the /peeked/ page of the web client is much richer in information.`,
	Args:    cobra.MaximumNArgs(1),
	PreRunE: subcommandPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		selectedContractId := ""
		if len(args) > 0 {
			selectedContractId = strings.ToLower(args[0])
		}

		contractIds, groups, err := db.GetPeekedGroupedByContract(time.Now().Add(-_peekedThreshold))
		if err != nil {
			return err
		}
		if selectedContractId != "" {
			_, exists := groups[selectedContractId]
			if !exists {
				log.Warnf("You haven't peeked any coop for contract %s recently.", selectedContractId)
				return nil
			}
			contractIds = []string{selectedContractId}
		}

		table := [][]string{
			{"Contract ID", "Code", "Spots", "Laid", "Rate/Required", "Time left"},
			{"-----------", "----", "-----", "----", "-------------", "---------"},
		}
		for groupIdx, contractId := range contractIds {
			for i, p := range groups[contractId] {
				contractIdField := p.ContractId
				if i != 0 {
					// Since rows are grouped by contract, we can merge the Contract ID cells.
					contractIdField = ""
				}
				spotsField := fmt.Sprintf("%d", p.Openings)
				laidField := util.Numfmt(p.EggsLaid)
				rateVsRequiredField := util.Numfmt(p.EggsPerHour)
				if p.RequiredEggsPerHour != 0 {
					rateVsRequiredField = fmt.Sprintf("%s / %s",
						util.Numfmt(p.EggsPerHour), util.Numfmt(p.RequiredEggsPerHour))
				}
				timeLeftField := util.FormatDurationNonNegative(p.TimeLeft)
				table = append(table, []string{
					contractIdField,
					p.Code,
					spotsField,
					laidField,
					rateVsRequiredField,
					timeLeftField,
				})
			}
			if groupIdx != len(contractIds)-1 {
				table = append(table, []string{
					"-----------", "----", "-----", "----", "-------------", "---------",
				})
			}
		}
		util.PrintTable(table)
		fmt.Println()
		log.Warn("the /peeked/ page of the web client is much richer in information")

		return nil
	},
}

func init() {
	_rootCmd.AddCommand(_peekedCommand)
}
