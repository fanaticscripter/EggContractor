package main

import (
	"fmt"
	"html/template"
	"os"
	"strconv"
	"time"

	"github.com/google/go-cmp/cmp"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/fanaticscripter/EggContractor/api"
	"github.com/fanaticscripter/EggContractor/util"
)

type ShipParameters = api.ArtifactsConfigurationResponse_MissionParameters

type ship struct {
	*ShipParameters
	Sensors                    string
	LaunchesToAdvance          uint32
	TimeToAdvanceStd           time.Duration
	TimeToAdvancePro           time.Duration
	CumulativeTimeToAdvanceStd time.Duration
	CumulativeTimeToAdvancePro time.Duration
}

type fuel struct {
	Egg    api.EggType
	Amount float64
}

func main() {
	config := &api.ArtifactsConfigurationResponse{}
	err := protojson.Unmarshal([]byte(_afxConfigJSON), config)
	if err != nil {
		log.Fatalf("failed to unmarshal config snapshot: %s", err)
	}

	// Verify with remote, if possible.
	req := &api.ArtifactsConfigurationRequestPayload{
		ClientVersion: api.ClientVersion,
	}
	resp := &api.ArtifactsConfigurationResponsePayload{}
	err = api.Request("/ei_afx/config", req, resp)
	if err != nil {
		log.Errorf("API request failed: %s", err)
	} else {
		remoteConfig := resp.Config
		remoteConfig.ArtifactParameters = nil
		if diff := cmp.Diff(config, remoteConfig, protocmp.Transform()); diff != "" {
			log.Fatalf("local and remote configs have diverged: %s", diff)
		}
	}

	var ships []*ship
	var cumulativeTimeToAdvance time.Duration
	stdConcurrency := time.Duration(1)
	proConcurrency := time.Duration(3)
	for _, s := range config.MissionParameters {
		launchesToAdvance := shipRequiredLaunchesToAdvance(s.Ship)
		var timeToAdvance time.Duration
		if s.Ship == api.MissionInfo_CHICKEN_ONE {
			// Forget about the nuance of 2 tutorial missions + 2 short missions
			timeToAdvance = 3 * 20 * time.Minute
			cumulativeTimeToAdvance = timeToAdvance
		} else if s.Ship == api.MissionInfo_HENERPRISE {
			// Do not display cumulative for Henerprise
			cumulativeTimeToAdvance = 0
		} else {
			for _, t := range s.Durations {
				if t.DurationType == api.MissionInfo_SHORT {
					timeToAdvance = time.Duration(launchesToAdvance) * util.DoubleToDuration(t.Seconds)
					break
				}
			}
			if timeToAdvance == 0 {
				panic(fmt.Sprintf("short mission not found for ship %s", s.Ship))
			}
			cumulativeTimeToAdvance += timeToAdvance
		}
		ships = append(ships, &ship{
			ShipParameters:             s,
			Sensors:                    shipSensors(s.Ship),
			LaunchesToAdvance:          launchesToAdvance,
			TimeToAdvanceStd:           timeToAdvance / stdConcurrency,
			TimeToAdvancePro:           timeToAdvance / proConcurrency,
			CumulativeTimeToAdvanceStd: cumulativeTimeToAdvance / stdConcurrency,
			CumulativeTimeToAdvancePro: cumulativeTimeToAdvance / proConcurrency,
		})
	}

	tmpl := template.Must(template.New("").Funcs(template.FuncMap{
		"eggiconpath":      eggIconPath,
		"fmtduration":      util.FormatDurationWhole,
		"fuels":            missionFuels,
		"iconurl":          iconURL,
		"numfmt":           util.NumfmtWhole,
		"shipiconpath":     shipIconPath,
		"seconds2duration": util.DoubleToDuration,
	}).ParseGlob("templates/*/*.html"))
	err = os.MkdirAll("src", 0o755)
	if err != nil {
		log.Fatalf("mkdir -p src failed: %s", err)
	}
	output, err := os.Create("src/index.html")
	if err != nil {
		log.Fatalf("failed to open src/index.html for writing: %s", err)
	}
	defer output.Close()
	err = tmpl.ExecuteTemplate(output, "index.html", struct {
		Ships []*ship
	}{
		Ships: ships,
	})
	if err != nil {
		log.Fatalf("failed to render template: %s", err)
	}
}

