package eiafx

import (
	_ "embed"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/fanaticscripter/EggContractor/api"
)

//go:embed eiafx-data.json
var _eiafxDataJSON []byte

var Data *Store

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
	BaseCraftingPrices    []float64                 `json:"base_crafting_prices"`
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

func LoadData() error {
	Data = &Store{}
	err := json.Unmarshal(_eiafxDataJSON, Data)
	if err != nil {
		return errors.Wrap(err, "error unmarshalling eiafx-data.json")
	}
	return nil
}

func GetTier(spec *api.ArtifactSpec) (*Tier, error) {
	afxId := spec.Name
	afxLevel := spec.Level
	familyAfxId := spec.Family()
	for _, f := range Data.ArtifactFamilies {
		if f.AfxId == familyAfxId {
			for _, t := range f.Tiers {
				if t.AfxId == afxId && t.AfxLevel == afxLevel {
					return t, nil
				}
			}
			break
		}
	}
	return nil, errors.Errorf("artifact (%s, %s) not found in data.json", afxId, afxLevel)
}
