package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/fanaticscripter/EggContractor/util"
)

func dumpContractDBToCSV(csvpath string) {
	contracts, err := getContractsFromDB()
	if err != nil {
		log.Error(err)
		return
	}
	f, err := os.Create(csvpath)
	if err != nil {
		log.Errorf("error creating %#v: %s", csvpath, err)
	}
	w := csv.NewWriter(f)
	if err := w.Write([]string{
		"ID",
		"Name",
		"Type",
		"Has Leggacy",
		"Offering (Estimated)",
		"Expiry",
		"Egg",
		"Duration",
		"Size",
		"Token",
		"Std#1",
		"Std#2",
		"Std#3",
		"Std Rate/hr",
		"Elt#1",
		"Elt#2",
		"Elt#3",
		"Elt Rate/hr",
		"#PE",
		"Std PE Rate/hr",
		"Elt PE Rate/hr",
		"JSON",
		"Protobuf (base64)",
	}); err != nil {
		log.Errorf("error writing header to %#v: %s", csvpath, err)
	}
	for _, c := range contracts {
		offeringTime := c.EstimatedOfferingTime().In(time.UTC).Format("2006-01-02Z")
		expiryTime := c.ExpiryTime().In(time.UTC).Format(time.RFC3339)
		if c.Id == "first-contract" {
			offeringTime = "-"
			expiryTime = "-"
		}
		maxCoopSize := "-"
		if c.MaxCoopSize > 1 {
			maxCoopSize = fmt.Sprintf("%d", c.MaxCoopSize)
		}
		tokenInterval := "-"
		if c.TokenInterval != 0 {
			tokenInterval = fmt.Sprintf("%.0fm", c.TokenInterval)
		}
		stdGoals := c.StandardGoalsStr()
		eltGoals := c.EliteGoalsStr()
		stdUltGoal := c.StandardUltimateGoal()
		eltUltGoal := c.EliteUltimateGoal()
		stdRate := ""
		eltRate := ""
		if stdUltGoal > 0 {
			stdRate = util.Numfmt(stdUltGoal / c.Duration().Hours())
		}
		if eltUltGoal > 0 {
			eltRate = util.Numfmt(eltUltGoal / c.Duration().Hours())
		}
		peCnt := c.ProphecyEggCount()
		stdPEGoal := c.StandardProphecyEggGoal()
		eltPEGoal := c.EliteProphecyEggGoal()
		stdPERate := ""
		eltPERate := ""
		if stdPEGoal > 0 {
			stdPERate = util.Numfmt(stdPEGoal / c.Duration().Hours())
		}
		if eltPEGoal > 0 {
			eltPERate = util.Numfmt(eltPEGoal / c.Duration().Hours())
		}
		jsonb, err := c.JSON()
		if err != nil {
			log.Errorf("error serializing contract %s to JSON: %s", c.Id, err)
			continue
		}
		protob, err := c.B64Protobuf()
		if err != nil {
			log.Errorf("error serializing contract %s to base64-encoded protobuf: %s", c.Id, err)
			continue
		}
		if err := w.Write([]string{
			c.Id,
			c.Name,
			c.Type(),
			fmt.Sprintf("%t", c.HasLeggacy),
			offeringTime,
			expiryTime,
			c.EggType.Display(),
			util.FormatDurationWhole(c.Duration()),
			maxCoopSize,
			tokenInterval,
			stdGoals[0],
			stdGoals[1],
			stdGoals[2],
			stdRate,
			eltGoals[0],
			eltGoals[1],
			eltGoals[2],
			eltRate,
			fmt.Sprintf("%d", peCnt),
			stdPERate,
			eltPERate,
			string(jsonb),
			string(protob),
		}); err != nil {
			log.Errorf("error writing record to %#v: %s", csvpath, err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		log.Errorf("error flushing records to %#v: %s", csvpath, err)
	}
}
