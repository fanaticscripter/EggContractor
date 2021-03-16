package artifacts

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/fanaticscripter/EggContractor/api"
)

// effect_numbers.go is generated.
//go:generate python3 effect_numbers.py

// Assumes non-enlightenment farm, for now.
func MultiplicativeEffect(artifacts []*api.CompleteArtifact, afxIds []api.ArtifactSpec_Name) float64 {
	effect := 1.0

	examineItem := func(spec *api.ArtifactSpec) {
		for _, id := range afxIds {
			if spec.Name == id {
				delta, err := effectDelta(spec)
				if err != nil {
					log.Warn(err)
					return
				}
				effect *= (1 + delta)
				return
			}
		}
	}

	for _, a := range artifacts {
		if a == nil {
			continue
		}
		examineItem(a.Spec)
		for _, stone := range a.Stones {
			if stone == nil {
				continue
			}
			examineItem(stone)
		}
	}

	return effect
}

func LayingRateEffect(artifacts []*api.CompleteArtifact) float64 {
	return MultiplicativeEffect(artifacts, []api.ArtifactSpec_Name{
		api.ArtifactSpec_QUANTUM_METRONOME,
		api.ArtifactSpec_TACHYON_STONE,
	})
}

// A stripped down ArtifactSpec that is comparable.
type item struct {
	Name   api.ArtifactSpec_Name
	Level  api.ArtifactSpec_Level
	Rarity api.ArtifactSpec_Rarity
}

func newItem(a *api.ArtifactSpec) item {
	return item{
		Name:   a.Name,
		Level:  a.Level,
		Rarity: a.Rarity,
	}
}

func effectDelta(a *api.ArtifactSpec) (float64, error) {
	item := newItem(a)
	delta, ok := _effectDeltas[item]
	if !ok {
		return 0, errors.Errorf("cannot find effect delta for %+v", item)
	}
	return delta, nil
}
