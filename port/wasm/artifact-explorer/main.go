package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"syscall/js"

	"github.com/fanaticscripter/EggContractor/api"
	"github.com/fanaticscripter/EggContractor/util"
)

var (
	_playerIdPattern   = regexp.MustCompile(`(?i)^EI\d+$`)
	_nonIdCharsPattern = regexp.MustCompile(`[^A-Za-z0-9_-]+`)
)

type result struct {
	Successful bool        `json:"successful"`
	Data       interface{} `json:"data"`
	Err        string      `json:"error"`
}

func dataResult(data interface{}) *result {
	return &result{
		Successful: true,
		Data:       data,
	}
}

func errorResult(err error) *result {
	return &result{
		Successful: false,
		Err:        err.Error(),
	}
}

type ship struct {
	Name       string `json:"name"`
	AbbrevName string `json:"abbrevName"`
	IconPath   string `json:"iconPath"`
}

type missionParams = api.ArtifactsConfigurationResponse_MissionParameters_Duration

type mission struct {
	Id           string         `json:"id"`
	Display      string         `json:"display"`
	ShipName     string         `json:"shipName"`
	ShipIconPath string         `json:"shipIconPath"`
	Type         string         `json:"type"`
	AbbrevType   string         `json:"abbrevType"`
	MinQuality   float64        `json:"minQuality"`
	MaxQuality   float64        `json:"maxQuality"`
	Params       *missionParams `json:"params"`
}

type artifactParams = api.ArtifactsConfigurationResponse_ArtifactParameters

type artifact struct {
	SortKey       string                  `json:"sortKey"` // Unique key that is also good for sorting
	ItemId        string                  `json:"itemId"`  // The ItemId does not include rarity
	Family        *family                 `json:"family"`
	Name          string                  `json:"name"`
	Level         api.ArtifactSpec_Level  `json:"level"`
	Rarity        api.ArtifactSpec_Rarity `json:"rarity"`
	TierNumber    int                     `json:"tierNumber"`
	TierName      string                  `json:"tierName"`
	RarityDisplay string                  `json:"rarityDisplay"`
	IconPath      string                  `json:"iconPath"`
	Quality       float64                 `json:"quality"`
	Params        *artifactParams         `json:"params"`
}

type family struct {
	Id       api.ArtifactSpec_Name `json:"id"`
	Name     string                `json:"name"`
	Effect   string                `json:"effect"`
	Type     api.ArtifactSpec_Type `json:"type"`
	TypeName string                `json:"typeName"`
}

