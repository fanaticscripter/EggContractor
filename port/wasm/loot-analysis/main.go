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
	"github.com/fanaticscripter/EggContractor/util"
)

const _appDataFile = "src/app-data.json"

type payload struct {
	Items []item    `json:"items"`
	Stats lootStats `json:"stats"`
}

type item struct {
	Id             string                  `json:"id"`
	Name           string                  `json:"name"`
	Tier           eiafx.CoreTier          `json:"tier"`
	AfxRarity      api.ArtifactSpec_Rarity `json:"afxRarity"`
	Rarity         string                  `json:"rarity"`
	Quality        float64                 `json:"quality"`
	OddsMultiplier float64                 `json:"oddsMultiplier"`

	categoryIndex int
}

const (
	_indexCategoryArtifactsCommon = iota
	_indexCategoryArtifactsRare
	_indexCategoryArtifactsEpic
	_indexCategoryArtifactsLegendary
	_indexCategoryStones
	_indexCategoryIngredients
	_indexCategoryStoneIngredients
	_numCategories
)

var _categoryNames = []string{
	"Artifacts (Common)",
	"Artifacts (Rare)",
	"Artifacts (Epic)",
	"Artifacts (Legendary)",
	"Stones",
	"Ingredients",
	"Stone fragments",
}

type lootStats []missionLootStats

type missionLootStats struct {
	Info         missionInfo                       `json:"info"`
	MissionCount int                               `json:"missionCount"`
	Categories   [_numCategories]categoryLootStats `json:"categories"`
}

type missionInfo struct {
	Id              string  `json:"id"`
	Display         string  `json:"display"`
	ShipIconPath    string  `json:"shipIconPath"`
	Capacity        int     `json:"capacity"`
	DurationSeconds float64 `json:"durationSeconds"`
	DurationDisplay string  `json:"durationDisplay"`
	Quality         float64 `json:"quality"`
	MinQuality      float64 `json:"minQuality"`
	MaxQuality      float64 `json:"maxQuality"`
}

type categoryLootStats struct {
	CategoryName string     `json:"categoryName"`
	Stats        []itemStat `json:"stats"`
}

type itemStat struct {
	ItemId string `json:"itemId"`
	Count  int    `json:"count"`
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

	items := make([]item, 0)
	for _, p := range eiafx.Config.ArtifactParameters {
		artifact := p.Spec
		tier, err := eiafx.GetTier(artifact)
		if err != nil {
			log.Fatal(err)
		}
		afxRarity := artifact.Rarity
		rarity := afxRarity.Display()
		id := tier.Id
		name := fmt.Sprintf("%s (T%d)", tier.Name, tier.TierNumber)
		if tier.AfxType == api.ArtifactSpec_ARTIFACT {
			id += ":" + strings.ToLower(rarity)
			if afxRarity != api.ArtifactSpec_COMMON {
				name += ", " + rarity
			}
		}
		categoryIndex := -1
		switch tier.AfxType {
		case api.ArtifactSpec_ARTIFACT:
			switch afxRarity {
			case api.ArtifactSpec_COMMON:
				categoryIndex = _indexCategoryArtifactsCommon
			case api.ArtifactSpec_RARE:
				categoryIndex = _indexCategoryArtifactsRare
			case api.ArtifactSpec_EPIC:
				categoryIndex = _indexCategoryArtifactsEpic
			case api.ArtifactSpec_LEGENDARY:
				categoryIndex = _indexCategoryArtifactsLegendary
			}
		case api.ArtifactSpec_STONE:
			categoryIndex = _indexCategoryStones
		case api.ArtifactSpec_INGREDIENT:
			categoryIndex = _indexCategoryIngredients
		case api.ArtifactSpec_STONE_INGREDIENT:
			categoryIndex = _indexCategoryStoneIngredients
		}
		if categoryIndex < 0 {
			log.Fatalf("the impossible happened: failed to categorize %+v", artifact)
		}
		items = append(items, item{
			Id:             id,
			Name:           name,
			Tier:           tier.CoreTier,
			AfxRarity:      afxRarity,
			Rarity:         rarity,
			Quality:        tier.Quality,
			OddsMultiplier: p.OddsMultiplier,
			categoryIndex:  categoryIndex,
		})
	}

	stats := make(lootStats, 0)
	for _, s := range eiafx.Config.MissionParameters {
		ship := s.Ship
		for _, duration := range s.Durations {
			typ := duration.DurationType
			if typ == api.MissionInfo_TUTORIAL {
				continue
			}
			info := missionInfo{
				Id:              missionId(ship, typ),
				Display:         ship.Name() + ", " + typ.Display(),
				ShipIconPath:    "egginc/" + ship.IconFilename(),
				Capacity:        int(duration.Capacity),
				DurationSeconds: duration.Seconds,
				DurationDisplay: util.FormatDurationWhole(util.DoubleToDuration(duration.Seconds)),
				Quality:         float64(duration.Quality),
				MinQuality:      float64(duration.MinQuality),
				MaxQuality:      float64(duration.MaxQuality),
			}
			missionData := loot.Data.MissionLoot(ship, typ)
			missionCount := missionData.TotalArtifactsCount / info.Capacity
			var categories [_numCategories]categoryLootStats
			for i := range categories {
				categories[i].CategoryName = _categoryNames[i]
				categories[i].Stats = []itemStat{}
			}
			for _, it := range items {
				if info.MinQuality <= it.Quality && it.Quality <= info.MaxQuality {
					categories[it.categoryIndex].Stats = append(categories[it.categoryIndex].Stats, itemStat{
						ItemId: it.Id,
						Count:  missionData.ItemRarityCount(it.Tier.AfxId, it.Tier.AfxLevel, it.AfxRarity),
					})
				}
			}
			stats = append(stats, missionLootStats{
				Info:         info,
				MissionCount: missionCount,
				Categories:   categories,
			})
		}
	}

	appData := &payload{
		Items: items,
		Stats: stats,
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
