package main

import (
	"fmt"
	"log"
	"os"
	"text/template"
)

type artifactClass struct {
	Id         string
	Name       string
	Effect     string
	LevelNames []string
}

var _artifacts = []artifactClass{
	{
		Id:         "lunar_totem",
		Name:       "Lunar totem",
		Effect:     "Modify away earnings",
		LevelNames: []string{"Basic", "Regular", "Powerful", "Eggceptional"},
	},
	{
		Id:         "neodymium_medallion",
		Name:       "Neodymium medallion",
		Effect:     "Increase drone frequency",
		LevelNames: []string{"Weak", "Regular", "Precise", "Eggceptional"},
	},
	{
		Id:         "beak_of_midas",
		Name:       "Beak of midas",
		Effect:     "Increase gold reward chance",
		LevelNames: []string{"Dull", "Regular", "Jeweled", "Glistening"},
	},
	{
		Id:         "light_of_eggendil",
		Name:       "Light of eggendil",
		Effect:     "Enlightenment egg value increase",
		LevelNames: []string{"Dim", "Shimmering", "Glowing", "Brilliant"},
	},
	{
		Id:         "demeters_necklace",
		Name:       "Demeters necklace",
		Effect:     "Increase egg value",
		LevelNames: []string{"Simple", "Jeweled", "Pristine", "Beggspoke"},
	},
	{
		Id:         "vial_of_martian_dust",
		Name:       "Vial of martian dust",
		Effect:     "Increase max running chicken bonus",
		LevelNames: []string{"Tiny", "Regular", "Hermetic", "Prime"},
	},
	{
		Id:         "ornate_gusset",
		Name:       "Gusset",
		Effect:     "Increase hen house capacity",
		LevelNames: []string{"Plain", "Ornate", "Distegguished", "Jeweled"},
	},
	{
		Id:         "the_chalice",
		Name:       "The chalice",
		Effect:     "Improved internal hatcheries",
		LevelNames: []string{"Plain", "Polished", "Jeweled", "Eggceptional"},
	},
	{
		Id:         "book_of_basan",
		Name:       "Book of basan",
		Effect:     "Increases effect of Eggs of Prophecy",
		LevelNames: []string{"Regular", "Collectors", "Fortified", "Gilded"},
	},
	{
		Id:         "phoenix_feather",
		Name:       "Phoenix feather",
		Effect:     "Increased soul egg collection rate",
		LevelNames: []string{"Tattered", "Regular", "Brilliant", "Blazing"},
	},
	{
		Id:         "tungsten_ankh",
		Name:       "Tungsten ankh",
		Effect:     "Increases egg value",
		LevelNames: []string{"Crude", "Regular", "Polished", "Brilliant"},
	},
	{
		Id:         "aurelian_brooch",
		Name:       "Aurelian brooch",
		Effect:     "Increase drone rewards",
		LevelNames: []string{"Plain", "Regular", "Jeweled", "Eggceptional"},
	},
	{
		Id:         "carved_rainstick",
		Name:       "Carved rainstick",
		Effect:     "Increase chance of cash rewards from gifts and drones",
		LevelNames: []string{"Simple", "Regular", "Ornate", "Meggnificent"},
	},
	{
		Id:         "puzzle_cube",
		Name:       "Puzzle cube",
		Effect:     "Lower research costs",
		LevelNames: []string{"Ancient", "Regular", "Mystical", "Unsolvable"},
	},
	{
		Id:         "quantum_metronome",
		Name:       "Quantum metronome",
		Effect:     "Increases egg laying rate",
		LevelNames: []string{"Misaligned", "Adequate", "Perfect", "Reggference"},
	},
	{
		Id:         "ship_in_a_bottle",
		Name:       "Ship in a bottle",
		Effect:     "Increase co-op mates earnings",
		LevelNames: []string{"Regular", "Detailed", "Complex", "Eggquisite"},
	},
	{
		Id:         "tachyon_deflector",
		Name:       "Tachyon deflector",
		Effect:     "Increase co-op mates egg laying rate",
		LevelNames: []string{"Weak", "Regular", "Robust", "Eggceptional"},
	},
	{
		Id:         "interstellar_compass",
		Name:       "Interstellar compass",
		Effect:     "Increase egg shipping rate",
		LevelNames: []string{"Miscalibrated", "Regular", "Precise", "Clairvoyant"},
	},
	{
		Id:         "dilithium_monocle",
		Name:       "Dilithium monocle",
		Effect:     "increases boost effectiveness",
		LevelNames: []string{"Regular", "Precise", "Eggsacting", "Flawless"},
	},
	{
		Id:         "titanium_actuator",
		Name:       "Titanium actuator",
		Effect:     "increase hold to hatch rate",
		LevelNames: []string{"Inconsistent", "Regular", "Precise", "Reggference"},
	},
	{
		Id:         "mercurys_lens",
		Name:       "Mercury's lens",
		Effect:     "increases farm value",
		LevelNames: []string{"Misaligned", "Regular", "Precise", "Meggnificent"},
	},
}