func iconURL(relpath string, size int) string {
	dir := strconv.Itoa(size)
	if size <= 0 {
		dir = "orig"
	}
	return fmt.Sprintf("https://eggincassets.tcl.sh/%s/%s", dir, relpath)
}

func shipIconPath(ship api.MissionInfo_Spaceship) string {
	return "egginc/" + ship.IconFilename()
}

func eggIconPath(egg api.EggType) string {
	return "egginc/" + egg.IconFilename()
}

func shipSensors(ship api.MissionInfo_Spaceship) string {
	switch ship {
	case api.MissionInfo_CHICKEN_ONE:
		fallthrough
	case api.MissionInfo_CHICKEN_NINE:
		return "Basic"
	case api.MissionInfo_CHICKEN_HEAVY:
		fallthrough
	case api.MissionInfo_BCR:
		return "Intermediate"
	case api.MissionInfo_MILLENIUM_CHICKEN:
		fallthrough
	case api.MissionInfo_CORELLIHEN_CORVETTE:
		fallthrough
	case api.MissionInfo_GALEGGTICA:
		return "Advanced"
	case api.MissionInfo_CHICKFIANT:
		fallthrough
	case api.MissionInfo_VOYEGGER:
		return "Cutting Edge"
	case api.MissionInfo_HENERPRISE:
		return "Next Generation"
	}
	return ""
}

func shipRequiredLaunchesToAdvance(ship api.MissionInfo_Spaceship) uint32 {
	switch ship {
	case api.MissionInfo_CHICKEN_ONE:
		return 4
	case api.MissionInfo_CHICKEN_NINE:
		return 6
	case api.MissionInfo_CHICKEN_HEAVY:
		return 12
	case api.MissionInfo_BCR:
		return 15
	case api.MissionInfo_MILLENIUM_CHICKEN:
		return 18
	case api.MissionInfo_CORELLIHEN_CORVETTE:
		return 21
	case api.MissionInfo_GALEGGTICA:
		return 24
	case api.MissionInfo_CHICKFIANT:
		return 27
	case api.MissionInfo_VOYEGGER:
		return 30
	}
	return 0
}

func missionFuels(ship api.MissionInfo_Spaceship, durationType api.MissionInfo_DurationType) []fuel {
	return _fuels[ship][durationType]
}

