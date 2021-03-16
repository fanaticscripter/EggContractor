package solo

import (
	"github.com/fanaticscripter/EggContractor/api"
)

type shippingCapacityResearch struct {
	Id            string
	PerLevel      float64
	MaxLevels     int
	HoverOnly     bool
	HyperloopOnly bool
}

// Max shipping rates per second without researches.
var _baseShippingCapacity = map[api.VehicleType]float64{
	api.VehicleType_TRIKE:               5e3 / 60,
	api.VehicleType_TRANSIT_VAN:         15e3 / 60,
	api.VehicleType_PICKUP:              50e3 / 60,
	api.VehicleType_TEN_FOOT:            100e3 / 60,
	api.VehicleType_TWENTY_FOUR_FOOT:    250e3 / 60,
	api.VehicleType_SEMI:                500e3 / 60,
	api.VehicleType_DOUBLE_SEMI:         1e6 / 60,
	api.VehicleType_FUTURE_SEMI:         5e6 / 60,
	api.VehicleType_MEGA_SEMI:           15e6 / 60,
	api.VehicleType_HOVER_SEMI:          30e6 / 60,
	api.VehicleType_QUANTUM_TRANSPORTER: 50e6 / 60,
	api.VehicleType_HYPERLOOP_TRAIN:     50e6 / 60,
}

var _shippingCapacityCommonResearches = []shippingCapacityResearch{
	{"leafsprings", 0.05, 30, false, false},
	{"lightweight_boxes", 0.1, 40, false, false},
	{"driver_training", 0.05, 30, false, false},
	{"super_alloy", 0.05, 50, false, false},
	{"quantum_storage", 0.05, 20, false, false},
	{"hover_upgrades", 0.05, 25, true, false},
	{"dark_containment", 0.05, 25, false, false},
	{"neural_net_refine", 0.05, 25, false, false},
	{"hyper_portalling", 0.05, 25, false, true},
}

var _shippingCapacityEpicResearches = []shippingCapacityResearch{
	{"transportation_lobbyist", 0.05, 30, false, false},
}

func maxEggsShippedPerSecond(farm *api.Farm, epicResearches []*api.EpicResearch) float64 {
	commonResearches := farm.Researches
	var rate float64
	perVehicleRates := make(map[api.VehicleType]float64)
	for i, vehicle := range farm.VehicleTypes {
		length := farm.TrainLengths[i]
		perVehicleRate, ok := perVehicleRates[vehicle]
		if !ok {
			perVehicleRate = vehicleMaxEggsShippedPerSecond(vehicle, commonResearches, epicResearches)
			perVehicleRates[vehicle] = perVehicleRate
		}
		rate += perVehicleRate * float64(length)
	}
	return rate
}

func vehicleMaxEggsShippedPerSecond(vehicle api.VehicleType, commonResearches []*api.Research, epicResearches []*api.EpicResearch) float64 {
	rate := _baseShippingCapacity[vehicle]
	for _, r := range _shippingCapacityCommonResearches {
		if r.HoverOnly && vehicle < api.VehicleType_HOVER_SEMI {
			continue
		}
		if r.HyperloopOnly && vehicle != api.VehicleType_HYPERLOOP_TRAIN {
			continue
		}
		for _, rr := range commonResearches {
			if r.Id == rr.Id {
				rate *= 1 + float64(rr.Level)*r.PerLevel
				break
			}
		}
	}
	for _, r := range _shippingCapacityEpicResearches {
		if r.HoverOnly && vehicle < api.VehicleType_HOVER_SEMI {
			continue
		}
		if r.HyperloopOnly && vehicle != api.VehicleType_HYPERLOOP_TRAIN {
			continue
		}
		for _, rr := range epicResearches {
			if r.Id == rr.Id {
				rate *= 1 + float64(rr.Level)*r.PerLevel
				break
			}
		}
	}
	return rate
}