var _stones = []artifactClass{
	{
		Id:         "tachyon_stone",
		Name:       "Tachyon stone",
		Effect:     "Increases egg laying rate when set",
		LevelNames: []string{"Fragment", "Regular", "Eggsquisite", "Brilliant"},
	},
	{
		Id:         "dilithium_stone",
		Name:       "Dilithium stone",
		Effect:     "Increases boost duration",
		LevelNames: []string{"Fragment", "Regular", "Eggsquisite", "Brilliant"},
	},
	{
		Id:         "shell_stone",
		Name:       "Shell stone",
		Effect:     "Increases egg value when set",
		LevelNames: []string{"Fragment", "Regular", "Eggsquisite", "Flawless"},
	},
	{
		Id:         "lunar_stone",
		Name:       "Lunar stone",
		Effect:     "Increases away earnings when set",
		LevelNames: []string{"Fragment", "Regular", "Eggsquisite", "Meggnificent"},
	},
	{
		Id:         "soul_stone",
		Name:       "Soul stone",
		Effect:     "Increases soul egg bonus when set",
		LevelNames: []string{"Fragment", "Regular", "Eggsquisite", "Radiant"},
	},
	{
		Id:         "prophecy_stone",
		Name:       "Prophecy stone",
		Effect:     "Increases egg of prophecy egg bonus when set",
		LevelNames: []string{"Fragment", "Regular", "Eggsquisite", "Radiant"},
	},
	{
		Id:         "quantum_stone",
		Name:       "Quantum stone",
		Effect:     "Increases shipping capacity when set",
		LevelNames: []string{"Fragment", "Regular", "Phased", "Meggnificent"},
	},
	{
		Id:         "terra_stone",
		Name:       "Terra stone",
		Effect:     "Increases max running chicken bonus when set",
		LevelNames: []string{"Fragment", "Regular", "Rich", "Eggceptional"},
	},
	{
		Id:         "life_stone",
		Name:       "Life stone",
		Effect:     "Improves internal hatcheries when set",
		LevelNames: []string{"Fragment", "Regular", "Good", "Eggceptional"},
	},
	{
		Id:         "clarity_stone",
		Name:       "Clarity stone",
		Effect:     "Enables effect of host artifact on enlightenment egg farm.",
		LevelNames: []string{"Fragment", "Regular", "Eggsquisite", "Eggceptional"},
	},
}

var _ingredients = []artifactClass{
	{
		Id:         "gold_meteorite",
		Name:       "Gold meteorite",
		LevelNames: []string{"Tiny", "Enriched", "Solid"},
	},
	{
		Id:         "tau_ceti_geode",
		Name:       "Tau ceti geode",
		LevelNames: []string{"Piece", "Glimmering", "Radiant"},
	},
	{
		Id:         "solar_titanium",
		Name:       "Solar titanium",
		LevelNames: []string{"Ore", "Bar", "Geogon"},
	},
}

var _unconfirmedIngredients = []artifactClass{
	{
		Id:         "et_aluminum",
		Name:       "Extraterrestrial aluminum",
		LevelNames: []string{"?", "?", "?"},
	},
	{
		Id:         "ancient_tungsten",
		Name:       "Ancient tungsten",
		LevelNames: []string{"?", "?", "?"},
	},
	{
		Id:         "space_rocks",
		Name:       "Space rocks",
		LevelNames: []string{"?", "?", "?"},
	},
	{
		Id:         "alien_wood",
		Name:       "Alien wood",
		LevelNames: []string{"?", "?", "?"},
	},
	{
		Id:         "centaurian_steel",
		Name:       "Centaurian steel",
		LevelNames: []string{"?", "?", "?"},
	},
	{
		Id:         "eridani_feather",
		Name:       "Eridani feather",
		LevelNames: []string{"?", "?", "?"},
	},
	// DRONE_PARTS doesn't even have an icon.
	// {
	// 	Id:         "drone_parts",
	// 	Name:       "Drone parts",
	// 	LevelNames: []string{"?", "?", "?"},
	// },
	{
		Id:         "celestial_bronze",
		Name:       "Celestial bronze",
		LevelNames: []string{"?", "?", "?"},
	},
	{
		Id:         "lalande_hide",
		Name:       "Lalande hide",
		LevelNames: []string{"?", "?", "?"},
	},
}

// idx should be 0-indexed.
func afxIconPath(artifact artifactClass, idx int) string {
	id := artifact.Id
	switch id {
	case "light_of_eggendil":
		id = "light_eggendil"
	case "neodymium_medallion":
		id = "neo_medallion"
	case "vial_of_martian_dust":
		id = "vial_martian_dust"
	}
	return fmt.Sprintf("static/afx_%s_%d.png", id, idx+1)
}

func main() {
	tmpl := template.Must(template.New("").Funcs(template.FuncMap{
		"afxiconpath": afxIconPath,
	}).ParseGlob("templates/*/*.html"))
	err := os.MkdirAll("src", 0o755)
	if err != nil {
		log.Fatalf("mkdir -p src failed: %s", err)
	}
	output, err := os.Create("src/index.html")
	if err != nil {
		log.Fatalf("failed to open src/index.html for writing: %s", err)
	}
	defer output.Close()
	err = tmpl.ExecuteTemplate(output, "index.html", struct {
		Artifacts              []artifactClass
		Stones                 []artifactClass
		Ingredients            []artifactClass
		UnconfirmedIngredients []artifactClass
	}{
		Artifacts:              _artifacts,
		Stones:                 _stones,
		Ingredients:            _ingredients,
		UnconfirmedIngredients: _unconfirmedIngredients,
	})
	if err != nil {
		log.Fatalf("failed to render template: %s", err)
	}
}
