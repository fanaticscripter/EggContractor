package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/fanaticscripter/EggContractor/api"
	"github.com/fanaticscripter/EggContractor/port/wasm/_common/eiafx"
	"github.com/fanaticscripter/EggContractor/port/wasm/_common/loot"
)

const _appDataFile = "src/app-data.json"

type payload struct {
	Items    []item    `json:"items"`
	Missions []mission `json:"missions"`
}

type item struct {
	eiafx.CoreTier
	Display  string       `json:"display"`
	IconPath string       `json:"iconPath"`
	Recipe   []ingredient `json:"recipe"`

	possibleAfxRarities []api.ArtifactSpec_Rarity
}

type ingredient struct {
	Id    string `json:"id"`
	Count uint32 `json:"count"`
}

type mission struct {
	Id              string     `json:"id"`
	Display         string     `json:"display"`
	IconPath        string     `json:"iconPath"`
	Capacity        int        `json:"capcity"`
	DurationSeconds float64    `json:"durationSeconds"`
	Count           int        `json:"count"`
	LootTotal       int        `json:"lootTotal"`
	Loot            []lootItem `json:"loot"`
}

type lootItem struct {
	Id    string `json:"id"`
	Count int    `json:"count"`
}

func main() {
	err := eiafx.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = eiafx.LoadData()
	if err != nil {
		log.Fatal(err)
	}

	err = loot.LoadData()
	if err != nil {
		log.Fatal(err)
	}

	items := []item{}
	for _, p := range eiafx.Config.ArtifactParameters {
		artifact := p.Spec
		if artifact.Rarity > 0 {
			continue
		}
		tier, err := eiafx.GetTier(artifact)
		if err != nil {
			log.Fatal(err)
		}
		var recipe []ingredient
		if tier.Craftable {
			for _, i := range tier.Recipe.Ingredients {
				recipe = append(recipe, ingredient{
					Id:    i.Id,
					Count: i.Count,
				})
			}
		}
		possibleAfxRarities := tier.PossibleAfxRarities
		if possibleAfxRarities == nil {
			possibleAfxRarities = []api.ArtifactSpec_Rarity{api.ArtifactSpec_COMMON}
		}
		items = append(items, item{
			CoreTier:            tier.CoreTier,
			Display:             fmt.Sprintf("%s (T%d)", tier.Name, tier.TierNumber),
			IconPath:            "egginc/" + tier.IconFilename,
			Recipe:              recipe,
			possibleAfxRarities: possibleAfxRarities,
		})
	}

	missions := []mission{}
	for _, s := range eiafx.Config.MissionParameters {
		ship := s.Ship
		for _, duration := range s.Durations {
			typ := duration.DurationType
			if typ == api.MissionInfo_TUTORIAL {
				continue
			}
			m := mission{
				Id:              missionId(ship, typ),
				Display:         ship.Name() + ", " + typ.Display(),
				IconPath:        "egginc/" + ship.IconFilename(),
				Capacity:        int(duration.Capacity),
				DurationSeconds: duration.Seconds,
			}

			missionData := loot.Data.MissionLoot(ship, typ)
			lootTotal := missionData.TotalArtifactsCount
			lootTotalTally := 0
			lootItems := []lootItem{}
			for _, item := range items {
				count, err := missionData.ItemCount(item.AfxId, item.AfxLevel, item.possibleAfxRarities)
				if err != nil {
					log.Fatal(err)
				}
				if count.Total == 0 {
					continue
				}
				lootItems = append(lootItems, lootItem{
					Id:    item.Id,
					Count: count.Total,
				})
				lootTotalTally += count.Total
			}
			if lootTotal != lootTotalTally {
				log.Fatalf("%s: wrong loot total: reported %d, actual tally %d", m.Display, lootTotal, lootTotalTally)
			}

			m.Count = lootTotal / m.Capacity
			m.LootTotal = lootTotal
			m.Loot = lootItems
			missions = append(missions, m)
		}
	}

	appData := &payload{
		Items:    items,
		Missions: missions,
	}

	encoded, err := json.MarshalIndent(appData, "", "  ")
	if err != nil {
		log.Fatalf("error serializing app payload: %s", err)
	}
	encoded = append(encoded, '\n')
	err = ioutil.WriteFile(_appDataFile, encoded, 0o644)
	if err != nil {
		log.Fatalf("error writing to %s: %s", _appDataFile, err)
	}
}

func missionId(ship api.MissionInfo_Spaceship, durationType api.MissionInfo_DurationType) string {
	return strings.ToLower(strings.ReplaceAll(ship.Name()+" "+durationType.Display(), " ", "-"))
}
