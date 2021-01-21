package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

// Strings extracted with string(1) from lib/armeabi-v7a/libegginc.so from Egg
// Inc v1.20.0 apk (downloaded from apkpure.com).
//
// I had to manually insert a line "CHALICES" after "the_chalice" since somehow
// the pattern breaks here.
//
// Also added stone fragments, "DILITHIUM MONOCLE" (regular) and "INTERSTELLAR
// COMPASS".
const _strings = `lunar_totem
LUNAR TOTEMS
Modify away earnings
BASIC LUNAR TOTEM
LUNAR TOTEM
POWERFUL LUNAR TOTEM
EGGCEPTIONAL LUNAR TOTEM
neodymium_medallion
NEODYMIUM MEDALLIONS
Increase drone frequency
WEAK NEODYMIUM MEDALLION
NEODYMIUM MEDALLION
PRECISE NEODYMIUM MEDALLION
EGGCEPTIONAL NEODYMIUM MEDALLION
beak_of_midas
BEAKS OF MIDAS
Increase gold reward chance
DULL BEAK OF MIDAS
BEAK OF MIDAS
JEWELED BEAK OF MIDAS
GLISTENING BEAK OF MIDAS
light_of_eggendil
LIGHT OF EGGENDILS
Enlightenment egg value increase
DIM LIGHT OF EGGENDIL
SHIMMERING LIGHT OF EGGENDIL
GLOWING LIGHT OF EGGENDIL
BRILLIANT LIGHT OF EGGENDIL
demeters_necklace
DEMETERS NECKLACES
Increase egg value
SIMPLE DEMETERS NECKLACE
JEWELED DEMETERS NECKLACE
PRISTINE DEMETERS NECKLACE
BEGGSPOKE DEMETERS NECKLACE
vial_of_martian_dust
VIALS OF MARTIAN DUST
Increase max running chicken bonus
TINY VIAL OF MARTIAN DUST
VIAL OF MARTIAN DUST
HERMETIC VIAL OF MARTIAN DUST
PRIME VIAL OF MARTIAN DUST
ornate_gusset
GUSSETS
Increase hen house capacity
PLAIN GUSSET
ORNATE GUSSET
DISTEGGUISHED GUSSET
JEWELED GUSSET
the_chalice
CHALICES
Improved internal hatcheries
PLAIN CHALICE
POLISHED CHALICE
JEWELED CHALICE
EGGCEPTIONAL CHALICE
book_of_basan
BOOKS OF BASAN
Increases effect of Eggs of Prophecy
BOOK OF BASAN
COLLECTORS BOOK OF BASAN
FORTIFIED BOOK OF BASAN
GILDED BOOK OF BASAN
phoenix_feather
PHOENIX FEATHERS
Increased soul egg collection rate
TATTERED PHOENIX FEATHER
PHOENIX FEATHER
BRILLIANT PHOENIX FEATHER
BLAZING PHOENIX FEATHER
tungsten_ankh
TUNGSTEN ANKHS
Increases egg value
CRUDE TUNGSTEN ANKH
TUNGSTEN ANKH
POLISHED TUNGSTEN ANKH
BRILLIANT TUNGSTEN ANKH
aurelian_brooch
AURELIAN BROOCHES
Increase drone rewards
PLAIN AURELIAN BROOCH
AURELIAN BROOCH
JEWELED AURELIAN BROOCH
EGGCEPTIONAL AURELIAN BROOCH
carved_rainstick
CARVED RAINSTICKS
Increase chance of cash rewards from gifts and drones
SIMPLE CARVED RAINSTICK
CARVED RAINSTICK
ORNATE CARVED RAINSTICK
MEGGNIFICENT CARVED RAINSTICK
puzzle_cube
PUZZLE CUBES
Lower research costs
ANCIENT PUZZLE CUBE
PUZZLE CUBE
MYSTICAL PUZZLE CUBE
UNSOLVABLE PUZZLE CUBE
quantum_metronome
QUANTUM METRONOMES
Increases egg laying rate
MISALIGNED QUANTUM METRONOME
ADEQUATE QUANTUM METRONOME
PERFECT QUANTUM METRONOME
REGGFERENCE QUANTUM METRONOME
ship_in_a_bottle
SHIPS IN BOTTLES
Increase co-op mates earings
SHIP IN A BOTTLE
DETAILED SHIP IN A BOTTLE
COMPLEX SHIP IN A BOTTLE
EGGQUISITE SHIP IN A BOTTLE
tachyon_deflector
TACHYON DEFLECTORS
Increase co-op mates egg laying rate
WEAK TACHYON DEFLECTOR
TACHYON DEFLECTOR
ROBUST TACHYON DEFLECTOR
EGGCEPTIONAL TACHYON DEFLECTOR
interstellar_compass
INTERSTELLAR COMPASS
Increase egg shipping rate
MISCALIBRATED INTERSTELLAR COMPASS
INTERSTELLAR COMPASS
PRECISE INTERSTELLAR COMPASS
CLAIRVOYANT INTERSTELLAR COMPASS
dilithium_monocle
DILITHIUM MONOCLE
increases boost effectiveness
DILITHIUM MONOCLE
PRECISE DILITHIUM MONOCLE
EGGSACTING DILITHIUM MONOCLE
FLAWLESS DILITHIUM MONOCLE
titanium_actuator
TITANIUM ACTUATORS
increase hold to hatch rate
INCONSISTENT TITANIUM ACTUATOR
TITANIUM ACTUATOR
PRECISE TITANIUM ACTUATOR
REGGFERENCE TITANIUM ACTUATOR
mercurys_lens
MERCURY'S LENSES
increases farm value
MISALIGNED MERCURY'S LENS
MERCURY'S LENS
PRECISE MERCURY'S LENS
MEGGNIFICENT MERCURY'S LENS
tachyon_stone
TACHYON STONES
Increases egg laying rate when set
FRAGMENT TACHYON STONE
TACHYON STONE
EGGSQUISITE TACHYON STONE
BRILLIANT TACHYON STONE
dilithium_stone
DILITHIUM STONES
Increases boost duration
FRAGMENT DILITHIUM STONE
DILITHIUM STONE
EGGSQUISITE DILITHIUM STONE
BRILLIANT DILITHIUM STONE
shell_stone
SHELL STONES
Increases egg value when set
FRAGMENT SHELL STONE
SHELL STONE
EGGSQUISITE SHELL STONE
FLAWLESS SHELL STONE
lunar_stone
LUNAR STONES
Increases away earnings when set
FRAGMENT LUNAR STONE
LUNAR STONE
EGGSQUISITE LUNAR STONE
MEGGNIFICENT LUNAR STONE
soul_stone
SOUL STONES
Increases soul egg bonus when set
FRAGMENT SOUL STONE
SOUL STONE
EGGSQUISITE SOUL STONE
RADIANT SOUL STONE
prophecy_stone
PROPHECY STONES
Increases egg of prophecy egg bonus when set
FRAGMENT PROPHECY STONE
PROPHECY STONE
EGGSQUISITE PROPHECY STONE
RADIANT PROPHECY STONE
quantum_stone
QUANTUM STONES
Increases shipping capacity when set
FRAGMENT QUANTUM STONE
QUANTUM STONE
PHASED QUANTUM STONE
MEGGNIFICENT QUANTUM STONE
terra_stone
TERRA STONES
Increases max running chicken bonus when set
FRAGMENT TERRA STONE
TERRA STONE
RICH TERRA STONE
EGGCEPTIONAL TERRA STONE
life_stone
LIFE STONES
Improves internal hatcheries when set
FRAGMENT LIFE STONE
LIFE STONE
GOOD LIFE STONE
EGGCEPTIONAL LIFE STONE
clarity_stone
CLARITY STONES
Enables effect of host artifact on enlightenment egg farm.
FRAGMENT CLARITY STONE
CLARITY STONE
EGGSQUISITE CLARITY STONE
EGGCEPTIONAL CLARITY STONE
`

