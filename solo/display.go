package solo

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/fanaticscripter/EggContractor/util"
)

func (c *SoloContract) Display(now time.Time, multiPlayerMode bool) {
	fmt.Printf("%s (%s)\n", c.GetName(), c.GetId())
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	if multiPlayerMode {
		playerName := c.GetPlayerNickname()
		if playerName == "" {
			playerName = "Unknown player"
		}
		fmt.Fprintf(w, "Player:\t%s\n", playerName)
	}
	contractType := "Standard"
	if c.GetIsElite() {
		contractType = "Elite"
	}
	fmt.Fprintf(w, "Type:\t%s\n", contractType)
	fmt.Fprintf(w, "Eggs laid:\t%s / %s\n",
		util.Numfmt(c.GetEggsLaid()), util.Numfmt(c.GetUltimateGoal()))
	fmt.Fprintf(w, "Eggs laid, offline-adjusted:\t%s\n",
		util.Numfmt(c.GetOfflineAdjustedEggsLaid(now)))
	fmt.Fprintf(w, "Hourly laying rate:\t%s current / %s required\n",
		util.Numfmt(c.GetEggsPerHour()), util.Numfmt(c.RequiredEggsPerHour()))
	fmt.Fprintf(w, "Time to complete:\t%s expected / %s remaining\n",
		util.FormatDuration(c.ExpectedDurationUntilFinish()),
		util.FormatDurationNonNegative(c.GetDurationUntilProductionDeadline()))
	fmt.Fprintf(w, "Time to complete, offline-adjusted:\t%s\n",
		util.FormatDuration(c.GetOfflineAdjustedExpectedDurationUntilFinish(now)))
	fmt.Fprintf(w, "Reported to server:\t%s (%s)\n",
		util.FormatDatetime(c.GetServerRefreshTime()),
		humanize.Time(c.GetServerRefreshTime()))
	w.Flush()
	fmt.Println()
}
