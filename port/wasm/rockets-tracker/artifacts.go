package main

import (
	"github.com/fanaticscripter/EggContractor/api"
)

type artifactClassInfo struct {
	Id        api.ArtifactSpec_Name
	Effect    string
	TierNames []string
}

type artifactsProgress struct {
	Artifacts   []*artifactClass `json:"artifacts"`
	Stones      []*artifactClass `json:"stones"`
	Ingredients []*artifactClass `json:"ingredients"`
}

type artifactClass struct {
	Id       api.ArtifactSpec_Name `json:"id"`
	Name     string                `json:"name"`
	Effect   string                `json:"effect"`
	Unlocked bool                  `json:"unlocked"`
	Tiers    []*artifactTier       `json:"tiers"`
}

type artifactTier struct {
	Name                 string    `json:"name"`
	IconPath             string    `json:"iconPath"`
	TierNumber           int       `json:"tierNumber"`
	Unlocked             bool      `json:"unlocked"`
	PreviousTierUnlocked bool      `json:"previousTierUnlocked"`
	Count                uint32    `json:"count"`
	CraftedCount         uint32    `json:"craftedCount"`
	RarityCounts         [4]uint32 `json:"rarityCounts"`
}

var _artifactsInfo = []*artifactClassInfo{
	{
		Id:        api.ArtifactSpec_PUZZLE_CUBE,
		Effect:    "Lower research costs",
		TierNames: []string{"Ancient", "Regular", "Mystical", "Unsolvable"},
	},
	{
		Id:        api.ArtifactSpec_LUNAR_TOTEM,
		Effect:    "Modify away earnings",
		TierNames: []string{"Basic", "Regular", "Powerful", "Eggceptional"},
	},
	{
		Id:        api.ArtifactSpec_DEMETERS_NECKLACE,
		Effect:    "Increase egg value",
		TierNames: []string{"Simple", "Jeweled", "Pristine", "Beggspoke"},
	},
	{
		Id:        api.ArtifactSpec_VIAL_MARTIAN_DUST,
		Effect:    "Increase max running chicken bonus",
		TierNames: []string{"Tiny", "Regular", "Hermetic", "Prime"},
	},
	{
		Id:        api.ArtifactSpec_AURELIAN_BROOCH,
		Effect:    "Increase drone rewards",
		TierNames: []string{"Plain", "Regular", "Jeweled", "Eggceptional"},
	},
	{
		Id:        api.ArtifactSpec_TUNGSTEN_ANKH,
		Effect:    "Increases egg value",
		TierNames: []string{"Crude", "Regular", "Polished", "Brilliant"},
	},
	{
		Id:        api.ArtifactSpec_ORNATE_GUSSET,
		Effect:    "Increase hen house capacity",
		TierNames: []string{"Plain", "Ornate", "Distegguished", "Jeweled"},
	},
	{
		Id:        api.ArtifactSpec_NEODYMIUM_MEDALLION,
		Effect:    "Increase drone frequency",
		TierNames: []string{"Weak", "Regular", "Precise", "Eggceptional"},
	},
	{
		Id:        api.ArtifactSpec_MERCURYS_LENS,
		Effect:    "increases farm value",
		TierNames: []string{"Misaligned", "Regular", "Precise", "Meggnificent"},
	},
	{
		Id:        api.ArtifactSpec_BEAK_OF_MIDAS,
		Effect:    "Increase gold reward chance",
		TierNames: []string{"Dull", "Regular", "Jeweled", "Glistening"},
	},
	{
		Id:        api.ArtifactSpec_CARVED_RAINSTICK,
		Effect:    "Increase chance of cash rewards from gifts and drones",
		TierNames: []string{"Simple", "Regular", "Ornate", "Meggnificent"},
	},
	{
		Id:        api.ArtifactSpec_INTERSTELLAR_COMPASS,
		Effect:    "Increase egg shipping rate",
		TierNames: []string{"Miscalibrated", "Regular", "Precise", "Clairvoyant"},
	},
	{
		Id:        api.ArtifactSpec_THE_CHALICE,
		Effect:    "Improved internal hatcheries",
		TierNames: []string{"Plain", "Polished", "Jeweled", "Eggceptional"},
	},
	{
		Id:        api.ArtifactSpec_PHOENIX_FEATHER,
		Effect:    "Increased soul egg collection rate",
		TierNames: []string{"Tattered", "Regular", "Brilliant", "Blazing"},
	},
	{
		Id:        api.ArtifactSpec_QUANTUM_METRONOME,
		Effect:    "Increases egg laying rate",
		TierNames: []string{"Misaligned", "Adequate", "Perfect", "Reggference"},
	},
	{
		Id:        api.ArtifactSpec_DILITHIUM_MONOCLE,
		Effect:    "increases boost effectiveness",
		TierNames: []string{"Regular", "Precise", "Eggsacting", "Flawless"},
	},
	{
		Id:        api.ArtifactSpec_TITANIUM_ACTUATOR,
		Effect:    "increase hold to hatch rate",
		TierNames: []string{"Inconsistent", "Regular", "Precise", "Reggference"},
	},
	{
		Id:        api.ArtifactSpec_SHIP_IN_A_BOTTLE,
		Effect:    "Increase co-op mates earnings",
		TierNames: []string{"Regular", "Detailed", "Complex", "Eggquisite"},
	},
	{
		Id:        api.ArtifactSpec_TACHYON_DEFLECTOR,
		Effect:    "Increase co-op mates egg laying rate",
		TierNames: []string{"Weak", "Regular", "Robust", "Eggceptional"},
	},
	{
		Id:        api.ArtifactSpec_BOOK_OF_BASAN,
		Effect:    "Increases effect of Eggs of Prophecy",
		TierNames: []string{"Regular", "Collectors", "Fortified", "Gilded"},
	},
	{
		Id:        api.ArtifactSpec_LIGHT_OF_EGGENDIL,
		Effect:    "Enlightenment egg value increase",
		TierNames: []string{"Dim", "Shimmering", "Glowing", "Brilliant"},
	},
}

