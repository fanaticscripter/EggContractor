package main

import (
	"encoding/base64"
	"encoding/json"
	"sort"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/fanaticscripter/EggContractor/api"
	"github.com/fanaticscripter/EggContractor/db"
	"github.com/fanaticscripter/EggContractor/util"
)

type contract struct {
	*api.ContractProperties
	IsLeggacy  bool
	HasLeggacy bool
}

func (c *contract) Type() string {
	if c.IsLeggacy {
		return "Leggacy"
	} else {
		return "Original"
	}
}

func (c *contract) EstimatedOfferingTime() time.Time {
	if c.Id == "first-contract" {
		return time.Time{}
	}
	// I believe original contracts expire in 21 days, where as Leggacy
	// contracts expire in 7 days.
	if c.IsLeggacy {
		return c.ExpiryTime().Add(-7 * 24 * time.Hour)
	} else {
		return c.ExpiryTime().Add(-21 * 24 * time.Hour)
	}
}

func (c *contract) StandardRewards() []*api.ContractProperties_Reward {
	if len(c.RewardTiers) >= 2 {
		return c.RewardTiers[1].Rewards
	}
	return c.Rewards
}

func (c *contract) EliteRewards() []*api.ContractProperties_Reward {
	if len(c.RewardTiers) >= 2 {
		return c.RewardTiers[0].Rewards
	}
	return nil
}

func (c *contract) StandardGoalsStr() []string {
	rewards := c.StandardRewards()
	goals := make([]string, 0)
	for _, r := range rewards {
		goals = append(goals, util.NumfmtWhole(r.Goal))
	}
	goals = extend(goals, 3)
	return goals
}

func (c *contract) EliteGoalsStr() []string {
	rewards := c.EliteRewards()
	goals := make([]string, 0)
	for _, r := range rewards {
		goals = append(goals, util.NumfmtWhole(r.Goal))
	}
	goals = extend(goals, 3)
	return goals
}

func (c *contract) StandardUltimateGoal() float64 {
	rewards := c.StandardRewards()
	if len(rewards) == 0 {
		return 0
	}
	return rewards[len(rewards)-1].Goal
}

func (c *contract) EliteUltimateGoal() float64 {
	rewards := c.EliteRewards()
	if len(rewards) == 0 {
		return 0
	}
	return rewards[len(rewards)-1].Goal
}

func (c *contract) ProphecyEggCount() int {
	for _, r := range c.Rewards {
		if r.Type == api.RewardType_PROPHECY_EGG {
			return int(r.Count)
		}
	}
	return 0
}

func (c *contract) StandardProphecyEggGoal() float64 {
	rewards := c.StandardRewards()
	for _, r := range rewards {
		if r.Type == api.RewardType_PROPHECY_EGG {
			return r.Goal
		}
	}
	return 0
}

func (c *contract) EliteProphecyEggGoal() float64 {
	rewards := c.EliteRewards()
	for _, r := range rewards {
		if r.Type == api.RewardType_PROPHECY_EGG {
			return r.Goal
		}
	}
	return 0
}

func (c *contract) JSON() ([]byte, error) {
	b, err := json.Marshal(c.ContractProperties)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (c *contract) B64Protobuf() ([]byte, error) {
	p, err := proto.Marshal(c.ContractProperties)
	if err != nil {
		return nil, err
	}
	enc := base64.StdEncoding
	b := make([]byte, enc.EncodedLen(len(p)))
	enc.Encode(b, p)
	return b, nil
}

func getContractsFromDB() ([]*contract, error) {
	contracts, err := db.GetContracts()
	if err != nil {
		return nil, err
	}
	sort.Slice(contracts, func(i, j int) bool {
		return contracts[i].ExpiryTimestamp < contracts[j].ExpiryTimestamp
	})
	wrappedContracts := make([]*contract, 0)
	seenIds := make(map[string]struct{})
	leggacyIds := make(map[string]struct{})
	for _, c := range contracts {
		_, isLeggacy := seenIds[c.Id]
		wrappedContracts = append(wrappedContracts, &contract{
			ContractProperties: c,
			IsLeggacy:          isLeggacy,
		})
		seenIds[c.Id] = struct{}{}
		if isLeggacy {
			leggacyIds[c.Id] = struct{}{}
		}
	}
	for _, c := range wrappedContracts {
		_, hasLeggacy := leggacyIds[c.Id]
		if hasLeggacy {
			c.HasLeggacy = true
		}
	}
	sort.Slice(wrappedContracts, func(i, j int) bool {
		return wrappedContracts[i].EstimatedOfferingTime().Before(wrappedContracts[j].EstimatedOfferingTime())
	})
	return wrappedContracts, nil
}

func extend(s []string, length int) []string {
	for len(s) < length {
		s = append(s, "")
	}
	return s
}
