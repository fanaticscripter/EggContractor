package contract

import "github.com/fanaticscripter/EggContractor/api"

type ProgressInfo struct {
	EggsLaid          float64
	ProjectedEggsLaid float64
	Rewards           []*Reward
	UltimateGoal      float64
}

type Reward struct {
	*api.Reward
	PercentageOfUltimateGoal float64
	PercentageCompleted      float64
}

func NewProgressInfo(
	rewards []*api.Reward,
	eggsLaid float64,
	projectedEggsLaid float64,
) *ProgressInfo {
	if len(rewards) == 0 {
		return nil
	}
	ultimateGoal := rewards[len(rewards)-1].Goal
	if ultimateGoal == 0 {
		panic("NewProgressInfo: ultimate goal is zero")
	}
	wrappedRewards := make([]*Reward, 0)
	for _, r := range rewards {
		if r.Goal == 0 {
			panic("NewProgressInfo: reward goal is zero")
		}
		wrappedRewards = append(wrappedRewards, &Reward{
			Reward:                   r,
			PercentageOfUltimateGoal: r.Goal / ultimateGoal * 100,
			PercentageCompleted:      eggsLaid / r.Goal * 100,
		})
	}
	return &ProgressInfo{
		EggsLaid:          eggsLaid,
		ProjectedEggsLaid: projectedEggsLaid,
		Rewards:           wrappedRewards,
		UltimateGoal:      ultimateGoal,
	}
}

func (p *ProgressInfo) PercentageCompleted() float64 {
	return p.EggsLaid / p.UltimateGoal * 100
}

func (p *ProgressInfo) ProjectedPercentageCompleted() float64 {
	return p.ProjectedEggsLaid / p.UltimateGoal * 100
}
