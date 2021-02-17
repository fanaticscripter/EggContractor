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
	*Tier
	SortKey   string                  `json:"sortKey"`
	AfxRarity api.ArtifactSpec_Rarity `json:"afxRarity"`
	Rarity    string                  `json:"rarity"`
	Quality   float64                 `json:"quality"`
	IconPath  string                  `json:"iconPath"`
	Params    *artifactParams         `json:"params"`
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
		resp := &api.ArtifactsConfigurationResponsePayload{}
		err := api.RequestWithContext(ctx, "/ei_afx/config", req, resp)
		if err != nil {
			errs <- err
			cancel()
			return
		}
		config = resp.Config
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
		a, err := newArtifact(p)
		if err != nil {
			return nil, err
		}
		artifacts = append(artifacts, a)
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
