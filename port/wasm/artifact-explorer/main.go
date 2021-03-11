package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"

	"github.com/fanaticscripter/EggContractor/api"
	"github.com/fanaticscripter/EggContractor/util"
)

const _appDataFile = "src/app-data.json"

type payload struct {
	Ships        []*ship     `json:"ships"`
	Missions     []*mission  `json:"missions"`
	Artifacts    []*artifact `json:"artifacts"`
	LootTable    lootTable   `json:"lootTable"`
	MissionsCSV  string      `json:"missionsCSV"`
	ArtifactsCSV string      `json:"artifactsCSV"`
}

type ship struct {
	Name       string `json:"name"`
	AbbrevName string `json:"abbrevName"`
	IconPath   string `json:"iconPath"`
}

type missionParams = api.ArtifactsConfigurationResponse_MissionParameters_Duration

type mission struct {
	Id              string                       `json:"id"`
	Display         string                       `json:"display"`
	ShipId          api.MissionInfo_Spaceship    `json:"shipId"`
	ShipName        string                       `json:"shipName"`
	ShipIconPath    string                       `json:"shipIconPath"`
	TypeId          api.MissionInfo_DurationType `json:"typeId"`
	Type            string                       `json:"type"`
	AbbrevType      string                       `json:"abbrevType"`
	Capacity        int                          `json:"capacity"`
	DurationSeconds float64                      `json:"durationSeconds"`
	DurationDisplay string                       `json:"durationDisplay"`
	MinQuality      float64                      `json:"minQuality"`
	MaxQuality      float64                      `json:"maxQuality"`
	Params          *missionParams               `json:"params"`
}

type artifactParams = api.ArtifactsConfigurationResponse_ArtifactParameters

type artifact struct {
	*Tier
	SortKey   string                  `json:"sortKey"`
	AfxRarity api.ArtifactSpec_Rarity `json:"afxRarity"`
	Rarity    string                  `json:"rarity"`
	Quality   float64                 `json:"quality"`
	IconPath  string                  `json:"iconPath"`
	Odds      *odds                   `json:"odds"`
	Params    *artifactParams         `json:"params"`
}

type odds struct {
	Total    float64                             `json:"total"`
	Rarities map[api.ArtifactSpec_Rarity]float64 `json:"rarities"`
}

// map key is mission.Id
type lootTable map[string]*missionLootTable

type missionLootTable struct {
	MissionCount int `json:"missionCount"`
	// map key is artifact.Id
	Items map[string]*ItemCount `json:"items"`
}