type artifactClass struct {
	Id         string
	Name       string
	Effect     string
	LevelNames []string
}

func generateArtifactClasses() []artifactClass {
	type scannerState int
	const (
		scanId scannerState = iota
		scanPlural
		scanEffect
		scanLevels
	)
	scanner := bufio.NewScanner(strings.NewReader(_strings))
	var state scannerState
	var id, name, nameCommonSuffix, effect string
	var levelNames []string
	classes := make([]artifactClass, 0)
	for scanner.Scan() {
		line := scanner.Text()
		switch state {
		case scanLevels:
			if nameCommonSuffix != "" && strings.Contains(line, nameCommonSuffix) {
				level := strings.TrimRight(strings.TrimSuffix(line, nameCommonSuffix), " ")
				level = strings.Title(strings.ToLower(level))
				if level == "" {
					level = "Regular"
				}
				levelNames = append(levelNames, level)
				break
			}
			// New artifact class encountered.
			classes = append(classes, artifactClass{
				Id:         id,
				Name:       name,
				Effect:     effect,
				LevelNames: levelNames,
			})
			state = scanId
			fallthrough
		case scanId:
			id = line
			levelNames = make([]string, 0)
			nameLower := strings.ToLower(strings.Replace(id, "_", " ", -1))
			nameUpper := strings.ToUpper(nameLower)
			name = nameUpper[:1] + nameLower[1:]
			nameCommonSuffix = nameUpper
			// Special cases
			switch id {
			case "ornate_gusset":
				name = "Gusset"
				nameCommonSuffix = "GUSSET"
			case "the_chalice":
				name = "The chalice"
				nameCommonSuffix = "CHALICE"
			case "mercurys_lens":
				name = "Mercury's lens"
				nameCommonSuffix = "MERCURY'S LENS"
			}
			state++
		case scanPlural:
			state++
		case scanEffect:
			effect = line
			state++
		}
	}
	classes = append(classes, artifactClass{
		Id:         id,
		Name:       name,
		Effect:     effect,
		LevelNames: levelNames,
	})
	return classes
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
	artifacts := generateArtifactClasses()

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
	err = tmpl.ExecuteTemplate(output, "index.html", artifacts)
	if err != nil {
		log.Fatalf("failed to render template: %s", err)
	}
}
