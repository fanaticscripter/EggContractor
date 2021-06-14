package web

import (
	"html/template"
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
	Warnings []template.HTML

	RefreshTime          time.Time
	Statuses             []*SoloCoopStatus
	ContractFilterConfig ContractFilterConfig
	HideSolos            bool
	HideFull             bool

	Peeker *peekerPayload
}

type SoloStatus struct {
	*solo.SoloContract
	ClientRefreshTime time.Time
	Filtered          bool
}

type CoopStatus struct {
	*coop.CoopStatus
	Activities map[string]*coop.CoopMemberActivity
	Filtered   bool
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

type ContractFilterConfig struct {
	Contracts []struct {
		Id   string
		Name string
	}
	Filter string
}

// GET /?by=<timestamp>
func indexHandler(c echo.Context) error {
	byThisTime := time.Now()
	by := c.QueryParam("by")
	byTimestamp, err := strconv.ParseFloat(by, 64)
	if err == nil {
		byThisTime = util.DoubleToTime(byTimestamp)
	}

	contractFilter := c.QueryParam("c")
	hideSolos := c.QueryParam("hide_solos") != ""
	hideFull := c.QueryParam("hide_full") != ""

	payload := getIndexPayload(byThisTime, contractFilter, hideSolos, hideFull)
	return c.Render(http.StatusOK, "index.html", payload)
}

func getIndexPayload(byThisTime time.Time, contractFilter string, hideSolos bool, hideFull bool) *indexPayload {
	errs := make([]error, 0)
	warnings := make([]template.HTML, 0)
	timestamp, solos, coops, err := db.GetSoloAndCoopStatusesFromRefresh(byThisTime, true)
	if err != nil {
		errs = append(errs, err)
		return &indexPayload{
			Errors: errs,
			Peeker: &peekerPayload{},
		}
	}
	if _configDeprecations.HasLegacyPlayerField {
		warnings = append(warnings, template.HTML(`Config key <code>player</code> has been deprecated
		in favor of <code>players</code> since the introduction of multi-account support. See
		<a href="https://github.com/fanaticscripter/EggContractor/wiki/Multi-account-migration" target="_blank" class="text-blue-500 hover:text-blue-400">this guide</a>
		for how to migrate and take advantage of new features.`))
	}
	if timestamp.IsZero() {
		warnings = append(warnings,
			"no refresh found in the database, try using the refresh subcommand of EggContractor")
	} else if len(solos) == 0 && len(coops) == 0 {
		warnings = append(warnings, util.HTMLMsgNoActiveContracts)
	}

	isFiltered := func(contractId string) bool { return false }
	contractFilterIsValid := false
	if contractFilter != "" {
		// Make sure contract filter is valid -- do not filter if there's nothing
		// matching it.
		for _, s := range solos {
			if contractFilter == s.GetId() {
				contractFilterIsValid = true
				break
			}
		}
		for _, c := range coops {
			if contractFilter == c.ContractId {
				contractFilterIsValid = true
				break
			}
		}
		if contractFilterIsValid {
			isFiltered = func(contractId string) bool { return contractId != contractFilter }
		}
	}

	wrappedSolos := make([]*SoloStatus, len(solos))
	for i, s := range solos {
		wrappedSolos[i] = &SoloStatus{
			SoloContract:      s,
			ClientRefreshTime: timestamp,
			Filtered:          hideSolos || isFiltered(s.GetId()),
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
			Filtered:   isFiltered(c.ContractId),
		}
	}

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

	filterConfig := ContractFilterConfig{}
	if contractFilterIsValid {
		filterConfig.Filter = contractFilter
	}
	lastContractId := ""
	for _, s := range statuses {
		if s.ContractId != lastContractId {
			filterConfig.Contracts = append(filterConfig.Contracts, struct {
				Id   string
				Name string
			}{
				Id:   s.ContractId,
				Name: s.ContractName,
			})
			lastContractId = s.ContractId
		}
	}

	contractIds := make([]string, len(statuses))
	for i, s := range statuses {
		contractIds[i] = s.ContractId
	}
	peeker, err := newPeekerPayloadFromPresetList(contractIds)
	if err != nil {
		errs = append(errs, err)
	}

	return &indexPayload{
		Errors:               errs,
		Warnings:             warnings,
		RefreshTime:          timestamp,
		Statuses:             statuses,
		ContractFilterConfig: filterConfig,
		HideSolos:            hideSolos,
		HideFull:             hideFull,
		Peeker:               peeker,
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

// Tests whether a SoloStatus/CoopStatus is filtered. Returns false for any
// other type.
func statusIsFiltered(s interface{}) bool {
	ss, ok := s.(*SoloStatus)
	if ok {
		return ss.Filtered
	}
	cc, ok := s.(*CoopStatus)
	if ok {
		return cc.Filtered
	}
	return false
}