var _stonesInfo = []*artifactClassInfo{
	{
		Id:        api.ArtifactSpec_LUNAR_STONE,
		Effect:    "Increases away earnings when set",
		TierNames: []string{"Fragment", "Regular", "Eggsquisite", "Meggnificent"},
	},
	{
		Id:        api.ArtifactSpec_SHELL_STONE,
		Effect:    "Increases egg value when set",
		TierNames: []string{"Fragment", "Regular", "Eggsquisite", "Flawless"},
	},
	{
		Id:        api.ArtifactSpec_TACHYON_STONE,
		Effect:    "Increases egg laying rate when set",
		TierNames: []string{"Fragment", "Regular", "Eggsquisite", "Brilliant"},
	},
	{
		Id:        api.ArtifactSpec_TERRA_STONE,
		Effect:    "Increases max running chicken bonus when set",
		TierNames: []string{"Fragment", "Regular", "Rich", "Eggceptional"},
	},
	{
		Id:        api.ArtifactSpec_SOUL_STONE,
		Effect:    "Increases soul egg bonus when set",
		TierNames: []string{"Fragment", "Regular", "Eggsquisite", "Radiant"},
	},
	{
		Id:        api.ArtifactSpec_DILITHIUM_STONE,
		Effect:    "Increases boost duration",
		TierNames: []string{"Fragment", "Regular", "Eggsquisite", "Brilliant"},
	},
	{
		Id:        api.ArtifactSpec_QUANTUM_STONE,
		Effect:    "Increases shipping capacity when set",
		TierNames: []string{"Fragment", "Regular", "Phased", "Meggnificent"},
	},
	{
		Id:        api.ArtifactSpec_LIFE_STONE,
		Effect:    "Improves internal hatcheries when set",
		TierNames: []string{"Fragment", "Regular", "Good", "Eggceptional"},
	},
	{
		Id:        api.ArtifactSpec_CLARITY_STONE,
		Effect:    "Enables effect of host artifact on enlightenment egg farm.",
		TierNames: []string{"Fragment", "Regular", "Eggsquisite", "Eggceptional"},
	},
	{
		Id:        api.ArtifactSpec_PROPHECY_STONE,
		Effect:    "Increases egg of prophecy egg bonus when set",
		TierNames: []string{"Fragment", "Regular", "Eggsquisite", "Radiant"},
	},
}

