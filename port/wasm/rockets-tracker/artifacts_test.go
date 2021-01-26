package main

import (
	"fmt"
	"testing"
)

func TestSingleCraftCost(t *testing.T) {
	tests := []struct {
		name          string
		params        craftingCostParams
		expectedCosts []uint64
	}{
		{
			"Lunar totem",
			craftingCostParams{47.913866, 4.791387, 300, 0.2},
			[]uint64{47, 34, 32, 30, 29, 28, 28, 27, 27, 26},
		},
		{
			"Eggceptional tachyon deflector",
			craftingCostParams{1220591.960920, 122059.196092, 300, 0.2},
			[]uint64{1220591, 869525, 817323, 783258, 757357, 736215, 718227, 702498, 688475, 675791},
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s crafting costs", tt.name)
		t.Run(testname, func(t *testing.T) {
			for i, expectedCost := range tt.expectedCosts {
				cost := singleCraftCost(&tt.params, uint32(i))
				if cost != expectedCost {
					t.Errorf("cost for %d-th craft: expected %d, got %d", i+1, expectedCost, cost)
				}
			}
		})
	}
}
