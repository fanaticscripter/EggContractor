package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"

	"github.com/fanaticscripter/EggContractor/api"
)

const _lootDataFile = "src/mission_reward_count.json"

var _lootData LootStore

type (
	ShipName         string
	DurationTypeName string
	ArtifactName     string
	RarityName       string
)

type LootStore map[ShipName]ShipLootStore

type ShipLootStore map[DurationTypeName]*MissionLootStore

type MissionLootStore struct {
	TotalArtifactsCount int                                   `json:"count"`
	LootCounts          map[ArtifactName][]map[RarityName]int `json:"rewards"`
}

type ItemCount struct {
	Total    int                             `json:"total"`
	Rarities map[api.ArtifactSpec_Rarity]int `json:"rarities"`
}

func (s LootStore) MissionLoot(
	ship api.MissionInfo_Spaceship,
	durationType api.MissionInfo_DurationType,
) *MissionLootStore {
	shipName := ShipName(strings.ToLower(ship.String()))
	durationTypeName := DurationTypeName(strings.ToLower(durationType.String()))
	return s[shipName][durationTypeName]
}

func (s *MissionLootStore) ItemCount(
	afxId api.ArtifactSpec_Name,
	afxLevel api.ArtifactSpec_Level,
	possibleAfxRarities []api.ArtifactSpec_Rarity,
) (*ItemCount, error) {
	artifactName := ArtifactName(strings.ToLower(afxId.String()))
	counts := s.LootCounts[artifactName][afxLevel]
	rarities := make(map[api.ArtifactSpec_Rarity]int)
	total := 0
	for rarityName, count := range counts {
		afxRarity := rarityName.AfxRarity()
		possible := false
		for _, r := range possibleAfxRarities {
			if r == afxRarity {
				possible = true
				rarities[afxRarity] = count
				total += count
				break
			}
		}
		if !possible && count > 0 {
			a := api.ArtifactSpec{
				Name:  afxId,
				Level: afxLevel,
			}
			return nil, fmt.Errorf("impossible item found in loot data: %s %s, x%d",
				afxRarity, a.GameName(), count)
		}
	}
	return &ItemCount{
		Total:    total,
		Rarities: rarities,
	}, nil
}

func (r RarityName) AfxRarity() api.ArtifactSpec_Rarity {
	s := strings.ToUpper(string(r))
	return api.ArtifactSpec_Rarity(api.ArtifactSpec_Rarity_value[s])
}

func loadLootData() error {
	body, err := ioutil.ReadFile(_lootDataFile)
	if err != nil {
		return errors.Wrapf(err, "error reading %s", _lootDataFile)
	}
	err = json.Unmarshal(body, &_lootData)
	if err != nil {
		return errors.Wrapf(err, "error unmarshalling %s", _lootDataFile)
	}
	return nil
}
