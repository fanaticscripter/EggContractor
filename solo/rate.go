package solo

import "github.com/fanaticscripter/EggContractor/api"

type layingRateResearch struct {
	Id        string
	PerLevel  float64
	MaxLevels int
}

const _baseLayingRate float64 = 1.0 / 30 // 1 egg per 30 seconds

var _layingRateCommonResearches = []layingRateResearch{
	{"comfy_nests", 0.10, 50},
	{"hen_house_ac", 0.05, 50},
	{"improved_genetics", 0.15, 30},
	{"time_compress", 0.10, 20},
	{"timeline_diversion", 0.02, 50},
	{"relativity_optimization", 0.10, 10},
}

var _layingRateEpicResearches = []layingRateResearch{
	{"epic_egg_laying", 0.05, 20},
}

func eggsPerChickenPerSecond(commonResearches []*api.Research, epicResearches []*api.EpicResearch) float64 {
	rate := _baseLayingRate
	for _, r := range _layingRateCommonResearches {
		for _, rr := range commonResearches {
			if r.Id == rr.Id {
				rate *= 1 + float64(rr.Level)*r.PerLevel
				break
			}
		}
	}
	for _, r := range _layingRateEpicResearches {
		for _, rr := range epicResearches {
			if r.Id == rr.Id {
				rate *= 1 + float64(rr.Level)*r.PerLevel
				break
			}
		}
	}
	return rate
}
