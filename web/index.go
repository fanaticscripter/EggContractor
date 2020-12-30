package web

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/fanaticscripter/EggContractor/coop"
	"github.com/fanaticscripter/EggContractor/db"
	"github.com/fanaticscripter/EggContractor/solo"
	"github.com/fanaticscripter/EggContractor/util"
)

type indexPayload struct {
	Errors   []error
	Warnings []string

	RefreshTime time.Time
	Solos       []*solo.SoloContract
	Coops       []*coop.CoopStatus

	Peeker *peekerPayload
}

// GET /?by=<timestamp>
func indexHandler(c echo.Context) error {
	byThisTime := time.Now()
	by := c.QueryParam("by")
	byTimestamp, err := strconv.ParseFloat(by, 64)
	if err == nil {
		byThisTime = util.DoubleToTime(byTimestamp)
	}

	payload := getIndexPayload(byThisTime)
	return c.Render(http.StatusOK, "index.html", payload)
}

func getIndexPayload(byThisTime time.Time) *indexPayload {
	errs := make([]error, 0)
	warnings := make([]string, 0)
	timestamp, solos, coops, err := db.GetSoloAndCoopStatusesFromRefresh(byThisTime)
	if err != nil {
		errs = append(errs, err)
		return &indexPayload{
			Errors: errs,
			Peeker: &peekerPayload{},
		}
	}
	if timestamp.IsZero() {
		warnings = append(warnings,
			"no refresh found in the database, try using the refresh subcommand of EggContractor")
	} else if len(solos) == 0 && len(coops) == 0 {
		warnings = append(warnings, util.MsgNoActiveContracts)
	}

	contractIds := make([]string, 0)
	for _, c := range solos {
		contractIds = append(contractIds, c.GetId())
	}
	for _, c := range coops {
		contractIds = append(contractIds, c.ContractId)
	}
	peeker, err := newPeekerPayloadFromPresetList(contractIds)
	if err != nil {
		errs = append(errs, err)
	}

	return &indexPayload{
		Errors:      errs,
		Warnings:    warnings,
		RefreshTime: timestamp,
		Solos:       solos,
		Coops:       coops,
		Peeker:      peeker,
	}
}
