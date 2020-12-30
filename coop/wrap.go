package coop

import (
	"github.com/fanaticscripter/EggContractor/api"
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

func (c *CoopStatus) Display(sortBy api.By) {
	c.CoopStatus.Display(sortBy, c.Contract)
}
