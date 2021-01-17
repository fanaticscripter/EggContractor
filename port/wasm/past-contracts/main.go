package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"sync"
	"syscall/js"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context/ctxhttp"
	"google.golang.org/protobuf/proto"

	"github.com/fanaticscripter/EggContractor/api"
	"github.com/fanaticscripter/EggContractor/util"
)

var _client *http.Client

type contract struct {
	*api.ContractProperties
	IsLeggacy             bool
	EstimatedOfferingTime time.Time
}

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

func retrieveContractList(playerId string) *result {
	ctx, cancel := context.WithCancel(context.Background())
	errs := make(chan error, 2)
	var wg sync.WaitGroup
	var fc *api.FirstContact
	var allContracts []*contract

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		fc, err = api.RequestFirstContactWithContext(ctx,
			&api.FirstContactRequestPayload{
				PlayerId: playerId,
				X3:       1,
			})
		if err != nil {
			errs <- err
			cancel()
			return
		}
		if fc.Data == nil || fc.Data.Contracts == nil {
			errs <- errors.Errorf("invalid API response: %+v", fc)
			cancel()
			return
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		allContracts, err = retrieveAllContracts(ctx)
		if err != nil {
			errs <- err
			cancel()
		}
	}()
	wg.Wait()

	select {
	case err := <-errs:
		return errorResult(err)
	default:
		// No error
	}

	type contractSummary struct {
		Id                      string `json:"id"`
		Name                    string `json:"name"`
		Date                    string `json:"date"`
		Attempted               bool   `json:"attempted"`
		Code                    string `json:"code"`
		GoalsInfo               string `json:"goalsInfo"`
		Incomplete              bool   `json:"incomplete"`
		ProphecyEggInfo         string `json:"prophecyEggInfo"`
		ProphecyEggCount        int    `json:"prophecyEggCount"`
		ProphecyEggNotCollected bool   `json:"prophecyEggNotCollected"`
	}

	contracts := make([]contractSummary, 0)
	seenIds := make(map[string]struct{})
	for _, c := range fc.Data.Contracts.PastContracts {
		seenIds[c.Props.Id] = struct{}{}

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
		var prophecyEggCount int
		for i, r := range rewards {
			if r.Type == api.RewardType_PROPHECY_EGG {
				prophecyEggIndex = i + 1
				prophecyEggCount = int(math.Round(r.Count))
				break
			}
		}
		var prophecyEggInfo string
		var prophecyEggNotCollected bool
		if prophecyEggIndex > 0 {
			prophecyEggInfo = fmt.Sprintf("%s #%d", contractType, prophecyEggIndex)
			if prophecyEggCount > 1 {
				prophecyEggInfo += fmt.Sprintf(" (%d)", prophecyEggCount)
			}
			prophecyEggNotCollected = numGoalsCompleted < prophecyEggIndex
		}

		contracts = append(contracts, contractSummary{
			Id:                      c.Props.Id,
			Name:                    c.Props.Name,
			Date:                    util.FormatDate(c.StartedTime()),
			Attempted:               true,
			Code:                    c.Code,
			GoalsInfo:               goalsInfo,
			Incomplete:              incomplete,
			ProphecyEggInfo:         prophecyEggInfo,
			ProphecyEggCount:        prophecyEggCount,
			ProphecyEggNotCollected: prophecyEggNotCollected,
		})
	}

	unattemptedContracts := make([]*contract, 0)
	// Loop through contract archive in reverse (in terms of offering date) so
	// that we only record the last incarnation of each contract.
	for i := len(allContracts) - 1; i >= 0; i-- {
		c := allContracts[i]
		_, exists := seenIds[c.Id]
		if exists {
			continue
		}
		unattemptedContracts = append(unattemptedContracts, c)
		seenIds[c.Id] = struct{}{}
	}
	// Loop reverse again so that earlier contracts come first.
	for i := len(unattemptedContracts) - 1; i >= 0; i-- {
		c := unattemptedContracts[i]

		var contractType string
		var rewards []*api.ContractProperties_Reward
		if len(c.RewardTiers) == 0 {
			rewards = c.Rewards
		} else {
			contractType = "elt"
			rewards = c.RewardTiers[0].Rewards
		}
		goalsInfo := fmt.Sprintf("-/%d", len(rewards))

		var prophecyEggIndex int
		var prophecyEggCount int
		for i, r := range rewards {
			if r.Type == api.RewardType_PROPHECY_EGG {
				prophecyEggIndex = i + 1
				prophecyEggCount = int(math.Round(r.Count))
				break
			}
		}
		var prophecyEggInfo string
		if prophecyEggIndex > 0 {
			prophecyEggInfo = fmt.Sprintf("%s #%d", contractType, prophecyEggIndex)
			if prophecyEggCount > 1 {
				prophecyEggInfo += fmt.Sprintf(" (%d)", prophecyEggCount)
			}
		}

		date := util.FormatDate(c.EstimatedOfferingTime)
		if c.Id == "first-contract" {
			date = "-"
		}

		contracts = append(contracts, contractSummary{
			Id:                      c.Id,
			Name:                    c.Name,
			Date:                    date,
			Attempted:               false,
			GoalsInfo:               goalsInfo,
			Incomplete:              true,
			ProphecyEggInfo:         prophecyEggInfo,
			ProphecyEggCount:        prophecyEggCount,
			ProphecyEggNotCollected: prophecyEggCount > 0,
		})
	}

	// Prepare CSV export.
	var b bytes.Buffer
	w := csv.NewWriter(&b)
	_ = w.Write([]string{
		"ID", "Name", "Date", "Code", "Goals", "PE", "Attempted", "Completed", "PE Uncollected",
	})
	for _, c := range contracts {
		_ = w.Write([]string{
			c.Id,
			c.Name,
			c.Date,
			c.Code,
			c.GoalsInfo,
			c.ProphecyEggInfo,
			fmt.Sprintf("%t", c.Attempted),
			fmt.Sprintf("%t", !c.Incomplete),
			fmt.Sprintf("%t", c.ProphecyEggNotCollected),
		})
	}
	w.Flush()

	return dataResult(struct {
		Contracts []contractSummary `json:"contracts"`
		CSV       string            `json:"csv"`
	}{
		Contracts: contracts,
		CSV:       b.String(),
	})
}

