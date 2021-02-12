package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/fanaticscripter/EggContractor/api"
	"github.com/pkg/errors"
)

const _eiafxDataFile = "data/eiafx-data.json"

var _eiafxData *Store

type Store struct {
	Schema           string    `json:"$schema"`
	ArtifactFamilies []*Family `json:"artifact_families"`
}

type Family struct {
	CoreFamily

	Effect       string  `json:"effect"`
	EffectTarget string  `json:"effect_target"`
	Tiers        []*Tier `json:"tiers"`
}

type CoreFamily struct {
	Id          string                  `json:"id"`
	AfxId       api.ArtifactSpec_Name   `json:"afx_id"`
	Name        string                  `json:"name"`
	AfxType     api.ArtifactSpec_Type   `json:"afx_type"`
	Type        string                  `json:"type"`
	SortKey     uint32                  `json:"sort_key"`
	ChildAfxIds []api.ArtifactSpec_Name `json:"child_afx_ids"`
}

type Tier struct {
	Family *CoreFamily `json:"family"`

	CoreTier

	Quality               float64                   `json:"quality"`
	Craftable             bool                      `json:"craftable"`
	HasRarities           bool                      `json:"has_rarities"`
	PossibleAfxRarities   []api.ArtifactSpec_Rarity `json:"possible_afx_rarities"`
	HasEffects            bool                      `json:"has_effects"`
	AvailableFromMissions bool                      `json:"available_from_missions"`

	Effects []*Effect `json:"effects"`
	Recipe  *Recipe   `json:"recipe"`

	IngredientsAvailableFromMissions bool         `json:"ingredients_available_from_missions"`
	HardDependencies                 []Ingredient `json:"hard_dependencies"`
}

type CoreTier struct {
	ItemIdentifiers
	TierNumber   int                   `json:"tier_number"`
	TierName     string                `json:"tier_name"`
	AfxType      api.ArtifactSpec_Type `json:"afx_type"`
	Type         string                `json:"type"`
	IconFilename string                `json:"icon_filename"`
}

type ItemIdentifiers struct {
	Id       string                 `json:"id"`
	AfxId    api.ArtifactSpec_Name  `json:"afx_id"`
	AfxLevel api.ArtifactSpec_Level `json:"afx_level"`
	Name     string                 `json:"name"`
}

type Effect struct {
	AfxRarity    api.ArtifactSpec_Rarity `json:"afx_rarity"`
	Rarity       string                  `json:"rarity"`
	Effect       string                  `json:"effect"`
	EffectTarget string                  `json:"effect_target"`
	EffectSize   string                  `json:"effect_size"`
	FamilyEffect string                  `json:"family_effect"`
	// May be null (for stones).
	Slots *uint32 `json:"slots"`
}

type Recipe struct {
	Ingredients   []Ingredient  `json:"ingredients"`
	CraftingPrice CraftingPrice `json:"crafting_price"`
}

type Ingredient struct {
	CoreTier
	Count uint32 `json:"count"`
}

type CraftingPrice struct {
	Base    float64 `json:"base"`
	Low     float64 `json:"low"`
	Domain  uint32  `json:"domain"`
	Curve   float64 `json:"curve"`
	Initial uint32  `json:"initial"`
	Minimum uint32  `json:"minimum"`
}

func loadEiAfxData() error {
	body, err := ioutil.ReadFile(_eiafxDataFile)
	if err != nil {
		return errors.Wrapf(err, "error reading %s", _eiafxDataFile)
	}
	_eiafxData = &Store{}
	err = json.Unmarshal(body, _eiafxData)
	if err != nil {
		return errors.Wrapf(err, "error unmarshalling %s", _eiafxDataFile)
	}
	return nil
}

func getTier(afxId api.ArtifactSpec_Name, afxLevel api.ArtifactSpec_Level) *Tier {
	for _, f := range _eiafxData.ArtifactFamilies {
		for _, t := range f.Tiers {
			if t.AfxId == afxId && t.AfxLevel == afxLevel {
				return t
			}
		}
	}
	return nil
}
