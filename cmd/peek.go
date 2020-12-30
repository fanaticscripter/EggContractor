package cmd

import (
	"strings"
	"time"

	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"

	"github.com/fanaticscripter/EggContractor/api"
	"github.com/fanaticscripter/EggContractor/db"
)

var _peekCommand = &cobra.Command{
	Use:     "peek <contract-id> <coop-code>",
	Short:   "Peek at a coop",
	Args:    cobra.ExactArgs(2),
	PreRunE: subcommandPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		contractId := strings.ToLower(args[0])
		code := strings.ToLower(args[1])

		now := time.Now()
		status, err := api.RequestCoopStatus(&api.CoopStatusRequestPayload{
			ContractId: contractId,
			Code:       code,
		})
		if err != nil {
			return err
		}
		wrapped, err := db.WrapCoopStatusWithDB(status)
		if err != nil {
			log.Error(err)
		} else if wrapped.Contract == nil {
			log.Warnf("contract %s not found in database, try using the refresh subcommand to populate the contract table",
				contractId)
		}
		wrapped.Display(_sortBy.by)

		err = db.InsertPeeked(db.NewPeeked(wrapped, now))
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	_rootCmd.AddCommand(_peekCommand)
}
