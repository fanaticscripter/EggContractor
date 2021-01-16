package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	log "github.com/sirupsen/logrus"

	"github.com/fanaticscripter/EggContractor/api"
	"github.com/fanaticscripter/EggContractor/util"
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

func retrievePastContracts(playerId string) *result {
	resp, err := api.RequestFirstContact(&api.FirstContactRequestPayload{
		PlayerId: playerId,
		X3:       1,
	})
	if err != nil {
		return errorResult(err)
	}
	if resp.Data == nil || resp.Data.Contracts == nil {
		return errorResult(fmt.Errorf("invalid API response: %+v", resp))
	}

	type contract struct {
		Id                      string `json:"id"`
		Name                    string `json:"name"`
		Date                    string `json:"date"`
		Code                    string `json:"code"`
		GoalsInfo               string `json:"goalsInfo"`
		Incomplete              bool   `json:"incomplete"`
		HasProphecyEgg          bool   `json:"hasProphecyEgg"`
		ProphecyEggInfo         string `json:"prophecyEggInfo"`
		ProphecyEggNotCollected bool   `json:"prophecyEggNotCollected"`
	}

	contracts := make([]contract, 0)
	for _, c := range resp.Data.Contracts.PastContracts {
		numGoalsCompleted := int(c.NumGoalsCompleted)
		totalGoals := len(c.Props.Rewards)
		var rewards []*api.ContractProperties_Reward

		var contractType string
		if len(c.Props.RewardTiers) == 0 {
			// Legacy contract without the standard/elite tier division.
			rewards = c.Props.Rewards
		} else if c.NumGoalsCompleted == 0 {
			// Can't tell standard or elite when none of the goals were completed.
			contractType = "elt"
			rewards = c.Props.RewardTiers[0].Rewards
		} else {
			eliteRewards := c.Props.RewardTiers[0].Rewards
			standardRewards := c.Props.RewardTiers[1].Rewards
			eliteCompleted := eliteRewards[c.NumGoalsCompleted-1].Goal
			standardCompleted := standardRewards[c.NumGoalsCompleted-1].Goal
			if util.NumfmtWhole(c.CompletedGoal) == util.NumfmtWhole(eliteCompleted) {
				contractType = "elt"
				rewards = eliteRewards
			} else if util.NumfmtWhole(c.CompletedGoal) == util.NumfmtWhole(standardCompleted) {
				contractType = "std"
				rewards = standardRewards
			} else {
				log.Errorf("%s: completed goal %s is neither standard nor elite\n",
					c.Props.Id, util.NumfmtWhole(c.CompletedGoal))
			}
		}

		goalsInfo := fmt.Sprintf("%d/%d", c.NumGoalsCompleted, totalGoals)
		incomplete := numGoalsCompleted < totalGoals

		var prophecyEggIndex int
		for i, r := range rewards {
			if r.Type == api.RewardType_PROPHECY_EGG {
				prophecyEggIndex = i + 1
				break
			}
		}
		hasProphecyEgg := prophecyEggIndex > 0
		var prophecyEggInfo string
		var prophecyEggNotCollected bool
		if prophecyEggIndex > 0 {
			prophecyEggInfo = fmt.Sprintf("%s #%d", contractType, prophecyEggIndex)
			prophecyEggNotCollected = numGoalsCompleted < prophecyEggIndex
		}

		contracts = append(contracts, contract{
			Id:                      c.Props.Id,
			Name:                    c.Props.Name,
			Date:                    util.FormatDate(c.StartedTime()),
			Code:                    c.Code,
			GoalsInfo:               goalsInfo,
			Incomplete:              incomplete,
			HasProphecyEgg:          hasProphecyEgg,
			ProphecyEggInfo:         prophecyEggInfo,
			ProphecyEggNotCollected: prophecyEggNotCollected,
		})
	}
	return dataResult(contracts)
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
	args := js.Global().Get("wasmArgs")
	playerId := args.Get("0").String()
	res := retrievePastContracts(playerId)
	encoded, _ := json.Marshal(res)
	js.Global().Call("wasmCallback", js.ValueOf(string(encoded)))
}
