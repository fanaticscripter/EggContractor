package web

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/fanaticscripter/EggContractor/api"
	"github.com/fanaticscripter/EggContractor/db"
	"github.com/fanaticscripter/EggContractor/util"
)

// Interface that accomodates api.CoopStatus, coop.CoopStatus, and web.CoopStatus.
type getMembers interface {
	GetMembers() []*api.CoopStatus_Member
}

type coopMemberPayload struct {
	Id                        string  `json:"id"`
	Name                      string  `json:"name"`
	EggsLaid                  float64 `json:"eggsLaid"`
	EggsLaidStr               string  `json:"eggsLaidStr"`
	EggsPerHour               float64 `json:"eggsPerHour"`
	EggsPerHourStr            string  `json:"eggsPerHourStr"`
	EarningBonusPercentage    float64 `json:"earningBonusPercentage"`
	EarningBonusPercentageStr string  `json:"earningBonusPercentageStr"`
	Tokens                    int32   `json:"tokens"`
	IsActive                  bool    `json:"isActive"`
	OfflineSeconds            float64 `json:"offlineSeconds"`
	OfflineTimeStr            string  `json:"offlineTimeStr"`
}

type peekerPayload struct {
	Contracts          []*api.ContractProperties
	PresetContractId   string
	PresetContractName string
	PresetCode         string
}

func newCoopMemberPayload(m *api.CoopStatus_Member) *coopMemberPayload {
	return &coopMemberPayload{
		Id:                        m.Id,
		Name:                      m.Name,
		EggsLaid:                  m.EggsLaid,
		EggsLaidStr:               util.Numfmt(m.EggsLaid),
		EggsPerHour:               m.EggsPerHour(),
		EggsPerHourStr:            util.Numfmt(m.EggsPerHour()),
		EarningBonusPercentage:    m.EarningBonusPercentage(),
		EarningBonusPercentageStr: util.Numfmt(m.EarningBonusPercentage()),
		Tokens:                    m.Tokens,
		IsActive:                  m.Active,
	}
}

func getMemberPayloads(c getMembers) []*coopMemberPayload {
	payloads := make([]*coopMemberPayload, len(c.GetMembers()))
	for i, m := range c.GetMembers() {
		payloads[i] = newCoopMemberPayload(m)
	}
	cc, ok := c.(*CoopStatus)
	if ok && cc.Activities != nil {
		for i, m := range c.GetMembers() {
			activity, ok := cc.Activities[m.Id]
			if ok {
				payloads[i].OfflineSeconds = activity.OfflineTime.Seconds()
				offlineTimeStr := util.FormatDurationHM(activity.OfflineTime)
				if activity.NoActivityRecorded {
					offlineTimeStr = "\u2265 " + offlineTimeStr
				}
				payloads[i].OfflineTimeStr = offlineTimeStr
			}
		}
	}
	return payloads
}

// Tests whether the coop has associated activity stats.
func hasActivityStats(c getMembers) bool {
	cc, ok := c.(*CoopStatus)
	if !ok {
		return false
	}
	return cc.Activities != nil
}

// On error, returns a payload with the passed in presets and no contracts.
func newPeekerPayload(presetContractId string, presetCode string) (*peekerPayload, error) {
	contracts, err := db.GetCoopContracts()
	if err != nil {
		return &peekerPayload{
			PresetContractId: presetContractId,
			PresetCode:       presetCode,
		}, errors.Wrap(err, "retrieve contract list from database")
	}
	presetContractName := ""
	if presetContractId != "" {
		for _, c := range contracts {
			if c.Id == presetContractId {
				presetContractName = c.Name
				break
			}
		}
	}
	return &peekerPayload{
		Contracts:          contracts,
		PresetContractId:   presetContractId,
		PresetContractName: presetContractName,
		PresetCode:         presetCode,
	}, nil
}

// Unlike newPeekerPayload, allows a list of contract IDs to be passed in, and
// the preset will be the first that's actually a coop contract registered in
// the database.
func newPeekerPayloadFromPresetList(presetContractIds []string) (*peekerPayload, error) {
	contracts, err := db.GetCoopContracts()
	if err != nil {
		err = errors.Wrap(err, "retrieve contract list from database")
		if len(presetContractIds) > 0 {
			return &peekerPayload{
				PresetContractId: presetContractIds[0],
			}, err
		}
		return &peekerPayload{}, err
	}
	presetContractId := ""
	presetContractName := ""
LoopPresetList:
	for _, id := range presetContractIds {
		for _, c := range contracts {
			if c.Id == id {
				presetContractId = id
				presetContractName = c.Name
				break LoopPresetList
			}
		}
	}
	return &peekerPayload{
		Contracts:          contracts,
		PresetContractId:   presetContractId,
		PresetContractName: presetContractName,
	}, nil
}

func marshalJSON(v interface{}) (string, error) {
	marshalled, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(marshalled), nil
}
