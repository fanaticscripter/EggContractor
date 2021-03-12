// Convert consumption-data.json to payload suitable for the app.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"

	"github.com/fanaticscripter/EggContractor/api"
	"github.com/fanaticscripter/EggContractor/port/wasm/_common/eiafx"
)

const _appDataFile = "src/app-data.json"

type AppPayload struct {
	Families []*AppFamily `json:"families"`
}

type AppFamily struct {
	eiafx.CoreFamily

	Tiers []*AppTier `json:"tiers"`
}

type AppTier struct {
	eiafx.CoreTier

	Rarities []ConsumptionOutcome `json:"rarities"`
	Sources  []Source             `json:"sources"`
}

type Source struct {
	eiafx.CoreTier
	AfxRarity     api.ArtifactSpec_Rarity `json:"afx_rarity"`
	Rarity        string                  `json:"rarity"`
	Deterministic bool                    `json:"deterministic"`
	ExpectedYield float64                 `json:"expected_yield"`
}

func main() {
	if err := eiafx.LoadData(); err != nil {
		log.Fatal(err)
	}
	if err := loadConsumptionData(); err != nil {
		log.Fatal(err)
	}
	for i, c := range _consumptionOutcomes {
		c.Complete()
		_consumptionOutcomes[i] = c
	}

	payload := &AppPayload{}
	for _, f := range eiafx.Data.ArtifactFamilies {
		af := &AppFamily{
			CoreFamily: f.CoreFamily,
		}
		for _, t := range f.Tiers {
			at := &AppTier{
				CoreTier: t.CoreTier,
			}
			for _, c := range _consumptionOutcomes {
				if c.Item.AfxId == t.AfxId && c.Item.AfxLevel == t.AfxLevel {
					at.Rarities = append(at.Rarities, c)
				}
			}
			numRarities := len(at.Rarities)
			expectedNumRarities := 1
			if t.HasRarities {
				expectedNumRarities = len(t.PossibleAfxRarities)
			}
			if numRarities != expectedNumRarities {
				panic(fmt.Sprintf("expected %d rarities, got %d rarities for %+v", expectedNumRarities, numRarities, t))
			}
			af.Tiers = append(af.Tiers, at)
		}
		payload.Families = append(payload.Families, af)
	}

	for _, f := range payload.Families {
		for _, t := range f.Tiers {
			for _, c := range t.Rarities {
				for _, bp := range c.ExpectedByproducts {
				LocateTarget:
					for _, ff := range payload.Families {
						for _, tt := range ff.Tiers {
							if bp.AfxId == tt.AfxId && bp.AfxLevel == tt.AfxLevel {
								tt.Sources = append(tt.Sources, Source{
									CoreTier:      t.CoreTier,
									AfxRarity:     c.Item.AfxRarity,
									Rarity:        c.Item.Rarity,
									Deterministic: c.Deterministic,
									ExpectedYield: bp.ExpectedCount,
								})
								break LocateTarget
							}
						}
					}
				}
			}
		}
	}

	encoded, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		log.Fatalf("error serializing app payload: %s", err)
	}
	encoded = append(encoded, '\n')
	err = ioutil.WriteFile(_appDataFile, encoded, 0o644)
	if err != nil {
		log.Fatalf("error writing to %s: %s", _appDataFile, err)
	}
}
