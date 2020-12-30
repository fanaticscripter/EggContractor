package solo

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/dustin/go-humanize"
	"github.com/fanaticscripter/EggContractor/util"
)

func (c *SoloContract) Display() {
	fmt.Printf("%s (%s)\n", c.GetName(), c.GetId())
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	contractType := "Standard"
	if c.GetIsElite() {
		contractType = "Elite"
	}
	fmt.Fprintf(w, "Type:\t%s\n", contractType)
	fmt.Fprintf(w, "Eggs laid:\t%s / %s\n",
		util.Numfmt(c.GetEggsLaid()), util.Numfmt(c.GetUltimateGoal()))
	fmt.Fprintf(w, "Hourly laying rate:\t%s current / %s required\n",
		util.Numfmt(c.GetEggsPerHour()), util.Numfmt(c.RequiredEggsPerHour()))
	fmt.Fprintf(w, "Time to complete:\t%s expected / %s remaining\n",
		util.FormatDuration(c.ExpectedDurationUntilFinish()),
		util.FormatDurationNonNegative(c.GetDurationUntilProductionDeadline()))
	fmt.Fprintf(w, "Reported to server:\t%s (%s)\n",
		util.FormatDatetime(c.GetLastRefreshedTime()),
		humanize.Time(c.GetLastRefreshedTime()))
	w.Flush()
	fmt.Println()
}