func assemblePayload() (*payload, error) {
	ctx, cancel := context.WithCancel(context.Background())
	errs := make(chan error, 2)
	var wg sync.WaitGroup
	var config *api.ArtifactsConfigurationResponse

	wg.Add(1)
	go func() {
		defer wg.Done()
		req := &api.ArtifactsConfigurationRequestPayload{
			ClientVersion: api.ClientVersion,
		}
		config = &api.ArtifactsConfigurationResponse{}
		err := api.RequestAuthenticatedWithContext(ctx, "/ei_afx/config", req, config)
		if err != nil {
			errs <- err
			cancel()
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := loadEiAfxData()
		if err != nil {
			errs <- err
			cancel()
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := loadLootData()
		if err != nil {
			errs <- err
			cancel()
			return
		}
	}()

	wg.Wait()

	select {
	case err := <-errs:
		return nil, err
	default:
		// No error
	}

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
				Id:              missionId(s.Ship, d.DurationType),
				Display:         shipName + ", " + typ,
				ShipId:          s.Ship,
				ShipName:        shipName,
				ShipIconPath:    "egginc/" + s.Ship.IconFilename(),
				TypeId:          d.DurationType,
				Type:            typ,
				AbbrevType:      abbreviatedMissionType(d.DurationType),
				Capacity:        int(d.Capacity),
				DurationSeconds: d.Seconds,
				DurationDisplay: util.FormatDurationWhole(util.DoubleToDuration(d.Seconds)),
				MinQuality:      float64(d.MinQuality),
				MaxQuality:      float64(d.MaxQuality),
				Params:          d,
			})
		}
	}

	artifacts := make([]*artifact, 0)
	id2odds := make(map[string]*odds)
	for _, p := range config.ArtifactParameters {
		a, err := newArtifact(p)
		if err != nil {
			return nil, err
		}
		artifacts = append(artifacts, a)
		if id2odds[a.Id] == nil {
			id2odds[a.Id] = &odds{Rarities: make(map[api.ArtifactSpec_Rarity]float64)}
		}
		id2odds[a.Id].Total += p.OddsMultiplier
		id2odds[a.Id].Rarities[a.AfxRarity] = p.OddsMultiplier
	}
	for _, a := range artifacts {
		a.Odds = id2odds[a.Id]
	}

	loot := make(lootTable)
	for _, m := range missions {
		missionId := m.Id
		data := _lootData.MissionLoot(m.ShipId, m.TypeId)
		if data.TotalArtifactsCount%m.Capacity != 0 {
			log.Fatalf("%s loot data: invalid total artifacts count: %d not divisible by %d",
				missionId, data.TotalArtifactsCount, m.Capacity)
		}
		missionCount := data.TotalArtifactsCount / m.Capacity
		items := make(map[string]*ItemCount)
		for _, a := range artifacts {
			if a.Quality < m.MinQuality || a.Quality > m.MaxQuality {
				continue
			}
			possibleAfxRarities := a.PossibleAfxRarities
			if possibleAfxRarities == nil {
				possibleAfxRarities = []api.ArtifactSpec_Rarity{api.ArtifactSpec_COMMON}
			}
			counts, err := data.ItemCount(a.AfxId, a.AfxLevel, possibleAfxRarities)
			if err != nil {
				log.Fatalf("%s: %s", missionId, err)
			}
			items[a.Id] = counts
		}
		loot[missionId] = &missionLootTable{
			MissionCount: missionCount,
			Items:        items,
		}
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
		if a.AfxRarity > 0 {
			name += ", " + a.Rarity
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

	return &payload{
		Ships:        ships,
		Missions:     missions,
		Artifacts:    artifacts,
		LootTable:    loot,
		MissionsCSV:  missionsCSV,
		ArtifactsCSV: artifactsCSV,
	}, nil
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

func newArtifact(p *artifactParams) (*artifact, error) {
	a := p.Spec
	afxId := a.Name
	afxLevel := a.Level
	familyAfxId := a.Family()
	var tier *Tier
	for _, f := range _eiafxData.ArtifactFamilies {
		if f.AfxId == familyAfxId {
			for _, t := range f.Tiers {
				if t.AfxId == afxId && t.AfxLevel == afxLevel {
					tier = t
					break
				}
			}
			break
		}
	}
	if tier == nil {
		return nil, fmt.Errorf("artifact (%s, %s) not found in data.json", afxId, afxLevel)
	}
	return &artifact{
		Tier:      tier,
		SortKey:   fmt.Sprintf("%3d-T%d-R%d-%s", tier.Family.SortKey, tier.TierNumber, a.Rarity, tier.Name),
		AfxRarity: a.Rarity,
		Rarity:    a.Rarity.Display(),
		Quality:   p.BaseQuality,
		IconPath:  "egginc/" + tier.IconFilename,
		Params:    p,
	}, nil
}

func main() {
	data, err := assemblePayload()
	if err != nil {
		log.Fatal(err)
	}
	encoded, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("error serializing app payload: %s", err)
	}
	encoded = append(encoded, '\n')
	err = ioutil.WriteFile(_appDataFile, encoded, 0o644)
	if err != nil {
		log.Fatalf("error writing to %s: %s", _appDataFile, err)
	}
}