func retrieveData() *result {
	req := &api.ArtifactsConfigurationRequestPayload{
		ClientVersion: api.ClientVersion,
	}
	resp := &api.ArtifactsConfigurationResponsePayload{}
	err := api.Request("/ei_afx/config", req, resp)
	if err != nil {
		return errorResult(err)
	}
	config := resp.Config

	ships := make([]*ship, 0)
	missions := make([]*mission, 0)
	for _, s := range config.MissionParameters {
		ships = append(ships, &ship{
			Name:       s.Ship.Name(),
			AbbrevName: abbreviatedShipName(s.Ship),
			IconPath:   "egginc/" + s.Ship.IconFilename(),
		})
		for _, d := range s.Durations {
			if d.DurationType == api.MissionInfo_TUTORIAL {
				continue
			}
			shipName := s.Ship.Name()
			typ := d.DurationType.Display()
			missions = append(missions, &mission{
				Id:           missionId(s.Ship, d.DurationType),
				Display:      shipName + ", " + typ,
				ShipName:     shipName,
				ShipIconPath: "egginc/" + s.Ship.IconFilename(),
				Type:         typ,
				AbbrevType:   abbreviatedMissionType(d.DurationType),
				MinQuality:   float64(d.MinQuality),
				MaxQuality:   float64(d.MaxQuality),
				Params:       d,
			})
		}
	}

	artifacts := make([]*artifact, 0)
	for _, p := range config.ArtifactParameters {
		a := p.Spec
		artifacts = append(artifacts, &artifact{
			SortKey:       artifactSortKey(a),
			ItemId:        artifactItemId(a),
			Family:        newFamily(a.Name),
			Name:          a.CasedName(),
			Level:         a.Level,
			Rarity:        a.Rarity,
			TierNumber:    a.TierNumber(),
			TierName:      a.CasedTierName(),
			RarityDisplay: a.Rarity.Display(),
			IconPath:      "egginc/" + a.IconFilename(),
			Quality:       p.BaseQuality,
			Params:        p,
		})
	}

	// CSV exports
	var b bytes.Buffer
	w := csv.NewWriter(&b)
	_ = w.Write([]string{
		"Ship", "Type", "Duration", "Duration seconds", "Capacity", "Quality", "Min quality", "Max quality",
	})
	for _, m := range missions {
		_ = w.Write([]string{
			m.ShipName,
			m.Params.DurationType.Display(),
			util.FormatDurationWhole(util.DoubleToDuration(m.Params.Seconds)),
			fmt.Sprintf("%.0f", m.Params.Seconds),
			fmt.Sprintf("%d", m.Params.Capacity),
			fmt.Sprintf("%f", m.Params.Quality),
			fmt.Sprintf("%f", m.MinQuality),
			fmt.Sprintf("%f", m.MaxQuality),
		})
	}
	w.Flush()
	missionsCSV := b.String()

	b.Reset()
	w = csv.NewWriter(&b)
	header := []string{"Item", "Tier", "Base quality", "Odds multiplier"}
	for _, s := range config.MissionParameters {
		for _, d := range s.Durations {
			header = append(header, abbreviatedShipName(s.Ship)+" "+abbreviatedMissionType(d.DurationType))
		}
	}
	header = append(header, []string{
		"value", "crafting price base", "crafting price low", "crafting price domain", "crafting price curve",
	}...)
	_ = w.Write(header)
	for _, a := range artifacts {
		name := a.Name
		if a.Rarity > 0 {
			name += ", " + a.RarityDisplay
		}
		row := []string{
			name,
			fmt.Sprintf("%d", a.TierNumber),
			fmt.Sprintf("%f", a.Quality),
			fmt.Sprintf("%f", a.Params.OddsMultiplier),
		}
		for _, m := range missions {
			withinRange := m.MinQuality <= a.Quality && a.Quality <= m.MaxQuality
			row = append(row, fmt.Sprintf("%t", withinRange))
		}
		row = append(row, []string{
			fmt.Sprintf("%f", a.Params.Value),
			fmt.Sprintf("%f", a.Params.CraftingPrice),
			fmt.Sprintf("%f", a.Params.CraftingPriceLow),
			fmt.Sprintf("%d", a.Params.CraftingPriceDomain),
			fmt.Sprintf("%f", a.Params.CraftingPriceCurve),
		}...)
		_ = w.Write(row)
	}
	w.Flush()
	artifactsCSV := b.String()

	return dataResult(struct {
		Ships        []*ship     `json:"ships"`
		Missions     []*mission  `json:"missions"`
		Artifacts    []*artifact `json:"artifacts"`
		MissionsCSV  string      `json:"missionsCSV"`
		ArtifactsCSV string      `json:"artifactsCSV"`
	}{
		Ships:        ships,
		Missions:     missions,
		Artifacts:    artifacts,
		MissionsCSV:  missionsCSV,
		ArtifactsCSV: artifactsCSV,
	})
}

func missionId(ship api.MissionInfo_Spaceship, durationType api.MissionInfo_DurationType) string {
	return strings.ToLower(strings.ReplaceAll(ship.Name()+" "+durationType.Display(), " ", "-"))
}

func abbreviatedShipName(s api.MissionInfo_Spaceship) string {
	switch s {
	case api.MissionInfo_CHICKEN_ONE:
		return "C1"
	case api.MissionInfo_CHICKEN_NINE:
		return "C9"
	case api.MissionInfo_CHICKEN_HEAVY:
		return "CH"
	case api.MissionInfo_BCR:
		return "BCR"
	case api.MissionInfo_MILLENIUM_CHICKEN:
		return "QC"
	case api.MissionInfo_CORELLIHEN_CORVETTE:
		return "CHC"
	case api.MissionInfo_GALEGGTICA:
		return "G"
	case api.MissionInfo_CHICKFIANT:
		return "D"
	case api.MissionInfo_VOYEGGER:
		return "V"
	case api.MissionInfo_HENERPRISE:
		return "H"
	}
	return ""
}

