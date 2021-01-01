package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/fanaticscripter/EggContractor/api"
	"github.com/fanaticscripter/EggContractor/db"
	"github.com/fanaticscripter/EggContractor/util"
)

var (
	_eventsNoRefresh bool
	_eventsCommand   = &cobra.Command{
		Use:     "events",
		Short:   "Print current and past events",
		Args:    cobra.NoArgs,
		PreRunE: subcommandPreRunE,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !_eventsNoRefresh {
				activeEvents, err := refreshEvents()
				if err != nil {
					log.Error(err)
				}
				if len(activeEvents) != 0 {
					fmt.Println("Active events:")
					w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
					for _, e := range activeEvents {
						timeLeft := util.DoubleToDuration((e.SecondsRemaining))
						fmt.Fprintf(w, "%s remaining\t%s\n", util.FormatDuration(timeLeft), e.Message)
					}
					w.Flush()
				} else {
					fmt.Println("No active events.")
				}
				fmt.Println()
			}

			events, err := db.GetEvents()
			if err != nil {
				return err
			}

			if _eventsNoRefresh {
				activeEvents := make([]*db.Event, 0)
				for _, e := range events {
					if e.HasTimeLeft() {
						activeEvents = append(activeEvents, e)
					}
				}
				if len(activeEvents) != 0 {
					fmt.Println("Active events:")
					w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
					for _, e := range activeEvents {
						fmt.Fprintf(w, "%s remaining\t%s\n", util.FormatDuration(e.TimeLeft()), e.Message)
					}
					w.Flush()
				} else {
					fmt.Println("No active events.")
				}
				fmt.Println()
			}

			fmt.Println("Recorded event history:")
			table := [][]string{
				{"Date", "Started", "Type", "Details", "Duration"},
				{"----", "-------", "----", "-------", "--------"},
			}
			for _, e := range events {
				table = append(table, []string{
					util.FormatDateCasual(e.FirstSeenTime),
					util.FormatTimeCasual(e.FirstSeenTime),
					e.EventType,
					e.UnhypedMessage(),
					util.FormatDuration(e.Duration()),
				})
			}
			util.PrintTable(table)

			return nil
		},
	}
)

func init() {
	_rootCmd.AddCommand(_eventsCommand)
	_eventsCommand.Flags().BoolVarP(&_eventsNoRefresh, "no-refresh", "n", false, "do not refresh current events")
}

func refreshEvents() (active []*api.Event, err error) {
	now := time.Now()
	p, err := api.RequestPeriodicals(&api.GetPeriodicalsRequestPayload{
		PlayerId:     _config.Player.Id,
		X2:           1,
		EarningBonus: 1e12, // Use a reasonably large EB just in case
	})
	if err != nil {
		return
	}
	active = p.Events.Events
	seen := now
	if p.Contracts.ResponseTimestamp != 0 {
		seen = util.DoubleToTime(p.Contracts.ResponseTimestamp)
	}
	for _, e := range active {
		if err := db.InsertEvent(seen, e); err != nil {
			log.Error(err)
		}
	}
	return
}
