package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/fanaticscripter/EggContractor/api"
	"github.com/pkg/errors"
)

const _consumptionDataFile = "data/consumption-data.json"

var _consumptionOutcomes []ConsumptionOutcome

type ConsumptionOutcome struct {
	Item               Item                `json:"item"`
	Deterministic      bool                `json:"deterministic"`
	Gold               float64             `json:"gold"`
	ExpectedByproducts []ExpectedByproduct `json:"expected_byproducts"`
	SampleByproducts   [][]Byproduct       `json:"sample_byproducts"`
}

type Item struct {
	AfxId        api.ArtifactSpec_Name   `json:"afx_id"`
	AfxLevel     api.ArtifactSpec_Level  `json:"afx_level"`
	AfxRarity    api.ArtifactSpec_Rarity `json:"afx_rarity"`
	Id           string                  `json:"id"` // Not in original, must be populated.
	Name         string                  `json:"name"`
	TierNumber   int                     `json:"tier_number"`
	Rarity       string                  `json:"rarity"`
	IconFilename string                  `json:"icon_filename"` // Not in original, must be populated.
}

type ExpectedByproduct struct {
	Item
	ExpectedCount float64 `json:"expected_count"`
}

type Byproduct struct {
	Item
	Count int `json:"count"`
}

func loadConsumptionData() error {
	body, err := ioutil.ReadFile(_consumptionDataFile)
	if err != nil {
		return errors.Wrapf(err, "error reading %s", _consumptionDataFile)
	}
	err = json.Unmarshal(body, &_consumptionOutcomes)
	if err != nil {
		return errors.Wrapf(err, "error unmarshalling %s", _consumptionDataFile)
	}
	return nil
}

func (it *Item) Complete() {
	tier := getTier(it.AfxId, it.AfxLevel)
	if tier == nil {
		panic(fmt.Sprintf("tier (%s, %s) not found", it.AfxId, it.AfxLevel))
	}
	it.Id = tier.Id
	it.IconFilename = tier.IconFilename
}

func (c *ConsumptionOutcome) Complete() {
	c.Item.Complete()
	for i, bp := range c.ExpectedByproducts {
		bp.Item.Complete()
		c.ExpectedByproducts[i] = bp
	}
	for _, sample := range c.SampleByproducts {
		for i, bp := range sample {
			bp.Item.Complete()
			sample[i] = bp
		}
	}
}