func abbreviatedMissionType(t api.MissionInfo_DurationType) string {
	switch t {
	case api.MissionInfo_SHORT:
		return "SH"
	case api.MissionInfo_LONG:
		return "ST"
	case api.MissionInfo_EPIC:
		return "EX"
	}
	return ""
}

func artifactSortKey(a *api.ArtifactSpec) string {
	return fmt.Sprintf(
		"%3d-T%d-R%d-%s",
		getArtifactFamilyOrder(a),
		a.TierNumber(),
		a.Rarity,
		a.GameName(),
	)
}

func artifactItemId(a *api.ArtifactSpec) string {
	id := fmt.Sprintf("%s-%d", a.Family().GameName(), a.TierNumber())
	id = strings.ToLower(strings.ReplaceAll(id, " ", "-"))
	return _nonIdCharsPattern.ReplaceAllString(id, "")
}

func newFamily(id api.ArtifactSpec_Name) *family {
	f := id.Family()
	typ := f.ArtifactType()
	var typeName string
	switch typ {
	case api.ArtifactSpec_ARTIFACT:
		typeName = "artifact"
	case api.ArtifactSpec_STONE:
		typeName = "stone"
	case api.ArtifactSpec_INGREDIENT:
		typeName = "ingredient"
	}
	return &family{
		Id:       f,
		Name:     f.CasedName(),
		Effect:   _artifactFamilyEffects[f],
		Type:     typ,
		TypeName: typeName,
	}
}

func main() {
	// I can't think of any communications mechanism other than global variables
	// and callbacks. (Note that we can't set a directly global variable for the
	// result, since when we do that the global variable seems to be somehow
	// "cached" for a while when accessed immediately, so if we run two
	// instances with different input args, when accessing the result of the
	// second run we would somehow still get the result of the first run... I
	// didn't investigate further since the callback route works despite the
	// increased complexity.)
	//
	// Related:
	// https://github.com/golang/go/issues/25612
	// https://stackoverflow.com/q/56398142
	res := retrieveData()
	encoded, _ := json.Marshal(res)
	js.Global().Call("wasmCallback", js.ValueOf(string(encoded)))
}

var _artifactFamilyOrder = []api.ArtifactSpec_Name{
	api.ArtifactSpec_PUZZLE_CUBE,
	api.ArtifactSpec_LUNAR_TOTEM,
	api.ArtifactSpec_DEMETERS_NECKLACE,
	api.ArtifactSpec_VIAL_MARTIAN_DUST,
	api.ArtifactSpec_AURELIAN_BROOCH,
	api.ArtifactSpec_TUNGSTEN_ANKH,
	api.ArtifactSpec_ORNATE_GUSSET,
	api.ArtifactSpec_NEODYMIUM_MEDALLION,
	api.ArtifactSpec_MERCURYS_LENS,
	api.ArtifactSpec_BEAK_OF_MIDAS,
	api.ArtifactSpec_CARVED_RAINSTICK,
	api.ArtifactSpec_INTERSTELLAR_COMPASS,
	api.ArtifactSpec_THE_CHALICE,
	api.ArtifactSpec_PHOENIX_FEATHER,
	api.ArtifactSpec_QUANTUM_METRONOME,
	api.ArtifactSpec_DILITHIUM_MONOCLE,
	api.ArtifactSpec_TITANIUM_ACTUATOR,
	api.ArtifactSpec_SHIP_IN_A_BOTTLE,
	api.ArtifactSpec_TACHYON_DEFLECTOR,
	api.ArtifactSpec_BOOK_OF_BASAN,
	api.ArtifactSpec_LIGHT_OF_EGGENDIL,
	api.ArtifactSpec_LUNAR_STONE,
	api.ArtifactSpec_SHELL_STONE,
	api.ArtifactSpec_TACHYON_STONE,
	api.ArtifactSpec_TERRA_STONE,
	api.ArtifactSpec_SOUL_STONE,
	api.ArtifactSpec_DILITHIUM_STONE,
	api.ArtifactSpec_QUANTUM_STONE,
	api.ArtifactSpec_LIFE_STONE,
	api.ArtifactSpec_CLARITY_STONE,
	api.ArtifactSpec_PROPHECY_STONE,
	api.ArtifactSpec_GOLD_METEORITE,
	api.ArtifactSpec_TAU_CETI_GEODE,
	api.ArtifactSpec_SOLAR_TITANIUM,
}