// This shit is typed by hand.
var _fuels = map[api.MissionInfo_Spaceship]map[api.MissionInfo_DurationType][]fuel{
	api.MissionInfo_CHICKEN_ONE: {
		api.MissionInfo_TUTORIAL: {
			{api.EggType_ROCKET_FUEL, 1e5},
		},
		api.MissionInfo_SHORT: {
			{api.EggType_ROCKET_FUEL, 2e6},
		},
		api.MissionInfo_LONG: {
			{api.EggType_ROCKET_FUEL, 3e6},
		},
		api.MissionInfo_EPIC: {
			{api.EggType_ROCKET_FUEL, 10e6},
		},
	},
	api.MissionInfo_CHICKEN_NINE: {
		api.MissionInfo_SHORT: {
			{api.EggType_ROCKET_FUEL, 10e6},
		},
		api.MissionInfo_LONG: {
			{api.EggType_ROCKET_FUEL, 15e6},
		},
		api.MissionInfo_EPIC: {
			{api.EggType_ROCKET_FUEL, 25e6},
		},
	},
	api.MissionInfo_CHICKEN_HEAVY: {
		api.MissionInfo_SHORT: {
			{api.EggType_ROCKET_FUEL, 100e6},
		},
		api.MissionInfo_LONG: {
			{api.EggType_ROCKET_FUEL, 50e6},
			{api.EggType_FUSION, 5e6},
		},
		api.MissionInfo_EPIC: {
			{api.EggType_ROCKET_FUEL, 75e6},
			{api.EggType_FUSION, 25e6},
		},
	},
	api.MissionInfo_BCR: {
		api.MissionInfo_SHORT: {
			{api.EggType_ROCKET_FUEL, 250e6},
			{api.EggType_FUSION, 50e6},
		},
		api.MissionInfo_LONG: {
			{api.EggType_ROCKET_FUEL, 400e6},
			{api.EggType_FUSION, 75e6},
		},
		api.MissionInfo_EPIC: {
			{api.EggType_SUPERFOOD, 5e6},
			{api.EggType_ROCKET_FUEL, 300e6},
			{api.EggType_FUSION, 100e6},
		},
	},
	api.MissionInfo_MILLENIUM_CHICKEN: {
		api.MissionInfo_SHORT: {
			{api.EggType_FUSION, 5e9},
			{api.EggType_GRAVITON, 1e9},
		},
		api.MissionInfo_LONG: {
			{api.EggType_FUSION, 7e9},
			{api.EggType_GRAVITON, 5e9},
		},
		api.MissionInfo_EPIC: {
			{api.EggType_SUPERFOOD, 10e6},
			{api.EggType_FUSION, 10e9},
			{api.EggType_GRAVITON, 15e9},
		},
	},
	api.MissionInfo_CORELLIHEN_CORVETTE: {
		api.MissionInfo_SHORT: {
			{api.EggType_FUSION, 15e9},
			{api.EggType_GRAVITON, 2e9},
		},
		api.MissionInfo_LONG: {
			{api.EggType_FUSION, 20e9},
			{api.EggType_GRAVITON, 3e9},
		},
		api.MissionInfo_EPIC: {
			{api.EggType_SUPERFOOD, 500e6},
			{api.EggType_FUSION, 25e9},
			{api.EggType_GRAVITON, 5e9},
		},
	},
	api.MissionInfo_GALEGGTICA: {
		api.MissionInfo_SHORT: {
			{api.EggType_FUSION, 50e9},
			{api.EggType_GRAVITON, 10e9},
		},
		api.MissionInfo_LONG: {
			{api.EggType_FUSION, 75e9},
			{api.EggType_GRAVITON, 25e9},
		},
		api.MissionInfo_EPIC: {
			{api.EggType_FUSION, 100e9},
			{api.EggType_GRAVITON, 50e9},
			{api.EggType_ANTIMATTER, 1e9},
		},
	},
	api.MissionInfo_CHICKFIANT: {
		api.MissionInfo_SHORT: {
			{api.EggType_DILITHIUM, 200e9},
			{api.EggType_ANTIMATTER, 50e9},
		},
		api.MissionInfo_LONG: {
			{api.EggType_DILITHIUM, 250e9},
			{api.EggType_ANTIMATTER, 150e9},
		},
		api.MissionInfo_EPIC: {
			{api.EggType_TACHYON, 25e9},
			{api.EggType_DILITHIUM, 250e9},
			{api.EggType_ANTIMATTER, 250e9},
		},
	},
	api.MissionInfo_VOYEGGER: {
		api.MissionInfo_SHORT: {
			{api.EggType_DILITHIUM, 1e12},
			{api.EggType_ANTIMATTER, 1e12},
		},
		api.MissionInfo_LONG: {
			{api.EggType_DILITHIUM, 1.5e12},
			{api.EggType_ANTIMATTER, 1.5e12},
		},
		api.MissionInfo_EPIC: {
			{api.EggType_TACHYON, 100e9},
			{api.EggType_DILITHIUM, 2e12},
			{api.EggType_ANTIMATTER, 2e12},
		},
	},
	api.MissionInfo_HENERPRISE: {
		api.MissionInfo_SHORT: {
			{api.EggType_DILITHIUM, 2e12},
			{api.EggType_ANTIMATTER, 2e12},
		},
		api.MissionInfo_LONG: {
			{api.EggType_DILITHIUM, 3e12},
			{api.EggType_ANTIMATTER, 3e12},
			{api.EggType_DARK_MATTER, 3e12},
		},
		api.MissionInfo_EPIC: {
			{api.EggType_TACHYON, 1e12},
			{api.EggType_DILITHIUM, 3e12},
			{api.EggType_ANTIMATTER, 3e12},
			{api.EggType_DARK_MATTER, 3e12},
		},
	},
}

