package main

import (
	"fmt"
	"math"

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
	CraftingCost         uint64    `json:"craftingCost"`
	RarityCounts         [4]uint32 `json:"rarityCounts"`
}

type artifactClassTier struct {
	Id         api.ArtifactSpec_Name
	TierNumber int
}

type craftingCostParams struct {
	Base   float64
	Low    float64
	Domain int
	Curve  float64
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

// Generated with `EggContractor afx-config`.
var _craftingCostInfo = map[artifactClassTier]craftingCostParams{
	artifactClassTier{api.ArtifactSpec_LUNAR_TOTEM, 2}:          craftingCostParams{47.913866, 4.791387, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_LUNAR_TOTEM, 3}:          craftingCostParams{4749.233211, 474.923321, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_LUNAR_TOTEM, 4}:          craftingCostParams{28986.362702, 2898.636270, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_NEODYMIUM_MEDALLION, 2}:  craftingCostParams{3982.334085, 398.233409, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_NEODYMIUM_MEDALLION, 3}:  craftingCostParams{16466.275749, 1646.627575, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_NEODYMIUM_MEDALLION, 4}:  craftingCostParams{58753.389780, 5875.338978, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_BEAK_OF_MIDAS, 2}:        craftingCostParams{6075.334222, 607.533422, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_BEAK_OF_MIDAS, 3}:        craftingCostParams{23885.787138, 2388.578714, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_BEAK_OF_MIDAS, 4}:        craftingCostParams{86083.096922, 8608.309692, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_LIGHT_OF_EGGENDIL, 2}:    craftingCostParams{67892.952668, 6789.295267, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_LIGHT_OF_EGGENDIL, 3}:    craftingCostParams{168121.517425, 16812.151743, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_LIGHT_OF_EGGENDIL, 4}:    craftingCostParams{330068.881399, 33006.888140, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_DEMETERS_NECKLACE, 2}:    craftingCostParams{383.059174, 38.305917, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_DEMETERS_NECKLACE, 3}:    craftingCostParams{7628.255210, 762.825521, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_DEMETERS_NECKLACE, 4}:    craftingCostParams{41267.360175, 4126.736017, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_VIAL_MARTIAN_DUST, 2}:    craftingCostParams{3302.979512, 330.297951, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_VIAL_MARTIAN_DUST, 3}:    craftingCostParams{26353.836098, 2635.383610, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_VIAL_MARTIAN_DUST, 4}:    craftingCostParams{139253.590590, 13925.359059, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_ORNATE_GUSSET, 2}:        craftingCostParams{5167.170022, 516.717002, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_ORNATE_GUSSET, 3}:        craftingCostParams{28986.362702, 2898.636270, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_ORNATE_GUSSET, 4}:        craftingCostParams{163774.023981, 16377.402398, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_THE_CHALICE, 2}:          craftingCostParams{8798.389852, 879.838985, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_THE_CHALICE, 3}:          craftingCostParams{30365.963161, 3036.596316, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_THE_CHALICE, 4}:          craftingCostParams{139253.590590, 13925.359059, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_BOOK_OF_BASAN, 2}:        craftingCostParams{114405.161123, 11440.516112, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_BOOK_OF_BASAN, 3}:        craftingCostParams{360472.914183, 36047.291418, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_BOOK_OF_BASAN, 4}:        craftingCostParams{934701.796280, 93470.179628, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_PHOENIX_FEATHER, 2}:      craftingCostParams{16466.275749, 1646.627575, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_PHOENIX_FEATHER, 3}:      craftingCostParams{50473.355229, 5047.335523, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_PHOENIX_FEATHER, 4}:      craftingCostParams{237296.526000, 23729.652600, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_TUNGSTEN_ANKH, 2}:        craftingCostParams{3982.334085, 398.233409, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_TUNGSTEN_ANKH, 3}:        craftingCostParams{25099.583520, 2509.958352, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_TUNGSTEN_ANKH, 4}:        craftingCostParams{110546.051566, 11054.605157, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_AURELIAN_BROOCH, 2}:      craftingCostParams{1186.987877, 118.698788, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_AURELIAN_BROOCH, 3}:      craftingCostParams{13827.006446, 1382.700645, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_AURELIAN_BROOCH, 4}:      craftingCostParams{58753.389780, 5875.338978, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_CARVED_RAINSTICK, 2}:     craftingCostParams{10082.597699, 1008.259770, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_CARVED_RAINSTICK, 3}:     craftingCostParams{39572.068752, 3957.206875, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_CARVED_RAINSTICK, 4}:     craftingCostParams{168121.517425, 16812.151743, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_PUZZLE_CUBE, 2}:          craftingCostParams{90.580782, 9.058078, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_PUZZLE_CUBE, 3}:          craftingCostParams{14672.717776, 1467.271778, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_PUZZLE_CUBE, 4}:          craftingCostParams{110353.422196, 11035.342220, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_QUANTUM_METRONOME, 2}:    craftingCostParams{18400.435579, 1840.043558, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_QUANTUM_METRONOME, 3}:    craftingCostParams{67892.952668, 6789.295267, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_QUANTUM_METRONOME, 4}:    craftingCostParams{306621.058816, 30662.105882, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_SHIP_IN_A_BOTTLE, 2}:     craftingCostParams{36322.087764, 3632.208776, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_SHIP_IN_A_BOTTLE, 3}:     craftingCostParams{113895.404464, 11389.540446, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_SHIP_IN_A_BOTTLE, 4}:     craftingCostParams{442510.337308, 44251.033731, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_TACHYON_DEFLECTOR, 2}:    craftingCostParams{77412.173390, 7741.217339, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_TACHYON_DEFLECTOR, 3}:    craftingCostParams{259516.135935, 25951.613593, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_TACHYON_DEFLECTOR, 4}:    craftingCostParams{1220591.960920, 122059.196092, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_INTERSTELLAR_COMPASS, 2}: craftingCostParams{9425.903017, 942.590302, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_INTERSTELLAR_COMPASS, 3}: craftingCostParams{37923.865687, 3792.386569, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_INTERSTELLAR_COMPASS, 4}: craftingCostParams{226570.172702, 22657.017270, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_DILITHIUM_MONOCLE, 2}:    craftingCostParams{24487.669756, 2448.766976, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_DILITHIUM_MONOCLE, 3}:    craftingCostParams{80590.781591, 8059.078159, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_DILITHIUM_MONOCLE, 4}:    craftingCostParams{271970.137871, 27197.013787, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_TITANIUM_ACTUATOR, 2}:    craftingCostParams{26353.836098, 2635.383610, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_TITANIUM_ACTUATOR, 3}:    craftingCostParams{88920.446829, 8892.044683, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_TITANIUM_ACTUATOR, 4}:    craftingCostParams{215900.725714, 21590.072571, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_MERCURYS_LENS, 2}:        craftingCostParams{7084.307300, 708.430730, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_MERCURYS_LENS, 3}:        craftingCostParams{30365.963161, 3036.596316, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_MERCURYS_LENS, 4}:        craftingCostParams{295510.569713, 29551.056971, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_TACHYON_STONE, 2}:        craftingCostParams{495.936563, 49.593656, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_TACHYON_STONE, 3}:        craftingCostParams{18048.119505, 1804.811950, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_TACHYON_STONE, 4}:        craftingCostParams{236946.431795, 23694.643179, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_DILITHIUM_STONE, 2}:      craftingCostParams{19671.158436, 1967.115844, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_DILITHIUM_STONE, 3}:      craftingCostParams{82572.111552, 8257.211155, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_DILITHIUM_STONE, 4}:      craftingCostParams{459362.657726, 45936.265773, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_SHELL_STONE, 2}:          craftingCostParams{262.014460, 26.201446, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_SHELL_STONE, 3}:          craftingCostParams{18048.119505, 1804.811950, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_SHELL_STONE, 4}:          craftingCostParams{183125.193310, 18312.519331, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_LUNAR_STONE, 2}:          craftingCostParams{406.733739, 40.673374, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_LUNAR_STONE, 3}:          craftingCostParams{12447.661756, 1244.766176, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_LUNAR_STONE, 4}:          craftingCostParams{155125.282426, 15512.528243, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_SOUL_STONE, 2}:           craftingCostParams{5196.192180, 519.619218, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_SOUL_STONE, 3}:           craftingCostParams{64770.664593, 6477.066459, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_SOUL_STONE, 4}:           craftingCostParams{250022.006603, 25002.200660, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_PROPHECY_STONE, 2}:       craftingCostParams{50259.443610, 5025.944361, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_PROPHECY_STONE, 3}:       craftingCostParams{246536.331347, 24653.633135, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_PROPHECY_STONE, 4}:       craftingCostParams{731674.386774, 73167.438677, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_QUANTUM_STONE, 2}:        craftingCostParams{7252.777969, 725.277797, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_QUANTUM_STONE, 3}:        craftingCostParams{25118.095517, 2511.809552, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_QUANTUM_STONE, 4}:        craftingCostParams{228278.681628, 22827.868163, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_TERRA_STONE, 2}:          craftingCostParams{5196.192180, 519.619218, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_TERRA_STONE, 3}:          craftingCostParams{29988.962680, 2998.896268, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_TERRA_STONE, 4}:          craftingCostParams{303837.601247, 30383.760125, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_LIFE_STONE, 2}:           craftingCostParams{10507.145320, 1050.714532, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_LIFE_STONE, 3}:           craftingCostParams{63177.402252, 6317.740225, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_LIFE_STONE, 4}:           craftingCostParams{325027.911927, 32502.791193, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_CLARITY_STONE, 2}:        craftingCostParams{36603.479778, 3660.347978, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_CLARITY_STONE, 3}:        craftingCostParams{185026.382927, 18502.638293, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_CLARITY_STONE, 4}:        craftingCostParams{774986.051383, 77498.605138, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_GOLD_METEORITE, 2}:       craftingCostParams{3982.530000, 398.253000, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_GOLD_METEORITE, 3}:       craftingCostParams{31373.351250, 3137.335125, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_TAU_CETI_GEODE, 2}:       craftingCostParams{7867.320000, 786.732000, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_TAU_CETI_GEODE, 3}:       craftingCostParams{53149.830000, 5314.983000, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_SOLAR_TITANIUM, 2}:       craftingCostParams{3753.000000, 375.300000, 300, 0.2},
	artifactClassTier{api.ArtifactSpec_SOLAR_TITANIUM, 3}:       craftingCostParams{37794.360000, 3779.436000, 300, 0.2},
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
				IconPath:   "egginc/" + a.IconFilename(),
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

		params, ok := _craftingCostInfo[artifactClassTier{id, tierNumber}]
		if !ok {
			panic(fmt.Sprintf("(%s, %d) not found in _craftingCostInfo", id, tierNumber))
		}
		cls.Tiers[tierNumber-1].CraftingCost = craftingCost(&params, a.Count)
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

func craftingCost(params *craftingCostParams, count uint32) uint64 {
	var total uint64
	for i := uint32(0); i < count; i++ {
		total += singleCraftCost(params, i)
	}
	return total
}

func singleCraftCost(params *craftingCostParams, previousCrafts uint32) uint64 {
	cost := math.Max(1, params.Base-(params.Base-params.Low)*
		math.Pow(math.Min(1, float64(previousCrafts)/float64(params.Domain)), params.Curve))
	return uint64(cost)
}
