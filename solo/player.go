package solo

import (
	"math"

	"github.com/fanaticscripter/EggContractor/api"
)

type Player struct {
	Progress *api.FirstContact_Payload_Progress
}

func GetPlayer(fc *api.FirstContact) *Player {
	return &Player{
		Progress: fc.Data.Progress,
	}
}

func (p *Player) EarningBonusPercentage() float64 {
	soulEggBonusPercentage := 10.0
	for _, r := range p.Progress.EpicResearches {
		if r.Id == "soul_eggs" {
			soulEggBonusPercentage += float64(r.Level)
			break
		}
	}
	prophecyEggBonusPercentage := 5.0
	for _, r := range p.Progress.EpicResearches {
		if r.Id == "prophecy_bonus" {
			prophecyEggBonusPercentage += float64(r.Level)
		}
	}
	soulEggBonusPercentage *= math.Pow(1+prophecyEggBonusPercentage/100, float64(p.Progress.ProphecyEggs))
	return soulEggBonusPercentage * p.Progress.SoulEggs
}

func (p *Player) EarningBonus() float64 {
	return p.EarningBonusPercentage() / 100
}