var _ingredientsInfo = []*artifactClassInfo{
	{
		Id:        api.ArtifactSpec_GOLD_METEORITE,
		TierNames: []string{"Tiny", "Enriched", "Solid"},
	},
	{
		Id:        api.ArtifactSpec_TAU_CETI_GEODE,
		TierNames: []string{"Piece", "Glimmering", "Radiant"},
	},
	{
		Id:        api.ArtifactSpec_SOLAR_TITANIUM,
		TierNames: []string{"Ore", "Bar", "Geogon"},
	},
}

func artifactFromClassTier(cls *artifactClassInfo, tierNumber int) *api.ArtifactSpec {
	if cls.Id.ArtifactType() == api.ArtifactSpec_STONE {
		if tierNumber == 1 {
			return &api.ArtifactSpec{
				Name: cls.Id.CorrespondingFragment(),
			}
		} else {
			return &api.ArtifactSpec{
				Name: cls.Id,
				// T2 => INFERIOR, T3 => LESSER, T4 => NORMAL
				Level: api.ArtifactSpec_Level(tierNumber - 2),
			}
		}
	}
	return &api.ArtifactSpec{
		Name:  cls.Id,
		Level: api.ArtifactSpec_Level(tierNumber - 1),
	}
}

func classTierFromArtifact(a *api.ArtifactSpec) (id api.ArtifactSpec_Name, tierNumber int) {
	switch a.Type() {
	case api.ArtifactSpec_STONE:
		id = a.Name
		// INFERIOR => T2, LESSER => T3, NORMAL => T4
		tierNumber = int(a.Level) + 2
		return
	case api.ArtifactSpec_STONE_INGREDIENT:
		id = a.Name.CorrespondingStone()
		tierNumber = 1
		return
	}
	id = a.Name
	tierNumber = int(a.Level) + 1
	return
}

func getArtifactsProgress(db *api.ArtifactsDB) *artifactsProgress {
	var info []*artifactClassInfo
	info = append(info, _artifactsInfo...)
	info = append(info, _stonesInfo...)
	info = append(info, _ingredientsInfo...)
	data := make(map[api.ArtifactSpec_Name]*artifactClass)
	for _, cls := range info {
		var tiers []*artifactTier
		for i, name := range cls.TierNames {
			a := artifactFromClassTier(cls, i+1)
			tiers = append(tiers, &artifactTier{
				Name:       name,
				IconPath:   "static/" + a.IconFilename(),
				TierNumber: i + 1,
			})
		}
		data[cls.Id] = &artifactClass{
			Id:     cls.Id,
			Name:   cls.Id.CasedName(),
			Effect: cls.Effect,
			Tiers:  tiers,
		}
	}

	for _, a := range db.DiscoveredArtifacts {
		id, tierNumber := classTierFromArtifact(a)
		cls, ok := data[id]
		if !ok {
			panic("class " + id.String() + " not found")
		}
		cls.Unlocked = true
		cls.Tiers[tierNumber-1].Unlocked = true
	}

	for _, a := range db.CraftingCounts {
		id, tierNumber := classTierFromArtifact(a.Spec)
		cls, ok := data[id]
		if !ok {
			panic("class " + id.String() + " not found")
		}
		cls.Tiers[tierNumber-1].CraftedCount = a.Count
	}

	recordItem := func(a *api.ArtifactSpec, cnt uint32) {
		id, tierNumber := classTierFromArtifact(a)
		cls, ok := data[id]
		if !ok {
			panic("class " + id.String() + " not found")
		}
		tier := cls.Tiers[tierNumber-1]
		tier.Count += cnt
		tier.RarityCounts[int(a.Rarity)] += cnt
	}

	for _, a := range db.InventoryItems {
		count := uint32(a.Quantity)
		recordItem(a.Artifact.Spec, count)
		for _, s := range a.Artifact.Stones {
			recordItem(s, count)
		}
	}

	for _, cls := range data {
		for i, tier := range cls.Tiers {
			if i == 0 {
				tier.PreviousTierUnlocked = true
			}
			if tier.Unlocked && i+1 < len(cls.Tiers) {
				cls.Tiers[i+1].PreviousTierUnlocked = true
			}
		}
	}

	progress := &artifactsProgress{}
	for _, a := range _artifactsInfo {
		progress.Artifacts = append(progress.Artifacts, data[a.Id])
	}
	for _, a := range _stonesInfo {
		progress.Stones = append(progress.Stones, data[a.Id])
	}
	for _, a := range _ingredientsInfo {
		progress.Ingredients = append(progress.Ingredients, data[a.Id])
	}
	return progress
}
