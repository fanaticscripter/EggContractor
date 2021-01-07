package coop

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/fanaticscripter/EggContractor/api"
	"github.com/fanaticscripter/EggContractor/util"
)

// Display prints a formatted report about the coop's current status.
// contract is optional, and without it, the goal, expected time to complete
// etc. cannot be calculated.
func (c *CoopStatus) Display(sortBy By) {
	contract := c.Contract
	if contract != nil {
		fmt.Printf("%s (%s)\n", contract.Name, c.ContractId)
	} else {
		fmt.Println(c.ContractId)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	contractType := "Standard"
	if c.IsElite() {
		contractType = "Elite"
	}
	fmt.Fprintf(w, "Type:\t%s\n", contractType)
	fmt.Fprintf(w, "Code:\t%s\n", c.Code)
	if c.Creator() != nil {
		fmt.Fprintf(w, "Created by:\t%s\n", c.Creator().Name)
	}
	playersField := strconv.Itoa(len(c.Members))
	if contract != nil {
		playersField = fmt.Sprintf("%d / %d", len(c.Members), contract.MaxCoopSize)
	}
	fmt.Fprintf(w, "Players:\t%s\n", playersField)
	eggsLaidField := util.Numfmt(c.EggsLaid)
	if contract != nil {
		eggsLaidField = fmt.Sprintf("%s / %s",
			util.Numfmt(c.EggsLaid), util.NumfmtWhole(contract.UltimateGoal(c.IsElite())))
	}
	fmt.Fprintf(w, "Eggs laid:\t%s\n", eggsLaidField)
	layingRateField := util.Numfmt(c.EggsPerHour())
	if contract != nil {
		layingRateField = fmt.Sprintf("%s current / %s required",
			util.Numfmt(c.EggsPerHour()), util.Numfmt(c.RequiredEggsPerHour(contract)))
	}
	fmt.Fprintf(w, "Hourly laying rate:\t%s\n", layingRateField)
	timeToCompleteField := fmt.Sprintf("%s remaining",
		util.FormatDurationNonNegative(c.DurationUntilProductionDeadline()))
	if contract != nil {
		timeToCompleteField = fmt.Sprintf("%s expected / %s remaining",
			util.FormatDuration(c.ExpectedDurationUntilFinish(contract)),
			util.FormatDurationNonNegative(c.DurationUntilProductionDeadline()))
	}
	fmt.Fprintf(w, "Time to complete:\t%s\n", timeToCompleteField)
	w.Flush()
	fmt.Println()

	members := make([]*api.CoopStatus_Member, len(c.Members))
	copy(members, c.Members)
	sortBy.Sort(members)
	table := [][]string{
		{"Player", "Laid", "Rate/hr", "EB%", "Tokens"},
		{"------", "----", "-------", "---", "------"},
	}
	for _, m := range members {
		table = append(table, []string{
			m.Name, util.Numfmt(m.EggsLaid), util.Numfmt(m.EggsPerHour()),
			util.Numfmt(m.EarningBonusPercentage()), strconv.Itoa(int(m.Tokens)),
		})
	}
	table = append(table, table[1])
	table = append(table, []string{"Total", util.Numfmt(c.EggsLaid), util.Numfmt(c.EggsPerHour()), "", ""})
	util.PrintTable(table)
	fmt.Println()
}
