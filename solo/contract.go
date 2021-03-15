package solo

import (
	"time"

	"github.com/fanaticscripter/EggContractor/api"
	"github.com/fanaticscripter/EggContractor/contract"
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
	GetEggType() api.EggType
	GetUltimateGoal() float64
	GetRewards() []*api.Reward
	GetEggsLaid() float64
	GetEggsPerSecond() float64
	GetDurationUntilProductionDeadline() time.Duration
	GetDurationUntilCollectionDeadline() time.Duration
	GetServerRefreshTime() time.Time
}

type SoloContract struct {
	BaseSoloContract
}

func (c *SoloContract) GetEggsPerHour() float64 {
	return c.GetEggsPerSecond() * 3600
}

// GetOfflineAdjustedEggsLaid returns server-reported EggsLaid plus expected
// amount laid during offline time (duration between last server refresh and
// specified client refresh), which is capped at 30hr.
func (c *SoloContract) GetOfflineAdjustedEggsLaid(clientRefreshTime time.Time) float64 {
	offlineHours := clientRefreshTime.Sub(c.GetServerRefreshTime()).Hours()
	if offlineHours > 30 {
		offlineHours = 30
	}
	return c.GetEggsLaid() + c.GetEggsPerHour()*offlineHours
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

// GetOfflineAdjustedExpectedDurationUntilFinish returns expected duration until
// finish minus the offline duration (duration between last server refresh and
// specified client refresh), capped at 30hr.
func (c *SoloContract) GetOfflineAdjustedExpectedDurationUntilFinish(clientRefreshTime time.Time) time.Duration {
	offline := clientRefreshTime.Sub(c.GetServerRefreshTime())
	if offline > 30*time.Hour {
		offline = 30 * time.Hour
	}
	expected := c.ExpectedDurationUntilFinish()
	if expected == util.InfDuration {
		return util.InfDuration
	}
	if expected <= offline {
		return 0
	}
	return expected - offline
}

func (c *SoloContract) ToPBSoloContract() *pb.SoloContract {
	return &pb.SoloContract{
		Id:                             c.GetId(),
		Name:                           c.GetName(),
		IsElite:                        c.GetIsElite(),
		EggType:                        c.GetEggType(),
		UltimateGoal:                   c.GetUltimateGoal(),
		Rewards:                        c.GetRewards(),
		EggsLaid:                       c.GetEggsLaid(),
		EggsPerSecond:                  c.GetEggsPerSecond(),
		SecondsUntilProductionDeadline: c.GetDurationUntilProductionDeadline().Seconds(),
		SecondsUntilCollectionDeadline: c.GetDurationUntilCollectionDeadline().Seconds(),
		ServerRefreshTimestamp:         util.TimeToDouble(c.GetServerRefreshTime()),
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

func (c *SoloContract) ProgressInfo() *contract.ProgressInfo {
	return c.ProgressInfoWithProjection(0)
}

func (c *SoloContract) ProgressInfoWithProjection(projectedEggsLaid float64) *contract.ProgressInfo {
	return contract.NewProgressInfo(c.GetRewards(), c.GetEggsLaid(), projectedEggsLaid)
}

type soloContract struct {
	Player   *Player
	Contract *api.Contract
	Farm     *api.Farm
}

func GetActiveSoloContracts(backup *api.FirstContact_Payload) []*SoloContract {
	player := GetPlayer(backup)
	contracts := backup.Contracts.ActiveContracts

	activeCoopContractIds := make(map[string]struct{})
	for _, c := range backup.Contracts.ActiveCoopStatuses {
		activeCoopContractIds[c.ContractId] = struct{}{}
	}

	solos := make([]*SoloContract, 0)
	for _, farm := range backup.Farms {
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

func (c *soloContract) GetEggType() api.EggType {
	return c.Contract.Props.EggType
}

func (c *soloContract) GetUltimateGoal() float64 {
	return c.Contract.Props.UltimateGoal(c.GetIsElite())
}

func (c *soloContract) GetRewards() []*api.Reward {
	if c.GetIsElite() {
		return c.Contract.Props.EliteRewards()
	} else {
		return c.Contract.Props.StandardRewards()
	}
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

func (c *soloContract) GetServerRefreshTime() time.Time {
	return c.Farm.LastSavedTime()
}