func getArtifactFamilyOrder(a *api.ArtifactSpec) int {
	family := a.Family()
	for i, f := range _artifactFamilyOrder {
		if family == f {
			return i
		}
	}
	return len(_artifactFamilyOrder)
}

var _artifactFamilyEffects = map[api.ArtifactSpec_Name]string{
	api.ArtifactSpec_PUZZLE_CUBE:          "Lower research costs",
	api.ArtifactSpec_LUNAR_TOTEM:          "Modify away earnings",
	api.ArtifactSpec_DEMETERS_NECKLACE:    "Increase egg value",
	api.ArtifactSpec_VIAL_MARTIAN_DUST:    "Increase max running chicken bonus",
	api.ArtifactSpec_AURELIAN_BROOCH:      "Increase drone rewards",
	api.ArtifactSpec_TUNGSTEN_ANKH:        "Increases egg value",
	api.ArtifactSpec_ORNATE_GUSSET:        "Increase hen house capacity",
	api.ArtifactSpec_NEODYMIUM_MEDALLION:  "Increase drone frequency",
	api.ArtifactSpec_MERCURYS_LENS:        "increases farm value",
	api.ArtifactSpec_BEAK_OF_MIDAS:        "Increase gold reward chance",
	api.ArtifactSpec_CARVED_RAINSTICK:     "Increase chance of cash rewards from gifts and drones",
	api.ArtifactSpec_INTERSTELLAR_COMPASS: "Increase egg shipping rate",
	api.ArtifactSpec_THE_CHALICE:          "Improved internal hatcheries",
	api.ArtifactSpec_PHOENIX_FEATHER:      "Increased soul egg collection rate",
	api.ArtifactSpec_QUANTUM_METRONOME:    "Increases egg laying rate",
	api.ArtifactSpec_DILITHIUM_MONOCLE:    "increases boost effectiveness",
	api.ArtifactSpec_TITANIUM_ACTUATOR:    "increase hold to hatch rate",
	api.ArtifactSpec_SHIP_IN_A_BOTTLE:     "Increase co-op mates earnings",
	api.ArtifactSpec_TACHYON_DEFLECTOR:    "Increase co-op mates egg laying rate",
	api.ArtifactSpec_BOOK_OF_BASAN:        "Increases effect of Eggs of Prophecy",
	api.ArtifactSpec_LIGHT_OF_EGGENDIL:    "Enlightenment egg value increase",
	api.ArtifactSpec_LUNAR_STONE:          "Increases away earnings when set",
	api.ArtifactSpec_SHELL_STONE:          "Increases egg value when set",
	api.ArtifactSpec_TACHYON_STONE:        "Increases egg laying rate when set",
	api.ArtifactSpec_TERRA_STONE:          "Increases max running chicken bonus when set",
	api.ArtifactSpec_SOUL_STONE:           "Increases soul egg bonus when set",
	api.ArtifactSpec_DILITHIUM_STONE:      "Increases boost duration",
	api.ArtifactSpec_QUANTUM_STONE:        "Increases shipping capacity when set",
	api.ArtifactSpec_LIFE_STONE:           "Improves internal hatcheries when set",
	api.ArtifactSpec_CLARITY_STONE:        "Enables effect of host artifact on enlightenment egg farm.",
	api.ArtifactSpec_PROPHECY_STONE:       "Increases egg of prophecy egg bonus when set",
	api.ArtifactSpec_GOLD_METEORITE:       "",
	api.ArtifactSpec_TAU_CETI_GEODE:       "",
	api.ArtifactSpec_SOLAR_TITANIUM:       "",
}
