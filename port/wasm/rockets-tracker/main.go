package main

import (
	"encoding/json"
	"syscall/js"

	"github.com/fanaticscripter/EggContractor/api"
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

func retrieveMissions(playerId string) *result {
	fc, err := api.RequestFirstContact(&api.FirstContactRequestPayload{
		EiUserId: playerId,
	})
	if err != nil {
		return errorResult(err)
	}

	artifactsDB := fc.Data.ArtifactsDb
	activeMissions := make([]*mission, 0)
	launched := make([]*api.MissionInfo, 0)
	for _, m := range artifactsDB.MissionInfos {
		activeMissions = append(activeMissions, newMission(m))
		if m.Status >= api.MissionInfo_EXPLORING {
			launched = append(launched, m)
		}
	}
	launched = append(launched, artifactsDB.MissionArchive...)
	stats, progress := generateStatsFromMissionArchive(launched)
	return dataResult(struct {
		ActiveMissions []*mission      `json:"activeMissions"`
		MissionStats   *missionStats   `json:"missionStats"`
		UnlockProgress *unlockProgress `json:"unlockProgress"`
	}{
		ActiveMissions: activeMissions,
		MissionStats:   stats,
		UnlockProgress: progress,
	})
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
	res := retrieveMissions(playerId)
	encoded, _ := json.Marshal(res)
	js.Global().Call("wasmCallback", js.ValueOf(string(encoded)))
}
