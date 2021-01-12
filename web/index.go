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
	Solos       []*SoloStatus
	Coops       []*CoopStatus

	Peeker *peekerPayload
}

type SoloStatus struct {
	*solo.SoloContract
	ClientRefreshTime time.Time
}

type CoopStatus struct {
	*coop.CoopStatus
	Activities map[string]*coop.CoopMemberActivity
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

	wrappedSolos := make([]*SoloStatus, len(solos))
	for i, s := range solos {
		wrappedSolos[i] = &SoloStatus{
			SoloContract:      s,
			ClientRefreshTime: timestamp,
		}
	}

	wrappedCoops := make([]*CoopStatus, len(coops))
	for i, c := range coops {
		activities, err := db.GetCoopMemberActivityStats(c, timestamp)
		if err != nil {
			errs = append(errs, err)
			activities = nil
		}
		wrappedCoops[i] = &CoopStatus{
			CoopStatus: c,
			Activities: activities,
		}
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
		Solos:       wrappedSolos,
		Coops:       wrappedCoops,
		Peeker:      peeker,
	}
}

func (s *SoloStatus) OfflineAdjustedEggsLaid() float64 {
	return s.GetOfflineAdjustedEggsLaid(s.ClientRefreshTime)
}

func (s *SoloStatus) OfflineAdjustedExpectedDurationUntilFinish() time.Duration {
	return s.GetOfflineAdjustedExpectedDurationUntilFinish(s.ClientRefreshTime)
}

func (c *CoopStatus) OfflineAdjustedEggsLaid() float64 {
	return c.GetOfflineAdjustedEggsLaid(c.Activities)
}