// Retrieve a list of all historical contracts from contracts.csv.
func retrieveAllContracts(ctx context.Context) ([]*contract, error) {
	resp, err := ctxhttp.Get(ctx, _client, "contracts.csv")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("GET contracts.csv: HTTP %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	r := csv.NewReader(resp.Body)
	var labels map[string]int // map column labels to column indices
	var typeColIdx int
	contracts := make([]*contract, 0)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return contracts, errors.Wrap(err, "error reading contracts.csv")
		}
		if labels == nil {
			labels = make(map[string]int)
			for i, label := range record {
				labels[label] = i
			}
			var ok bool
			typeColIdx, ok = labels["Type"]
			if !ok {
				return contracts, errors.Wrap(err, "contracts.csv: Type column not found")
			}
			continue
		}
		id := record[0]
		typ := record[typeColIdx]
		isLeggacy := false
		if typ == "Leggacy" {
			isLeggacy = true
		}
		b64proto := record[len(record)-1]
		c, err := decodeB64Protobuf(b64proto)
		if err != nil {
			return contracts, errors.Wrapf(err,
				"contracts.csv: error decoding protobuf for contract %#v (%s):", id, b64proto)
		}
		var estimatedOfferingTime time.Time
		if isLeggacy {
			estimatedOfferingTime = c.ExpiryTime().Add(-7 * 24 * time.Hour)
		} else {
			estimatedOfferingTime = c.ExpiryTime().Add(-21 * 24 * time.Hour)
		}
		contracts = append(contracts, &contract{
			ContractProperties:    c,
			IsLeggacy:             isLeggacy,
			EstimatedOfferingTime: estimatedOfferingTime,
		})
	}
	return contracts, nil
}

func decodeB64Protobuf(s string) (*api.ContractProperties, error) {
	protob, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	contract := &api.ContractProperties{}
	err = proto.Unmarshal(protob, contract)
	if err != nil {
		return nil, err
	}
	return contract, nil
}

func init() {
	_client = &http.Client{
		Timeout: 5 * time.Second,
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
	args := js.Global().Get("wasmArgs")
	playerId := args.Get("0").String()
	res := retrieveContractList(playerId)
	encoded, _ := json.Marshal(res)
	js.Global().Call("wasmCallback", js.ValueOf(string(encoded)))
}
