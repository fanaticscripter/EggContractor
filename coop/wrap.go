package coop

import (
	"time"

	"github.com/fanaticscripter/EggContractor/api"
	"github.com/fanaticscripter/EggContractor/util"
)

// CoopStatus wraps api.CoopStatus and the corresponding api.ContractProperties.
type CoopStatus struct {
	*api.CoopStatus
	Contract *api.ContractProperties
}

// Contract is set to nil if no matching contract is found in the list.
func WrapCoopStatusWithContractList(c *api.CoopStatus, contracts []*api.ContractProperties) *CoopStatus {
	var contract *api.ContractProperties
	for _, cc := range contracts {
		if c.ContractId == cc.Id {
			contract = cc
			break
		}
	}
	return &CoopStatus{
		CoopStatus: c,
		Contract:   contract,
	}
}

func (c *CoopStatus) HasCompleted() bool {
	if c.Contract == nil {
		return false
	}
	return c.EggsLaid >= c.Contract.UltimateGoal(c.IsElite())
}

func (c *CoopStatus) HasNoTimeLeft() bool {
	return c.DurationUntilProductionDeadline() <= 0
}

func (c *CoopStatus) IsOnTrackToFinish() bool {
	if c.Contract == nil {
		return false
	}
	return c.EggsPerHour() >= c.RequiredEggsPerHour(c.Contract)
}

// OfflineAdjustedEggsLaid returns confirmed EggsLaid plus expected amount laid
// during each member's offline time. Offline time is capped at 30hr, the max
// away time from fully upgraded pro permit silos (although it is possible that
// a player refilled their silos in between without syncing up to the server, we
// shouldn't get our hopes up).
func (c *CoopStatus) GetOfflineAdjustedEggsLaid(activities map[string]*CoopMemberActivity) float64 {
	eggs := c.EggsLaid
	for _, m := range c.Members {
		activity, ok := activities[m.Id]
		if !ok {
			continue
		}
		offlineTime := activity.OfflineTime
		if offlineTime > 30*time.Hour {
			offlineTime = 30 * time.Hour
		}
		eggs += offlineTime.Hours() * activity.EggsPerHourSince
	}
	return eggs
}

// GetOfflineAdjustedExpectedDurationUntilFinish returns expected duration until
// finish assuming offline-adjusted eggs laid instead of confirmed eggs laid.
func (c *CoopStatus) GetOfflineAdjustedExpectedDurationUntilFinish(activities map[string]*CoopMemberActivity) time.Duration {
	eggsLaid := c.GetOfflineAdjustedEggsLaid(activities)
	eggsToLay := c.Contract.UltimateGoal(c.IsElite()) - eggsLaid
	if eggsToLay <= 0 {
		return 0
	} else if c.EggsPerSecond() <= 0 {
		return util.InfDuration // Forever
	} else {
		return util.DoubleToDuration(eggsToLay / c.EggsPerSecond())
	}
}
