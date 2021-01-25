package api

import "fmt"

const _iconUnknown = "icon_help.png"

func (a *ArtifactSpec) TierNumber() int {
	switch a.Type() {
	case ArtifactSpec_ARTIFACT:
		// 0, 1, 2, 3 => T1, T2, T3, T4
		return int(a.Level) + 1
	case ArtifactSpec_STONE:
		// 0, 1, 2 => T2, T3, T4 (fragment as T1)
		return int(a.Level) + 2
	case ArtifactSpec_STONE_INGREDIENT:
		return 1
	case ArtifactSpec_INGREDIENT:
		// 0, 1, 2 => T1, T2, T3
		return int(a.Level) + 1
	}
	return 1
}

func (a *ArtifactSpec) IconFilename() string {
	var base string
	switch a.Name {
	// Artifacts
	case ArtifactSpec_LUNAR_TOTEM:
		base = "lunar_totem"
	case ArtifactSpec_NEODYMIUM_MEDALLION:
		base = "neo_medallion"
	case ArtifactSpec_BEAK_OF_MIDAS:
		base = "beak_of_midas"
	case ArtifactSpec_LIGHT_OF_EGGENDIL:
		base = "light_eggendil"
	case ArtifactSpec_DEMETERS_NECKLACE:
		base = "demeters_necklace"
	case ArtifactSpec_VIAL_MARTIAN_DUST:
		base = "vial_martian_dust"
	case ArtifactSpec_ORNATE_GUSSET:
		base = "ornate_gusset"
	case ArtifactSpec_THE_CHALICE:
		base = "the_chalice"
	case ArtifactSpec_BOOK_OF_BASAN:
		base = "book_of_basan"
	case ArtifactSpec_PHOENIX_FEATHER:
		base = "phoenix_feather"
	case ArtifactSpec_TUNGSTEN_ANKH:
		base = "tungsten_ankh"
	case ArtifactSpec_AURELIAN_BROOCH:
		base = "aurelian_brooch"
	case ArtifactSpec_CARVED_RAINSTICK:
		base = "carved_rainstick"
	case ArtifactSpec_PUZZLE_CUBE:
		base = "puzzle_cube"
	case ArtifactSpec_QUANTUM_METRONOME:
		base = "quantum_metronome"
	case ArtifactSpec_SHIP_IN_A_BOTTLE:
		base = "ship_in_a_bottle"
	case ArtifactSpec_TACHYON_DEFLECTOR:
		base = "tachyon_deflector"
	case ArtifactSpec_INTERSTELLAR_COMPASS:
		base = "interstellar_compass"
	case ArtifactSpec_DILITHIUM_MONOCLE:
		base = "dilithium_monocle"
	case ArtifactSpec_TITANIUM_ACTUATOR:
		base = "titanium_actuator"
	case ArtifactSpec_MERCURYS_LENS:
		base = "mercurys_lens"
	// Stones & fragments
	case ArtifactSpec_TACHYON_STONE:
		fallthrough
	case ArtifactSpec_TACHYON_STONE_FRAGMENT:
		base = "tachyon_stone"
	case ArtifactSpec_DILITHIUM_STONE:
		fallthrough
	case ArtifactSpec_DILITHIUM_STONE_FRAGMENT:
		base = "dilithium_stone"
	case ArtifactSpec_SHELL_STONE:
		fallthrough
	case ArtifactSpec_SHELL_STONE_FRAGMENT:
		base = "shell_stone"
	case ArtifactSpec_LUNAR_STONE:
		fallthrough
	case ArtifactSpec_LUNAR_STONE_FRAGMENT:
		base = "lunar_stone"
	case ArtifactSpec_SOUL_STONE:
		fallthrough
	case ArtifactSpec_SOUL_STONE_FRAGMENT:
		base = "soul_stone"
	case ArtifactSpec_PROPHECY_STONE:
		fallthrough
	case ArtifactSpec_PROPHECY_STONE_FRAGMENT:
		base = "prophecy_stone"
	case ArtifactSpec_QUANTUM_STONE:
		fallthrough
	case ArtifactSpec_QUANTUM_STONE_FRAGMENT:
		base = "quantum_stone"
	case ArtifactSpec_TERRA_STONE:
		fallthrough
	case ArtifactSpec_TERRA_STONE_FRAGMENT:
		base = "terra_stone"
	case ArtifactSpec_LIFE_STONE:
		fallthrough
	case ArtifactSpec_LIFE_STONE_FRAGMENT:
		base = "life_stone"
	case ArtifactSpec_CLARITY_STONE:
		fallthrough
	case ArtifactSpec_CLARITY_STONE_FRAGMENT:
		base = "clarity_stone"
	// Ingredients
	case ArtifactSpec_GOLD_METEORITE:
		base = "gold_meteorite"
	case ArtifactSpec_TAU_CETI_GEODE:
		base = "tau_ceti_geode"
	case ArtifactSpec_SOLAR_TITANIUM:
		base = "solar_titanium"
	// Unconfirmed ingredients
	case ArtifactSpec_EXTRATERRESTRIAL_ALUMINUM:
		base = "et_aluminum"
	case ArtifactSpec_ANCIENT_TUNGSTEN:
		base = "ancient_tungsten"
	case ArtifactSpec_SPACE_ROCKS:
		base = "space_rocks"
	case ArtifactSpec_ALIEN_WOOD:
		base = "alien_wood"
	case ArtifactSpec_CENTAURIAN_STEEL:
		base = "centaurian_steel"
	case ArtifactSpec_ERIDANI_FEATHER:
		base = "eridani_feather"
	case ArtifactSpec_DRONE_PARTS:
		return _iconUnknown
	case ArtifactSpec_CELESTIAL_BRONZE:
		base = "celestial_bronze"
	case ArtifactSpec_LALANDE_HIDE:
		base = "lalande_hide"
	}
	return fmt.Sprintf("afx_%s_%d.png", base, a.TierNumber())
}
