package solo

import (
	"time"

	"github.com/fanaticscripter/EggContractor/api"
	"github.com/fanaticscripter/EggContractor/solo/pb"
	"github.com/fanaticscripter/EggContractor/util"
)

// Implementation note: unlike a coop status which is largely self-contained in
// API response (in both /first_contact and /coop_status) and only optionally
// requires a contract properties object -- which are also self-contained
// themselves -- to fill in additional context, hence allowing easy database
// storage, the status of a solo contract has to be painstakingly pieced
// together from all over the /first_contact response, as made clear by this
// package. We cannot serialize and store in the entire /first_contact response
// in the database since it's usually over 10KB in size. Therefore, to
// reconstruct the status of a solo contract later, we have to store a frozen
// struct with computed properties in the database.
//
// To reduce duplicate code, an underlying interface, BaseSoloContract, is
// shared between the implementation derived from /first_contact, and the
// implementation based on frozen and serialized data. I chose protobuf for
// serialization instead of say, SQLite JSON1, simply because we're already
// storing other structures (contract properties, coop statuses) as protobuf in
// the database.

type BaseSoloContract interface {
	GetId() string
	GetName() string
	GetIsElite() bool
	GetUltimateGoal() float64
	GetEggsLaid() float64
	GetEggsPerSecond() float64
	GetDurationUntilProductionDeadline() time.Duration
	GetDurationUntilCollectionDeadline() time.Duration
	GetLastRefreshedTime() time.Time
}

type SoloContract struct {
	BaseSoloContract
}

func (c *SoloContract) GetEggsPerHour() float64 {
	return c.GetEggsPerSecond() * 3600
}

func (c *SoloContract) RequiredEggsPerHour() float64 {
	eggsToLay := c.GetUltimateGoal() - c.GetEggsLaid()
	hoursLeft := c.GetDurationUntilProductionDeadline().Hours()
	if eggsToLay <= 0 || hoursLeft <= 0 {
		return 0
	} else {
		return eggsToLay / hoursLeft
	}
}

func (c *SoloContract) ExpectedDurationUntilFinish() time.Duration {
	eggsToLay := c.GetUltimateGoal() - c.GetEggsLaid()
	if eggsToLay <= 0 {
		return 0
	} else if c.GetEggsPerSecond() <= 0 {
		return util.InfDuration // Forever
	} else {
		return util.DoubleToDuration(eggsToLay / c.GetEggsPerSecond())
	}
}

func (c *SoloContract) ToPBSoloContract() *pb.SoloContract {
	return &pb.SoloContract{
		Id:                             c.GetId(),
		Name:                           c.GetName(),
		IsElite:                        c.GetIsElite(),
		UltimateGoal:                   c.GetUltimateGoal(),
		EggsLaid:                       c.GetEggsLaid(),
		EggsPerSecond:                  c.GetEggsPerSecond(),
		SecondsUntilProductionDeadline: c.GetDurationUntilProductionDeadline().Seconds(),
		SecondsUntilCollectionDeadline: c.GetDurationUntilCollectionDeadline().Seconds(),
		LastRefreshedTimestamp:         util.TimeToDouble(c.GetLastRefreshedTime()),
	}
}

func (c *SoloContract) HasCompleted() bool {
	return c.GetEggsLaid() >= c.GetUltimateGoal()
}

func (c *SoloContract) HasNoTimeLeft() bool {
	return c.GetDurationUntilProductionDeadline() <= 0
}

func (c *SoloContract) IsOnTrackToFinish() bool {
	return c.GetEggsPerHour() >= c.RequiredEggsPerHour()
}

type soloContract struct {
	Player   *Player
	Contract *api.Contract
	Farm     *api.Farm
}

func GetActiveSoloContracts(fc *api.FirstContact) []*SoloContract {
	player := GetPlayer(fc)
	contracts := fc.Data.Contracts.ActiveContracts

	activeCoopContractIds := make(map[string]struct{})
	for _, c := range fc.Data.Contracts.ActiveCoopStatuses {
		activeCoopContractIds[c.ContractId] = struct{}{}
	}

	solos := make([]*SoloContract, 0)
	for _, farm := range fc.Data.Farms {
		if farm.ContractId == "" {
			continue
		}
		_, isActiveCoop := activeCoopContractIds[farm.ContractId]
		if isActiveCoop {
			continue
		}
		var contract *api.Contract
		for _, c := range contracts {
			if farm.ContractId == c.Props.Id {
				contract = c
			}
		}
		if contract == nil {
			continue
		}
		solos = append(solos, &SoloContract{
			&soloContract{
				Player:   player,
				Contract: contract,
				Farm:     farm,
			},
		})
	}
	return solos
}

func (c *soloContract) GetId() string {
	return c.Contract.Props.Id
}

func (c *soloContract) GetName() string {
	return c.Contract.Props.Name
}

func (c *soloContract) GetIsElite() bool {
	return c.Player.EarningBonus() >= api.EliteEarningBonusThreshold
}

func (c *soloContract) GetUltimateGoal() float64 {
	return c.Contract.Props.UltimateGoal(c.GetIsElite())
}

func (c *soloContract) GetEggsLaid() float64 {
	return c.Farm.EggsLaid
}

func (c *soloContract) GetEggsPerSecond() float64 {
	return float64(c.Farm.ChickenCount) *
		eggsPerChickenPerSecond(c.Farm.Researches, c.Player.Progress.EpicResearches)
}

func (c *soloContract) GetDurationUntilProductionDeadline() time.Duration {
	return time.Until(c.Contract.ProductionDeadlineTime())
}

func (c *soloContract) GetDurationUntilCollectionDeadline() time.Duration {
	return time.Until(c.Contract.CollectionDeadlineTime())
}

func (c *soloContract) GetLastRefreshedTime() time.Time {
	return c.Farm.LastSavedTime()
}
