package web

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/fanaticscripter/EggContractor/contract"
	"github.com/fanaticscripter/EggContractor/coop"
	"github.com/fanaticscripter/EggContractor/db"
	"github.com/fanaticscripter/EggContractor/solo"
	"github.com/fanaticscripter/EggContractor/util"
)

type indexPayload struct {
	Errors   []error
	Warnings []string

	RefreshTime time.Time
	Statuses    []*SoloCoopStatus

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

// SoloCoopStatus is a unified type for solos and coops to facilitate sorting.
// For solos, Solo is set and Coop is nil; and vice versa.
type SoloCoopStatus struct {
	ContractId   string
	ContractName string
	IsSolo       bool
	Solo         *SoloStatus
	Coop         *CoopStatus
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

	// sort.Slice(wrappedSolos, func(i, j int) bool {
	// 	s1 := wrappedSolos[i]
	// 	s2 := wrappedSolos[j]
	// 	switch strings.Compare(s1.GetName(), s2.GetName()) {
	// 	case -1:
	// 		return true
	// 	case 1:
	// 		return false
	// 	default:
	// 		return s1.GetPlayerNickname() < s2.GetPlayerNickname()
	// 	}
	// })

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

	// sort.Slice(wrappedCoops, func(i, j int) bool {
	// 	c1 := wrappedCoops[i]
	// 	c2 := wrappedCoops[j]
	// 	switch strings.Compare(c1.Contract.Name, c2.Contract.Name) {
	// 	case -1:
	// 		return true
	// 	case 1:
	// 		return false
	// 	default:
	// 		return c1.Code < c2.Code
	// 	}
	// })

	statuses := make([]*SoloCoopStatus, 0, len(solos)+len(coops))
	for _, s := range wrappedSolos {
		statuses = append(statuses, &SoloCoopStatus{
			ContractId:   s.GetId(),
			ContractName: s.GetName(),
			IsSolo:       true,
			Solo:         s,
		})
	}
	for _, c := range wrappedCoops {
		statuses = append(statuses, &SoloCoopStatus{
			ContractId:   c.ContractId,
			ContractName: c.Contract.Name,
			IsSolo:       false,
			Coop:         c,
		})
	}
	sort.Slice(statuses, func(i, j int) bool {
		s1 := statuses[i]
		s2 := statuses[j]
		switch strings.Compare(s1.ContractName, s2.ContractName) {
		case -1:
			return true
		case 1:
			return false
		}
		switch {
		case s1.IsSolo && s2.IsSolo:
			// Both are solos, compare player nickname.
			return s1.Solo.GetPlayerNickname() < s2.Solo.GetPlayerNickname()
		case s1.IsSolo && !s2.IsSolo:
			// Coops before solos.
			return false
		case !s1.IsSolo && s2.IsSolo:
			// Coops before solos.
			return true
		case !s1.IsSolo && !s2.IsSolo:
			// Both are coops, compare coop code.
			return s1.Coop.Code < s2.Coop.Code
		}
		return false
	})

	contractIds := make([]string, len(statuses))
	for i, s := range statuses {
		contractIds[i] = s.ContractId
	}
	peeker, err := newPeekerPayloadFromPresetList(contractIds)
	if err != nil {
		errs = append(errs, err)
	}

	return &indexPayload{
		Errors:      errs,
		Warnings:    warnings,
		RefreshTime: timestamp,
		Statuses:    statuses,
		Peeker:      peeker,
	}
}

func (s *SoloStatus) OfflineAdjustedEggsLaid() float64 {
	return s.GetOfflineAdjustedEggsLaid(s.ClientRefreshTime)
}

func (s *SoloStatus) OfflineAdjustedExpectedDurationUntilFinish() time.Duration {
	return s.GetOfflineAdjustedExpectedDurationUntilFinish(s.ClientRefreshTime)
}

func (s *SoloStatus) ProgressInfo() *contract.ProgressInfo {
	return s.SoloContract.ProgressInfoWithProjection(s.OfflineAdjustedEggsLaid())
}

func (c *CoopStatus) OfflineAdjustedEggsLaid() float64 {
	return c.GetOfflineAdjustedEggsLaid(c.Activities)
}

func (c *CoopStatus) OfflineAdjustedExpectedDurationUntilFinish() time.Duration {
	return c.GetOfflineAdjustedExpectedDurationUntilFinish(c.Activities)
}

func (c *CoopStatus) ProgressInfo() *contract.ProgressInfo {
	return c.CoopStatus.ProgressInfoWithProjection(c.OfflineAdjustedEggsLaid())
}
