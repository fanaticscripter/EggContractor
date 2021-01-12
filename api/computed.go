package api

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/fanaticscripter/EggContractor/util"
)

// It is said that elite contract is unlocked at 10T% EB.
// https://egg-inc.fandom.com/wiki/Contracts says so, and it certainly matches
// my personal experience. But it's hard to confirm.
const EliteEarningBonusThreshold = 1e11

func (f *FirstContact_Payload) AllContractProperties() []*ContractProperties {
	s := make([]*ContractProperties, 0)
	for _, c := range f.Contracts.ActiveContracts {
		s = append(s, c.Props)
	}
	for _, c := range f.Contracts.PastContracts {
		s = append(s, c.Props)
	}
	return s
}

func (c *Contract) StartedTime() time.Time {
	return util.DoubleToTime(c.Started)
}

func (c *Contract) ProductionDeadlineTime() time.Time {
	// The production_deadline field may not be available for solo contracts.
	if c.ProductionDeadline != 0 {
		return util.DoubleToTime(c.ProductionDeadline)
	} else if !c.StartedTime().IsZero() {
		return c.StartedTime().Add(c.Props.Duration())
	}
	return time.Time{}
}

func (c *Contract) CollectionDeadlineTime() time.Time {
	// The collection_deadline field may not be available for solo contracts.
	if c.CollectionDeadline != 0 {
		return util.DoubleToTime(c.CollectionDeadline)
	}
	prodDeadline := c.ProductionDeadlineTime()
	if prodDeadline.IsZero() {
		return prodDeadline.Add(48 * time.Hour)
	}
	return time.Time{}
}

func (c *ContractProperties) Duration() time.Duration {
	return util.DoubleToDuration(c.DurationSeconds)
}

func (c *ContractProperties) ExpiryTime() time.Time {
	return util.DoubleToTime(c.ExpiryTimestamp)
}

func (c *ContractProperties) UltimateGoal(isElite bool) float64 {
	ultimateGoals := make([]float64, 0, 3)
	if len(c.Rewards) > 0 {
		ultimateGoals = append(ultimateGoals, c.Rewards[len(c.Rewards)-1].Goal)
	}
	for _, t := range c.RewardTiers {
		if len(t.Rewards) > 0 {
			ultimateGoals = append(ultimateGoals, t.Rewards[len(t.Rewards)-1].Goal)
		}
	}
	var min, max float64
	for i, g := range ultimateGoals {
		if g > max {
			max = g
		}
		if i == 0 || g < min {
			min = g
		}
	}
	if isElite {
		return max
	} else {
		return min
	}
}

func (f *Farm) LastSavedTime() time.Time {
	return util.DoubleToTime(f.LastSaved)
}

func (c *CoopStatus) EggsPerSecond() float64 {
	var sum float64
	for _, m := range c.Members {
		sum += m.EggsPerSecond
	}
	return sum
}

func (c *CoopStatus) EggsPerHour() float64 {
	return 3600 * c.EggsPerSecond()
}

func (c *CoopStatus) DurationUntilProductionDeadline() time.Duration {
	return util.DoubleToDuration(c.SecondsUntilProductionDeadline)
}

func (c *CoopStatus) DurationUntilCollectionDeadline() time.Duration {
	return util.DoubleToDuration(c.SecondsUntilCollectionDeadline)
}

func (c *CoopStatus) IsElite() bool {
	var belowThresholdCnt, aboveThresholdCnt uint
	for _, m := range c.Members {
		if m.EarningBonus() >= EliteEarningBonusThreshold {
			aboveThresholdCnt++
		} else {
			belowThresholdCnt++
		}
	}
	// Ideally either one should be zero, but I can't be sure about the
	// threshold (in fact I can't even be sure the threshold is static), so play
	// it safe.
	return aboveThresholdCnt > belowThresholdCnt
}

func (c *CoopStatus) Creator() *CoopStatus_Member {
	for _, m := range c.Members {
		if m.Id == c.CreatorId {
			return m
		}
	}
	return nil
}

// RequiredEggsPerHour returns the laying rate required to complete the ultimate
// goal before the production deadline.
func (c *CoopStatus) RequiredEggsPerHour(contract *ContractProperties) float64 {
	eggsToLay := contract.UltimateGoal(c.IsElite()) - c.EggsLaid
	hoursLeft := c.DurationUntilProductionDeadline().Hours()
	if eggsToLay <= 0 || hoursLeft <= 0 {
		return 0
	} else {
		return eggsToLay / hoursLeft
	}
}

func (c *CoopStatus) ExpectedDurationUntilFinish(contract *ContractProperties) time.Duration {
	eggsToLay := contract.UltimateGoal(c.IsElite()) - c.EggsLaid
	if eggsToLay <= 0 {
		return 0
	} else if c.EggsPerSecond() <= 0 {
		return util.InfDuration // Forever
	} else {
		return util.DoubleToDuration(eggsToLay / c.EggsPerSecond())
	}
}

func (m *CoopStatus_Member) EggsPerHour() float64 {
	return 3600 * m.EggsPerSecond
}

func (m *CoopStatus_Member) EarningBonus() float64 {
	return math.Pow(10, m.EarningBonusOom)
}

func (m *CoopStatus_Member) EarningBonusPercentage() float64 {
	return m.EarningBonus() * 100
}

func (e EggType) Display() string {
	switch {
	case e == EggType_AI:
		return "AI"
	default:
		return strings.Title(strings.ReplaceAll(strings.ToLower(e.String()), "_", " "))
	}
}

func (e EggType) Value() float64 {
	switch e {
	case EggType_EDIBLE:
		return 0.1
	case EggType_SUPERFOOD:
		return 1.25
	case EggType_MEDICAL:
		return 6.25
	case EggType_ROCKET_FUEL:
		return 30
	case EggType_SUPER_MATERIAL:
		return 150
	case EggType_FUSION:
		return 700
	case EggType_QUANTUM:
		return 3_000
	case EggType_IMMORTALITY:
		return 12_500
	case EggType_TACHYON:
		return 50_000
	case EggType_GRAVITON:
		return 175_000
	case EggType_DILITHIUM:
		return 525_000
	case EggType_PRODIGY:
		return 1_500_000
	case EggType_TERRAFORM:
		return 10_000_000
	case EggType_ANTIMATTER:
		return 1e9
	case EggType_DARK_MATTER:
		return 1e11
	case EggType_AI:
		return 1e12
	case EggType_NEBULA:
		return 1.5e13
	case EggType_UNIVERSE:
		return 1e14
	case EggType_ENLIGHTENMENT:
		return 1e-7
	// Contract-only eggs.
	case EggType_CHOCOLATE:
		return 5
	case EggType_EASTER:
		return 0.05
	case EggType_WATERBALLOON:
		return 0.1
	case EggType_FIREWORK:
		return 4.99
	case EggType_PUMPKIN:
		return 0.99
	default:
		return 0
	}
}

func (e EggType) ValueDisplay() string {
	if e == EggType_ENLIGHTENMENT {
		return "0.00"
	}
	value := e.Value()
	if value >= 1e6 {
		return util.NumfmtWhole(value)
	}
	s := fmt.Sprintf("%f", value)
	if strings.Contains(s, ".") {
		// Trim unnecessary trailing zeros, and possibly also the decimal point
		// if the number turns out to be an integer.
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
	}
	return s
}