// Snapshot of JSON-encoded /ei_afx/config response (with artifactParameters
// stripped due to irrelevance), so that we can render offline.
//
// Generated by `EggContractor afx-config`.
const _afxConfigJSON = `
{
  "missionParameters": [
    {
      "ship": "CHICKEN_ONE",
      "durations": [
        {
          "durationType": "TUTORIAL",
          "seconds": 60,
          "quality": 0.9,
          "minQuality": 0,
          "maxQuality": 1.3,
          "capacity": 4
        },
        {
          "durationType": "SHORT",
          "seconds": 1200,
          "quality": 1,
          "minQuality": 0,
          "maxQuality": 1.4,
          "capacity": 4
        },
        {
          "durationType": "LONG",
          "seconds": 3600,
          "quality": 1.1,
          "minQuality": 0,
          "maxQuality": 1.7,
          "capacity": 5
        },
        {
          "durationType": "EPIC",
          "seconds": 7200,
          "quality": 1.2,
          "minQuality": 0,
          "maxQuality": 2.2,
          "capacity": 6
        }
      ],
      "capacityDEPRECATED": 0
    },
    {
      "ship": "CHICKEN_NINE",
      "durations": [
        {
          "durationType": "SHORT",
          "seconds": 1800,
          "quality": 1.4,
          "minQuality": 0,
          "maxQuality": 1.9,
          "capacity": 7
        },
        {
          "durationType": "LONG",
          "seconds": 3600,
          "quality": 1.55,
          "minQuality": 0,
          "maxQuality": 2.35,
          "capacity": 8
        },
        {
          "durationType": "EPIC",
          "seconds": 10800,
          "quality": 1.68,
          "minQuality": 0,
          "maxQuality": 2.7,
          "capacity": 9
        }
      ],
      "capacityDEPRECATED": 0
    },
    {
      "ship": "CHICKEN_HEAVY",
      "durations": [
        {
          "durationType": "SHORT",
          "seconds": 2700,
          "quality": 1.7,
          "minQuality": 0,
          "maxQuality": 2.2,
          "capacity": 12
        },
        {
          "durationType": "LONG",
          "seconds": 5400,
          "quality": 1.85,
          "minQuality": 0,
          "maxQuality": 2.7,
          "capacity": 14
        },
        {
          "durationType": "EPIC",
          "seconds": 14400,
          "quality": 1.9,
          "minQuality": 0,
          "maxQuality": 3.2,
          "capacity": 15
        }
      ],
      "capacityDEPRECATED": 0
    },
    {
      "ship": "BCR",
      "durations": [
        {
          "durationType": "SHORT",
          "seconds": 5400,
          "quality": 2,
          "minQuality": 0,
          "maxQuality": 2.7,
          "capacity": 18
        },
        {
          "durationType": "LONG",
          "seconds": 14400,
          "quality": 2.2,
          "minQuality": 0,
          "maxQuality": 3.3,
          "capacity": 20
        },
        {
          "durationType": "EPIC",
          "seconds": 28800,
          "quality": 2.4,
          "minQuality": 0,
          "maxQuality": 3.8,
          "capacity": 22
        }
      ],
      "capacityDEPRECATED": 0
    },
    {
      "ship": "MILLENIUM_CHICKEN",
      "durations": [
        {
          "durationType": "SHORT",
          "seconds": 10800,
          "quality": 3,
          "minQuality": 0.5,
          "maxQuality": 3.9,
          "capacity": 10
        },
        {
          "durationType": "LONG",
          "seconds": 21600,
          "quality": 3.3,
          "minQuality": 0.5,
          "maxQuality": 4.4,
          "capacity": 12
        },
        {
          "durationType": "EPIC",
          "seconds": 43200,
          "quality": 3.5,
          "minQuality": 0.5,
          "maxQuality": 5.2,
          "capacity": 14
        }
      ],
      "capacityDEPRECATED": 0
    },
    {
      "ship": "CORELLIHEN_CORVETTE",
      "durations": [
        {
          "durationType": "SHORT",
          "seconds": 14400,
          "quality": 3.5,
          "minQuality": 1.5,
          "maxQuality": 4.2,
          "capacity": 18
        },
        {
          "durationType": "LONG",
          "seconds": 43200,
          "quality": 3.8,
          "minQuality": 1.6,
          "maxQuality": 4.9,
          "capacity": 21
        },
        {
          "durationType": "EPIC",
          "seconds": 86400,
          "quality": 4.1,
          "minQuality": 1.75,
          "maxQuality": 5.8,
          "capacity": 24
        }
      ],
      "capacityDEPRECATED": 0
    },
    {
      "ship": "GALEGGTICA",
      "durations": [
        {
          "durationType": "SHORT",
          "seconds": 21600,
          "quality": 4,
          "minQuality": 2,
          "maxQuality": 5.2,
          "capacity": 27
        },
        {
          "durationType": "LONG",
          "seconds": 57600,
          "quality": 4.3,
          "minQuality": 2.3,
          "maxQuality": 6.1,
          "capacity": 30
        },
        {
          "durationType": "EPIC",
          "seconds": 108000,
          "quality": 4.6,
          "minQuality": 2.5,
          "maxQuality": 7.2,
          "capacity": 35
        }
      ],
      "capacityDEPRECATED": 0
    },
    {
      "ship": "CHICKFIANT",
      "durations": [
        {
          "durationType": "SHORT",
          "seconds": 28800,
          "quality": 5,
          "minQuality": 3,
          "maxQuality": 7.2,
          "capacity": 20
        },
        {
          "durationType": "LONG",
          "seconds": 86400,
          "quality": 5.6,
          "minQuality": 3.2,
          "maxQuality": 8,
          "capacity": 24
        },
        {
          "durationType": "EPIC",
          "seconds": 172800,
          "quality": 6.3,
          "minQuality": 3.4,
          "maxQuality": 9.2,
          "capacity": 28
        }
      ],
      "capacityDEPRECATED": 0
    },
    {
      "ship": "VOYEGGER",
      "durations": [
        {
          "durationType": "SHORT",
          "seconds": 43200,
          "quality": 5.8,
          "minQuality": 3.3,
          "maxQuality": 8.5,
          "capacity": 30
        },
        {
          "durationType": "LONG",
          "seconds": 129600,
          "quality": 6.4,
          "minQuality": 3.8,
          "maxQuality": 9.7,
          "capacity": 35
        },
        {
          "durationType": "EPIC",
          "seconds": 259200,
          "quality": 7.1,
          "minQuality": 3.9,
          "maxQuality": 12,
          "capacity": 40
        }
      ],
      "capacityDEPRECATED": 0
    },
    {
      "ship": "HENERPRISE",
      "durations": [
        {
          "durationType": "SHORT",
          "seconds": 86400,
          "quality": 6.6,
          "minQuality": 3.8,
          "maxQuality": 9.5,
          "capacity": 45
        },
        {
          "durationType": "LONG",
          "seconds": 172800,
          "quality": 7.3,
          "minQuality": 4.1,
          "maxQuality": 11.5,
          "capacity": 50
        },
        {
          "durationType": "EPIC",
          "seconds": 345600,
          "quality": 8,
          "minQuality": 4.5,
          "maxQuality": 14,
          "capacity": 56
        }
      ],
      "capacityDEPRECATED": 0
    }
  ],
  "artifactParameters": []
}`
