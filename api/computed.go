package api

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/fanaticscripter/EggContractor/util"
)

// It is said that elite contract is unlocked at 10T% EB.
// https://egg-inc.fandom.com/wiki/Contracts says so, and it certainly matches
// my personal experience. But it's hard to confirm.
const EliteEarningBonusThreshold = 1e11

func (f *FirstContact_Payload) ApproxTime() time.Time {
	return util.DoubleToTime(f.ApproxTimestamp)
}

func (f *FirstContact_Payload) AllContractProperties() []*ContractProperties {
	s := make([]*ContractProperties, 0)
	for _, c := range f.Contracts.ActiveContracts {
		s = append(s, c.Props)
	}
	for _, c := range f.Contracts.PastContracts {
		s = append(s, c.Props)
	}
	return s
}

func (c *Contract) StartedTime() time.Time {
	return util.DoubleToTime(c.Started)
}

func (c *Contract) ProductionDeadlineTime() time.Time {
	// The production_deadline field may not be available for solo contracts.
	if c.ProductionDeadline != 0 {
		return util.DoubleToTime(c.ProductionDeadline)
	} else if !c.StartedTime().IsZero() {
		return c.StartedTime().Add(c.Props.Duration())
	}
	return time.Time{}
}

func (c *Contract) CollectionDeadlineTime() time.Time {
	// The collection_deadline field may not be available for solo contracts.
	if c.CollectionDeadline != 0 {
		return util.DoubleToTime(c.CollectionDeadline)
	}
	prodDeadline := c.ProductionDeadlineTime()
	if prodDeadline.IsZero() {
		return prodDeadline.Add(48 * time.Hour)
	}
	return time.Time{}
}

func (c *ContractProperties) Duration() time.Duration {
	return util.DoubleToDuration(c.DurationSeconds)
}

func (c *ContractProperties) ExpiryTime() time.Time {
	return util.DoubleToTime(c.ExpiryTimestamp)
}

func (c *ContractProperties) EliteRewards() []*Reward {
	if len(c.RewardTiers) == 0 {
		return c.Rewards
	}
	return c.RewardTiers[0].Rewards
}

func (c *ContractProperties) StandardRewards() []*Reward {
	if len(c.RewardTiers) == 0 {
		return c.Rewards
	}
	return c.RewardTiers[1].Rewards
}

func (c *ContractProperties) UltimateGoal(isElite bool) float64 {
	var rewards []*Reward
	if isElite {
		rewards = c.EliteRewards()
	} else {
		rewards = c.StandardRewards()
	}
	if len(rewards) == 0 {
		return 0
	}
	return rewards[len(rewards)-1].Goal
}

func (f *Farm) LastSavedTime() time.Time {
	return util.DoubleToTime(f.LastSaved)
}

func (c *CoopStatus) EggsPerSecond() float64 {
	var sum float64
	for _, m := range c.Members {
		sum += m.EggsPerSecond
	}
	return sum
}

func (c *CoopStatus) EggsPerHour() float64 {
	return 3600 * c.EggsPerSecond()
}

func (c *CoopStatus) DurationUntilProductionDeadline() time.Duration {
	return util.DoubleToDuration(c.SecondsUntilProductionDeadline)
}

func (c *CoopStatus) DurationUntilCollectionDeadline() time.Duration {
	return util.DoubleToDuration(c.SecondsUntilCollectionDeadline)
}

func (c *CoopStatus) IsElite() bool {
	var belowThresholdCnt, aboveThresholdCnt uint
	for _, m := range c.Members {
		if m.EarningBonus() >= EliteEarningBonusThreshold {
			aboveThresholdCnt++
		} else {
			belowThresholdCnt++
		}
	}
	// Ideally either one should be zero, but I can't be sure about the
	// threshold (in fact I can't even be sure the threshold is static), so play
	// it safe.
	return aboveThresholdCnt > belowThresholdCnt
}

func (c *CoopStatus) Creator() *CoopStatus_Member {
	for _, m := range c.Members {
		if m.Id == c.CreatorId {
			return m
		}
	}
	return nil
}

// RequiredEggsPerHour returns the laying rate required to complete the ultimate
// goal before the production deadline.
func (c *CoopStatus) RequiredEggsPerHour(contract *ContractProperties) float64 {
	eggsToLay := contract.UltimateGoal(c.IsElite()) - c.EggsLaid
	hoursLeft := c.DurationUntilProductionDeadline().Hours()
	if eggsToLay <= 0 || hoursLeft <= 0 {
		return 0
	} else {
		return eggsToLay / hoursLeft
	}
}

func (c *CoopStatus) ExpectedDurationUntilFinish(contract *ContractProperties) time.Duration {
	eggsToLay := contract.UltimateGoal(c.IsElite()) - c.EggsLaid
	if eggsToLay <= 0 {
		return 0
	} else if c.EggsPerSecond() <= 0 {
		return util.InfDuration // Forever
	} else {
		return util.DoubleToDuration(eggsToLay / c.EggsPerSecond())
	}
}

func (m *CoopStatus_Member) EggsPerHour() float64 {
	return 3600 * m.EggsPerSecond
}

func (m *CoopStatus_Member) EarningBonus() float64 {
	return math.Pow(10, m.EarningBonusOom)
}

func (m *CoopStatus_Member) EarningBonusPercentage() float64 {
	return m.EarningBonus() * 100
}

func (e EggType) Display() string {
	switch {
	case e == EggType_AI:
		return "AI"
	default:
		return strings.Title(strings.ReplaceAll(strings.ToLower(e.String()), "_", " "))
	}
}

func (e EggType) Value() float64 {
	switch e {
	case EggType_EDIBLE:
		return 0.1
	case EggType_SUPERFOOD:
		return 1.25
	case EggType_MEDICAL:
		return 6.25
	case EggType_ROCKET_FUEL:
		return 30
	case EggType_SUPER_MATERIAL:
		return 150
	case EggType_FUSION:
		return 700
	case EggType_QUANTUM:
		return 3_000
	case EggType_IMMORTALITY:
		return 12_500
	case EggType_TACHYON:
		return 50_000
	case EggType_GRAVITON:
		return 175_000
	case EggType_DILITHIUM:
		return 525_000
	case EggType_PRODIGY:
		return 1_500_000
	case EggType_TERRAFORM:
		return 10_000_000
	case EggType_ANTIMATTER:
		return 1e9
	case EggType_DARK_MATTER:
		return 1e11
	case EggType_AI:
		return 1e12
	case EggType_NEBULA:
		return 1.5e13
	case EggType_UNIVERSE:
		return 1e14
	case EggType_ENLIGHTENMENT:
		return 1e-7
	// Contract-only eggs.
	case EggType_CHOCOLATE:
		return 5
	case EggType_EASTER:
		return 0.05
	case EggType_WATERBALLOON:
		return 0.1
	case EggType_FIREWORK:
		return 4.99
	case EggType_PUMPKIN:
		return 0.99
	default:
		return 0
	}
}

func (e EggType) ValueDisplay() string {
	if e == EggType_ENLIGHTENMENT {
		return "0.00"
	}
	value := e.Value()
	if value >= 1e6 {
		return util.NumfmtWhole(value)
	}
	s := fmt.Sprintf("%f", value)
	if strings.Contains(s, ".") {
		// Trim unnecessary trailing zeros, and possibly also the decimal point
		// if the number turns out to be an integer.
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
	}
	return s
}

func (m *MissionInfo) StartTime() time.Time {
	return util.DoubleToTime(m.StartTimeDerived)
}

func (s MissionInfo_Spaceship) Name() string {
	switch s {
	case MissionInfo_CHICKEN_ONE:
		return "Chicken One"
	case MissionInfo_CHICKEN_NINE:
		return "Chicken Nine"
	case MissionInfo_CHICKEN_HEAVY:
		return "Chicken Heavy"
	case MissionInfo_BCR:
		return "BCR"
	case MissionInfo_MILLENIUM_CHICKEN:
		return "Quintillion Chicken"
	case MissionInfo_CORELLIHEN_CORVETTE:
		return "Cornish-Hen Corvette"
	case MissionInfo_GALEGGTICA:
		return "Galeggtica"
	case MissionInfo_CHICKFIANT:
		return "Defihent"
	case MissionInfo_VOYEGGER:
		return "Voyegger"
	case MissionInfo_HENERPRISE:
		return "Henerprise"
	}
	return "Unknown"
}

func (s MissionInfo_Status) Display() string {
	lower := strings.ReplaceAll(strings.ToLower(s.String()), "_", " ")
	return strings.ToUpper(lower[:1]) + lower[1:]
}

func (d MissionInfo_DurationType) Display() string {
	switch d {
	case MissionInfo_TUTORIAL:
		return "Tutorial"
	case MissionInfo_SHORT:
		return "Short"
	case MissionInfo_LONG:
		return "Standard"
	case MissionInfo_EPIC:
		return "Extended"
	}
	return "Unknown"
}

// GameName is in all caps. Use CasedName for cased version.
func (a ArtifactSpec_Name) GameName() string {
	name := strings.ReplaceAll(a.String(), "_", " ")
	switch a {
	case ArtifactSpec_VIAL_MARTIAN_DUST:
		name = "VIAL OF MARTIAN DUST"
	case ArtifactSpec_ORNATE_GUSSET:
		name = "GUSSET"
	case ArtifactSpec_MERCURYS_LENS:
		name = "MERCURY'S LENS"
	}
	return name
}

func (a ArtifactSpec_Name) CasedName() string {
	return capitalizeArtifactName(strings.ToLower(a.GameName()))
}

// GameName is in all caps. Use CasedName for cased version.
func (a *ArtifactSpec) GameName() string {
	switch a.Name {
	// Artifacts
	case ArtifactSpec_LUNAR_TOTEM:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "BASIC LUNAR TOTEM"
		case ArtifactSpec_LESSER:
			return "LUNAR TOTEM"
		case ArtifactSpec_NORMAL:
			return "POWERFUL LUNAR TOTEM"
		case ArtifactSpec_GREATER:
			return "EGGCEPTIONAL LUNAR TOTEM"
		}
	case ArtifactSpec_NEODYMIUM_MEDALLION:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "WEAK NEODYMIUM MEDALLION"
		case ArtifactSpec_LESSER:
			return "NEODYMIUM MEDALLION"
		case ArtifactSpec_NORMAL:
			return "PRECISE NEODYMIUM MEDALLION"
		case ArtifactSpec_GREATER:
			return "EGGCEPTIONAL NEODYMIUM MEDALLION"
		}
	case ArtifactSpec_BEAK_OF_MIDAS:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "DULL BEAK OF MIDAS"
		case ArtifactSpec_LESSER:
			return "BEAK OF MIDAS"
		case ArtifactSpec_NORMAL:
			return "JEWELED BEAK OF MIDAS"
		case ArtifactSpec_GREATER:
			return "GLISTENING BEAK OF MIDAS"
		}
	case ArtifactSpec_LIGHT_OF_EGGENDIL:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "DIM LIGHT OF EGGENDIL"
		case ArtifactSpec_LESSER:
			return "SHIMMERING LIGHT OF EGGENDIL"
		case ArtifactSpec_NORMAL:
			return "GLOWING LIGHT OF EGGENDIL"
		case ArtifactSpec_GREATER:
			return "BRILLIANT LIGHT OF EGGENDIL"
		}
	case ArtifactSpec_DEMETERS_NECKLACE:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "SIMPLE DEMETERS NECKLACE"
		case ArtifactSpec_LESSER:
			return "JEWELED DEMETERS NECKLACE"
		case ArtifactSpec_NORMAL:
			return "PRISTINE DEMETERS NECKLACE"
		case ArtifactSpec_GREATER:
			return "BEGGSPOKE DEMETERS NECKLACE"
		}
	case ArtifactSpec_VIAL_MARTIAN_DUST:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "TINY VIAL OF MARTIAN DUST"
		case ArtifactSpec_LESSER:
			return "VIAL OF MARTIAN DUST"
		case ArtifactSpec_NORMAL:
			return "HERMETIC VIAL OF MARTIAN DUST"
		case ArtifactSpec_GREATER:
			return "PRIME VIAL OF MARTIAN DUST"
		}
	case ArtifactSpec_ORNATE_GUSSET:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "PLAIN GUSSET"
		case ArtifactSpec_LESSER:
			return "ORNATE GUSSET"
		case ArtifactSpec_NORMAL:
			return "DISTEGGUISHED GUSSET"
		case ArtifactSpec_GREATER:
			return "JEWELED GUSSET"
		}
	case ArtifactSpec_THE_CHALICE:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "PLAIN CHALICE"
		case ArtifactSpec_LESSER:
			return "POLISHED CHALICE"
		case ArtifactSpec_NORMAL:
			return "JEWELED CHALICE"
		case ArtifactSpec_GREATER:
			return "EGGCEPTIONAL CHALICE"
		}
	case ArtifactSpec_BOOK_OF_BASAN:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "BOOK OF BASAN"
		case ArtifactSpec_LESSER:
			return "COLLECTORS BOOK OF BASAN"
		case ArtifactSpec_NORMAL:
			return "FORTIFIED BOOK OF BASAN"
		case ArtifactSpec_GREATER:
			return "GILDED BOOK OF BASAN"
		}
	case ArtifactSpec_PHOENIX_FEATHER:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "TATTERED PHOENIX FEATHER"
		case ArtifactSpec_LESSER:
			return "PHOENIX FEATHER"
		case ArtifactSpec_NORMAL:
			return "BRILLIANT PHOENIX FEATHER"
		case ArtifactSpec_GREATER:
			return "BLAZING PHOENIX FEATHER"
		}
	case ArtifactSpec_TUNGSTEN_ANKH:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "CRUDE TUNGSTEN ANKH"
		case ArtifactSpec_LESSER:
			return "TUNGSTEN ANKH"
		case ArtifactSpec_NORMAL:
			return "POLISHED TUNGSTEN ANKH"
		case ArtifactSpec_GREATER:
			return "BRILLIANT TUNGSTEN ANKH"
		}
	case ArtifactSpec_AURELIAN_BROOCH:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "PLAIN AURELIAN BROOCH"
		case ArtifactSpec_LESSER:
			return "AURELIAN BROOCH"
		case ArtifactSpec_NORMAL:
			return "JEWELED AURELIAN BROOCH"
		case ArtifactSpec_GREATER:
			return "EGGCEPTIONAL AURELIAN BROOCH"
		}
	case ArtifactSpec_CARVED_RAINSTICK:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "SIMPLE CARVED RAINSTICK"
		case ArtifactSpec_LESSER:
			return "CARVED RAINSTICK"
		case ArtifactSpec_NORMAL:
			return "ORNATE CARVED RAINSTICK"
		case ArtifactSpec_GREATER:
			return "MEGGNIFICENT CARVED RAINSTICK"
		}
	case ArtifactSpec_PUZZLE_CUBE:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "ANCIENT PUZZLE CUBE"
		case ArtifactSpec_LESSER:
			return "PUZZLE CUBE"
		case ArtifactSpec_NORMAL:
			return "MYSTICAL PUZZLE CUBE"
		case ArtifactSpec_GREATER:
			return "UNSOLVABLE PUZZLE CUBE"
		}
	case ArtifactSpec_QUANTUM_METRONOME:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "MISALIGNED QUANTUM METRONOME"
		case ArtifactSpec_LESSER:
			return "ADEQUATE QUANTUM METRONOME"
		case ArtifactSpec_NORMAL:
			return "PERFECT QUANTUM METRONOME"
		case ArtifactSpec_GREATER:
			return "REGGFERENCE QUANTUM METRONOME"
		}
	case ArtifactSpec_SHIP_IN_A_BOTTLE:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "SHIP IN A BOTTLE"
		case ArtifactSpec_LESSER:
			return "DETAILED SHIP IN A BOTTLE"
		case ArtifactSpec_NORMAL:
			return "COMPLEX SHIP IN A BOTTLE"
		case ArtifactSpec_GREATER:
			return "EGGQUISITE SHIP IN A BOTTLE"
		}
	case ArtifactSpec_TACHYON_DEFLECTOR:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "WEAK TACHYON DEFLECTOR"
		case ArtifactSpec_LESSER:
			return "TACHYON DEFLECTOR"
		case ArtifactSpec_NORMAL:
			return "ROBUST TACHYON DEFLECTOR"
		case ArtifactSpec_GREATER:
			return "EGGCEPTIONAL TACHYON DEFLECTOR"
		}
	case ArtifactSpec_INTERSTELLAR_COMPASS:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "MISCALIBRATED INTERSTELLAR COMPASS"
		case ArtifactSpec_LESSER:
			return "INTERSTELLAR COMPASS"
		case ArtifactSpec_NORMAL:
			return "PRECISE INTERSTELLAR COMPASS"
		case ArtifactSpec_GREATER:
			return "CLAIRVOYANT INTERSTELLAR COMPASS"
		}
	case ArtifactSpec_DILITHIUM_MONOCLE:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "DILITHIUM MONOCLE"
		case ArtifactSpec_LESSER:
			return "PRECISE DILITHIUM MONOCLE"
		case ArtifactSpec_NORMAL:
			return "EGGSACTING DILITHIUM MONOCLE"
		case ArtifactSpec_GREATER:
			return "FLAWLESS DILITHIUM MONOCLE"
		}
	case ArtifactSpec_TITANIUM_ACTUATOR:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "INCONSISTENT TITANIUM ACTUATOR"
		case ArtifactSpec_LESSER:
			return "TITANIUM ACTUATOR"
		case ArtifactSpec_NORMAL:
			return "PRECISE TITANIUM ACTUATOR"
		case ArtifactSpec_GREATER:
			return "REGGFERENCE TITANIUM ACTUATOR"
		}
	case ArtifactSpec_MERCURYS_LENS:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "MISALIGNED MERCURY'S LENS"
		case ArtifactSpec_LESSER:
			return "MERCURY'S LENS"
		case ArtifactSpec_NORMAL:
			return "PRECISE MERCURY'S LENS"
		case ArtifactSpec_GREATER:
			return "MEGGNIFICENT MERCURY'S LENS"
		}
	// Stones
	case ArtifactSpec_TACHYON_STONE:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "TACHYON STONE"
		case ArtifactSpec_LESSER:
			return "EGGSQUISITE TACHYON STONE"
		case ArtifactSpec_NORMAL:
			return "BRILLIANT TACHYON STONE"
		}
	case ArtifactSpec_DILITHIUM_STONE:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "DILITHIUM STONE"
		case ArtifactSpec_LESSER:
			return "EGGSQUISITE DILITHIUM STONE"
		case ArtifactSpec_NORMAL:
			return "BRILLIANT DILITHIUM STONE"
		}
	case ArtifactSpec_SHELL_STONE:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "SHELL STONE"
		case ArtifactSpec_LESSER:
			return "EGGSQUISITE SHELL STONE"
		case ArtifactSpec_NORMAL:
			return "FLAWLESS SHELL STONE"
		}
	case ArtifactSpec_LUNAR_STONE:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "LUNAR STONE"
		case ArtifactSpec_LESSER:
			return "EGGSQUISITE LUNAR STONE"
		case ArtifactSpec_NORMAL:
			return "MEGGNIFICENT LUNAR STONE"
		}
	case ArtifactSpec_SOUL_STONE:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "SOUL STONE"
		case ArtifactSpec_LESSER:
			return "EGGSQUISITE SOUL STONE"
		case ArtifactSpec_NORMAL:
			return "RADIANT SOUL STONE"
		}
	case ArtifactSpec_PROPHECY_STONE:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "PROPHECY STONE"
		case ArtifactSpec_LESSER:
			return "EGGSQUISITE PROPHECY STONE"
		case ArtifactSpec_NORMAL:
			return "RADIANT PROPHECY STONE"
		}
	case ArtifactSpec_QUANTUM_STONE:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "QUANTUM STONE"
		case ArtifactSpec_LESSER:
			return "PHASED QUANTUM STONE"
		case ArtifactSpec_NORMAL:
			return "MEGGNIFICENT QUANTUM STONE"
		}
	case ArtifactSpec_TERRA_STONE:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "TERRA STONE"
		case ArtifactSpec_LESSER:
			return "RICH TERRA STONE"
		case ArtifactSpec_NORMAL:
			return "EGGCEPTIONAL TERRA STONE"
		}
	case ArtifactSpec_LIFE_STONE:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "LIFE STONE"
		case ArtifactSpec_LESSER:
			return "GOOD LIFE STONE"
		case ArtifactSpec_NORMAL:
			return "EGGCEPTIONAL LIFE STONE"
		}
	case ArtifactSpec_CLARITY_STONE:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "CLARITY STONE"
		case ArtifactSpec_LESSER:
			return "EGGSQUISITE CLARITY STONE"
		case ArtifactSpec_NORMAL:
			return "EGGCEPTIONAL CLARITY STONE"
		}
	// Stone fragments
	case ArtifactSpec_TACHYON_STONE_FRAGMENT:
		fallthrough
	case ArtifactSpec_DILITHIUM_STONE_FRAGMENT:
		fallthrough
	case ArtifactSpec_SHELL_STONE_FRAGMENT:
		fallthrough
	case ArtifactSpec_LUNAR_STONE_FRAGMENT:
		fallthrough
	case ArtifactSpec_SOUL_STONE_FRAGMENT:
		fallthrough
	case ArtifactSpec_PROPHECY_STONE_FRAGMENT:
		fallthrough
	case ArtifactSpec_QUANTUM_STONE_FRAGMENT:
		fallthrough
	case ArtifactSpec_TERRA_STONE_FRAGMENT:
		fallthrough
	case ArtifactSpec_LIFE_STONE_FRAGMENT:
		fallthrough
	case ArtifactSpec_CLARITY_STONE_FRAGMENT:
		return strings.ReplaceAll(a.Name.String(), "_", " ")
	// Ingredients
	case ArtifactSpec_GOLD_METEORITE:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "TINY GOLD METEORITE"
		case ArtifactSpec_LESSER:
			return "ENRICHED GOLD METEORITE"
		case ArtifactSpec_NORMAL:
			return "SOLID GOLD METEORITE"
		}
	case ArtifactSpec_TAU_CETI_GEODE:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "TAU CETI GEODE PIECE"
		case ArtifactSpec_LESSER:
			return "GLIMMERING TAU CETI GEODE"
		case ArtifactSpec_NORMAL:
			return "RADIANT TAU CETI GEODE"
		}
	case ArtifactSpec_SOLAR_TITANIUM:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "SOLAR TITANIUM ORE"
		case ArtifactSpec_LESSER:
			return "SOLAR TITANIUM BAR"
		case ArtifactSpec_NORMAL:
			return "SOLAR TITANIUM GEOGON"
		}
	// Unconfirmed ingredients
	case ArtifactSpec_EXTRATERRESTRIAL_ALUMINUM:
		fallthrough
	case ArtifactSpec_ANCIENT_TUNGSTEN:
		fallthrough
	case ArtifactSpec_SPACE_ROCKS:
		fallthrough
	case ArtifactSpec_ALIEN_WOOD:
		fallthrough
	case ArtifactSpec_CENTAURIAN_STEEL:
		fallthrough
	case ArtifactSpec_ERIDANI_FEATHER:
		fallthrough
	case ArtifactSpec_DRONE_PARTS:
		fallthrough
	case ArtifactSpec_CELESTIAL_BRONZE:
		fallthrough
	case ArtifactSpec_LALANDE_HIDE:
		return "? " + a.Name.String()
	}
	return a.Level.String() + " " + a.Name.GameName()
}

func (a *ArtifactSpec) CasedName() string {
	return capitalizeArtifactName(strings.ToLower(a.GameName()))
}

func capitalizeArtifactName(n string) string {
	n = strings.ToUpper(n[:1]) + n[1:]
	// Captalize proper nouns.
	for s, repl := range map[string]string{
		"midas":    "Midas",
		"eggendil": "Eggendil",
		"martian":  "Martian",
		"basan":    "Basan",
		"aurelian": "Aurelian",
		"mercury":  "Mercury",
		"tau ceti": "Tau Ceti",
		"Tau ceti": "Tau Ceti",
	} {
		n = strings.ReplaceAll(n, s, repl)
	}
	return n
}

func (a *ArtifactSpec) Type() ArtifactSpec_Type {
	return a.Name.ArtifactType()
}

func (a ArtifactSpec_Name) ArtifactType() ArtifactSpec_Type {
	switch a {
	// Artifacts
	case ArtifactSpec_LUNAR_TOTEM:
		fallthrough
	case ArtifactSpec_NEODYMIUM_MEDALLION:
		fallthrough
	case ArtifactSpec_BEAK_OF_MIDAS:
		fallthrough
	case ArtifactSpec_LIGHT_OF_EGGENDIL:
		fallthrough
	case ArtifactSpec_DEMETERS_NECKLACE:
		fallthrough
	case ArtifactSpec_VIAL_MARTIAN_DUST:
		fallthrough
	case ArtifactSpec_ORNATE_GUSSET:
		fallthrough
	case ArtifactSpec_THE_CHALICE:
		fallthrough
	case ArtifactSpec_BOOK_OF_BASAN:
		fallthrough
	case ArtifactSpec_PHOENIX_FEATHER:
		fallthrough
	case ArtifactSpec_TUNGSTEN_ANKH:
		fallthrough
	case ArtifactSpec_AURELIAN_BROOCH:
		fallthrough
	case ArtifactSpec_CARVED_RAINSTICK:
		fallthrough
	case ArtifactSpec_PUZZLE_CUBE:
		fallthrough
	case ArtifactSpec_QUANTUM_METRONOME:
		fallthrough
	case ArtifactSpec_SHIP_IN_A_BOTTLE:
		fallthrough
	case ArtifactSpec_TACHYON_DEFLECTOR:
		fallthrough
	case ArtifactSpec_INTERSTELLAR_COMPASS:
		fallthrough
	case ArtifactSpec_DILITHIUM_MONOCLE:
		fallthrough
	case ArtifactSpec_TITANIUM_ACTUATOR:
		fallthrough
	case ArtifactSpec_MERCURYS_LENS:
		return ArtifactSpec_ARTIFACT
	// Stones
	case ArtifactSpec_TACHYON_STONE:
		fallthrough
	case ArtifactSpec_DILITHIUM_STONE:
		fallthrough
	case ArtifactSpec_SHELL_STONE:
		fallthrough
	case ArtifactSpec_LUNAR_STONE:
		fallthrough
	case ArtifactSpec_SOUL_STONE:
		fallthrough
	case ArtifactSpec_PROPHECY_STONE:
		fallthrough
	case ArtifactSpec_QUANTUM_STONE:
		fallthrough
	case ArtifactSpec_TERRA_STONE:
		fallthrough
	case ArtifactSpec_LIFE_STONE:
		fallthrough
	case ArtifactSpec_CLARITY_STONE:
		return ArtifactSpec_STONE
	// Stone fragments
	case ArtifactSpec_TACHYON_STONE_FRAGMENT:
		fallthrough
	case ArtifactSpec_DILITHIUM_STONE_FRAGMENT:
		fallthrough
	case ArtifactSpec_SHELL_STONE_FRAGMENT:
		fallthrough
	case ArtifactSpec_LUNAR_STONE_FRAGMENT:
		fallthrough
	case ArtifactSpec_SOUL_STONE_FRAGMENT:
		fallthrough
	case ArtifactSpec_PROPHECY_STONE_FRAGMENT:
		fallthrough
	case ArtifactSpec_QUANTUM_STONE_FRAGMENT:
		fallthrough
	case ArtifactSpec_TERRA_STONE_FRAGMENT:
		fallthrough
	case ArtifactSpec_LIFE_STONE_FRAGMENT:
		fallthrough
	case ArtifactSpec_CLARITY_STONE_FRAGMENT:
		return ArtifactSpec_STONE_INGREDIENT
	// Ingredients
	case ArtifactSpec_GOLD_METEORITE:
		fallthrough
	case ArtifactSpec_TAU_CETI_GEODE:
		fallthrough
	case ArtifactSpec_SOLAR_TITANIUM:
		fallthrough
	// Unconfirmed ingredients
	case ArtifactSpec_EXTRATERRESTRIAL_ALUMINUM:
		fallthrough
	case ArtifactSpec_ANCIENT_TUNGSTEN:
		fallthrough
	case ArtifactSpec_SPACE_ROCKS:
		fallthrough
	case ArtifactSpec_ALIEN_WOOD:
		fallthrough
	case ArtifactSpec_CENTAURIAN_STEEL:
		fallthrough
	case ArtifactSpec_ERIDANI_FEATHER:
		fallthrough
	case ArtifactSpec_DRONE_PARTS:
		fallthrough
	case ArtifactSpec_CELESTIAL_BRONZE:
		fallthrough
	case ArtifactSpec_LALANDE_HIDE:
		return ArtifactSpec_INGREDIENT
	}
	return ArtifactSpec_ARTIFACT
}

// Family returns the family the artifact belongs to, which is the corresponding
// stone for stone fragments.
func (a *ArtifactSpec) Family() ArtifactSpec_Name {
	return a.Name.Family()
}

// Family returns the family of the artifact, which is simply itself other than
// when it is a stone fragment, in which case the corresponding stone is
// returned.
func (a ArtifactSpec_Name) Family() ArtifactSpec_Name {
	if a.ArtifactType() == ArtifactSpec_STONE_INGREDIENT {
		return a.CorrespondingStone()
	}
	return a
}

// CorrespondingStone returns the corresponding stone for a stone fragment.
// Result is undefined for non-stone fragments.
func (a ArtifactSpec_Name) CorrespondingStone() ArtifactSpec_Name {
	switch a {
	case ArtifactSpec_TACHYON_STONE_FRAGMENT:
		return ArtifactSpec_TACHYON_STONE
	case ArtifactSpec_DILITHIUM_STONE_FRAGMENT:
		return ArtifactSpec_DILITHIUM_STONE
	case ArtifactSpec_SHELL_STONE_FRAGMENT:
		return ArtifactSpec_SHELL_STONE
	case ArtifactSpec_LUNAR_STONE_FRAGMENT:
		return ArtifactSpec_LUNAR_STONE
	case ArtifactSpec_SOUL_STONE_FRAGMENT:
		return ArtifactSpec_SOUL_STONE
	case ArtifactSpec_PROPHECY_STONE_FRAGMENT:
		return ArtifactSpec_PROPHECY_STONE
	case ArtifactSpec_QUANTUM_STONE_FRAGMENT:
		return ArtifactSpec_QUANTUM_STONE
	case ArtifactSpec_TERRA_STONE_FRAGMENT:
		return ArtifactSpec_TERRA_STONE
	case ArtifactSpec_LIFE_STONE_FRAGMENT:
		return ArtifactSpec_LIFE_STONE
	case ArtifactSpec_CLARITY_STONE_FRAGMENT:
		return ArtifactSpec_CLARITY_STONE
	}
	return ArtifactSpec_UNKNOWN
}

// CorrespondingFragment returns the corresponding stone fragment for a stone.
// Result is undefined for non-stones.
func (a ArtifactSpec_Name) CorrespondingFragment() ArtifactSpec_Name {
	switch a {
	case ArtifactSpec_TACHYON_STONE:
		return ArtifactSpec_TACHYON_STONE_FRAGMENT
	case ArtifactSpec_DILITHIUM_STONE:
		return ArtifactSpec_DILITHIUM_STONE_FRAGMENT
	case ArtifactSpec_SHELL_STONE:
		return ArtifactSpec_SHELL_STONE_FRAGMENT
	case ArtifactSpec_LUNAR_STONE:
		return ArtifactSpec_LUNAR_STONE_FRAGMENT
	case ArtifactSpec_SOUL_STONE:
		return ArtifactSpec_SOUL_STONE_FRAGMENT
	case ArtifactSpec_PROPHECY_STONE:
		return ArtifactSpec_PROPHECY_STONE_FRAGMENT
	case ArtifactSpec_QUANTUM_STONE:
		return ArtifactSpec_QUANTUM_STONE_FRAGMENT
	case ArtifactSpec_TERRA_STONE:
		return ArtifactSpec_TERRA_STONE_FRAGMENT
	case ArtifactSpec_LIFE_STONE:
		return ArtifactSpec_LIFE_STONE_FRAGMENT
	case ArtifactSpec_CLARITY_STONE:
		return ArtifactSpec_CLARITY_STONE_FRAGMENT
	}
	return ArtifactSpec_UNKNOWN
}

func (a *ArtifactSpec) TierNumber() int {
	switch a.Type() {
	case ArtifactSpec_ARTIFACT:
		// 0, 1, 2, 3 => T1, T2, T3, T4
		return int(a.Level) + 1
	case ArtifactSpec_STONE:
		// 0, 1, 2 => T2, T3, T4 (fragment as T1)
		return int(a.Level) + 2
	case ArtifactSpec_STONE_INGREDIENT:
		return 1
	case ArtifactSpec_INGREDIENT:
		// 0, 1, 2 => T1, T2, T3
		return int(a.Level) + 1
	}
	return 1
}

func (a *ArtifactSpec) TierName() string {
	switch a.Name {
	// Artifacts
	case ArtifactSpec_LUNAR_TOTEM:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "BASIC"
		case ArtifactSpec_LESSER:
			return "REGULAR"
		case ArtifactSpec_NORMAL:
			return "POWERFUL"
		case ArtifactSpec_GREATER:
			return "EGGCEPTIONAL"
		}
	case ArtifactSpec_NEODYMIUM_MEDALLION:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "WEAK"
		case ArtifactSpec_LESSER:
			return "REGULAR"
		case ArtifactSpec_NORMAL:
			return "PRECISE"
		case ArtifactSpec_GREATER:
			return "EGGCEPTIONAL"
		}
	case ArtifactSpec_BEAK_OF_MIDAS:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "DULL"
		case ArtifactSpec_LESSER:
			return "REGULAR"
		case ArtifactSpec_NORMAL:
			return "JEWELED"
		case ArtifactSpec_GREATER:
			return "GLISTENING"
		}
	case ArtifactSpec_LIGHT_OF_EGGENDIL:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "DIM"
		case ArtifactSpec_LESSER:
			return "SHIMMERING"
		case ArtifactSpec_NORMAL:
			return "GLOWING"
		case ArtifactSpec_GREATER:
			return "BRILLIANT"
		}
	case ArtifactSpec_DEMETERS_NECKLACE:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "SIMPLE"
		case ArtifactSpec_LESSER:
			return "JEWELED"
		case ArtifactSpec_NORMAL:
			return "PRISTINE"
		case ArtifactSpec_GREATER:
			return "BEGGSPOKE"
		}
	case ArtifactSpec_VIAL_MARTIAN_DUST:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "TINY"
		case ArtifactSpec_LESSER:
			return "REGULAR"
		case ArtifactSpec_NORMAL:
			return "HERMETIC"
		case ArtifactSpec_GREATER:
			return "PRIME"
		}
	case ArtifactSpec_ORNATE_GUSSET:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "PLAIN"
		case ArtifactSpec_LESSER:
			return "ORNATE"
		case ArtifactSpec_NORMAL:
			return "DISTEGGUISHED"
		case ArtifactSpec_GREATER:
			return "JEWELED"
		}
	case ArtifactSpec_THE_CHALICE:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "PLAIN"
		case ArtifactSpec_LESSER:
			return "POLISHED"
		case ArtifactSpec_NORMAL:
			return "JEWELED"
		case ArtifactSpec_GREATER:
			return "EGGCEPTIONAL"
		}
	case ArtifactSpec_BOOK_OF_BASAN:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "REGULAR"
		case ArtifactSpec_LESSER:
			return "COLLECTORS"
		case ArtifactSpec_NORMAL:
			return "FORTIFIED"
		case ArtifactSpec_GREATER:
			return "GILDED"
		}
	case ArtifactSpec_PHOENIX_FEATHER:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "TATTERED"
		case ArtifactSpec_LESSER:
			return "REGULAR"
		case ArtifactSpec_NORMAL:
			return "BRILLIANT"
		case ArtifactSpec_GREATER:
			return "BLAZING"
		}
	case ArtifactSpec_TUNGSTEN_ANKH:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "CRUDE"
		case ArtifactSpec_LESSER:
			return "REGULAR"
		case ArtifactSpec_NORMAL:
			return "POLISHED"
		case ArtifactSpec_GREATER:
			return "BRILLIANT"
		}
	case ArtifactSpec_AURELIAN_BROOCH:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "PLAIN"
		case ArtifactSpec_LESSER:
			return "REGULAR"
		case ArtifactSpec_NORMAL:
			return "JEWELED"
		case ArtifactSpec_GREATER:
			return "EGGCEPTIONAL"
		}
	case ArtifactSpec_CARVED_RAINSTICK:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "SIMPLE"
		case ArtifactSpec_LESSER:
			return "REGULAR"
		case ArtifactSpec_NORMAL:
			return "ORNATE"
		case ArtifactSpec_GREATER:
			return "MEGGNIFICENT"
		}
	case ArtifactSpec_PUZZLE_CUBE:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "ANCIENT"
		case ArtifactSpec_LESSER:
			return "REGULAR"
		case ArtifactSpec_NORMAL:
			return "MYSTICAL"
		case ArtifactSpec_GREATER:
			return "UNSOLVABLE"
		}
	case ArtifactSpec_QUANTUM_METRONOME:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "MISALIGNED"
		case ArtifactSpec_LESSER:
			return "ADEQUATE"
		case ArtifactSpec_NORMAL:
			return "PERFECT"
		case ArtifactSpec_GREATER:
			return "REGGFERENCE"
		}
	case ArtifactSpec_SHIP_IN_A_BOTTLE:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "REGULAR"
		case ArtifactSpec_LESSER:
			return "DETAILED"
		case ArtifactSpec_NORMAL:
			return "COMPLEX"
		case ArtifactSpec_GREATER:
			return "EGGQUISITE"
		}
	case ArtifactSpec_TACHYON_DEFLECTOR:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "WEAK"
		case ArtifactSpec_LESSER:
			return "REGULAR"
		case ArtifactSpec_NORMAL:
			return "ROBUST"
		case ArtifactSpec_GREATER:
			return "EGGCEPTIONAL"
		}
	case ArtifactSpec_INTERSTELLAR_COMPASS:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "MISCALIBRATED"
		case ArtifactSpec_LESSER:
			return "REGULAR"
		case ArtifactSpec_NORMAL:
			return "PRECISE"
		case ArtifactSpec_GREATER:
			return "CLAIRVOYANT"
		}
	case ArtifactSpec_DILITHIUM_MONOCLE:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "REGULAR"
		case ArtifactSpec_LESSER:
			return "PRECISE"
		case ArtifactSpec_NORMAL:
			return "EGGSACTING"
		case ArtifactSpec_GREATER:
			return "FLAWLESS"
		}
	case ArtifactSpec_TITANIUM_ACTUATOR:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "INCONSISTENT"
		case ArtifactSpec_LESSER:
			return "REGULAR"
		case ArtifactSpec_NORMAL:
			return "PRECISE"
		case ArtifactSpec_GREATER:
			return "REGGFERENCE"
		}
	case ArtifactSpec_MERCURYS_LENS:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "MISALIGNED"
		case ArtifactSpec_LESSER:
			return "REGULAR"
		case ArtifactSpec_NORMAL:
			return "PRECISE"
		case ArtifactSpec_GREATER:
			return "MEGGNIFICENT"
		}
	// Stones
	case ArtifactSpec_TACHYON_STONE:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "REGULAR"
		case ArtifactSpec_LESSER:
			return "EGGSQUISITE"
		case ArtifactSpec_NORMAL:
			return "BRILLIANT"
		}
	case ArtifactSpec_DILITHIUM_STONE:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "REGULAR"
		case ArtifactSpec_LESSER:
			return "EGGSQUISITE"
		case ArtifactSpec_NORMAL:
			return "BRILLIANT"
		}
	case ArtifactSpec_SHELL_STONE:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "REGULAR"
		case ArtifactSpec_LESSER:
			return "EGGSQUISITE"
		case ArtifactSpec_NORMAL:
			return "FLAWLESS"
		}
	case ArtifactSpec_LUNAR_STONE:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "REGULAR"
		case ArtifactSpec_LESSER:
			return "EGGSQUISITE"
		case ArtifactSpec_NORMAL:
			return "MEGGNIFICENT"
		}
	case ArtifactSpec_SOUL_STONE:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "REGULAR"
		case ArtifactSpec_LESSER:
			return "EGGSQUISITE"
		case ArtifactSpec_NORMAL:
			return "RADIANT"
		}
	case ArtifactSpec_PROPHECY_STONE:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "REGULAR"
		case ArtifactSpec_LESSER:
			return "EGGSQUISITE"
		case ArtifactSpec_NORMAL:
			return "RADIANT"
		}
	case ArtifactSpec_QUANTUM_STONE:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "REGULAR"
		case ArtifactSpec_LESSER:
			return "PHASED"
		case ArtifactSpec_NORMAL:
			return "MEGGNIFICENT"
		}
	case ArtifactSpec_TERRA_STONE:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "REGULAR"
		case ArtifactSpec_LESSER:
			return "RICH"
		case ArtifactSpec_NORMAL:
			return "EGGCEPTIONAL"
		}
	case ArtifactSpec_LIFE_STONE:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "REGULAR"
		case ArtifactSpec_LESSER:
			return "GOOD"
		case ArtifactSpec_NORMAL:
			return "EGGCEPTIONAL"
		}
	case ArtifactSpec_CLARITY_STONE:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "REGULAR"
		case ArtifactSpec_LESSER:
			return "EGGSQUISITE"
		case ArtifactSpec_NORMAL:
			return "EGGCEPTIONAL"
		}
	// Stone fragments
	case ArtifactSpec_TACHYON_STONE_FRAGMENT:
		fallthrough
	case ArtifactSpec_DILITHIUM_STONE_FRAGMENT:
		fallthrough
	case ArtifactSpec_SHELL_STONE_FRAGMENT:
		fallthrough
	case ArtifactSpec_LUNAR_STONE_FRAGMENT:
		fallthrough
	case ArtifactSpec_SOUL_STONE_FRAGMENT:
		fallthrough
	case ArtifactSpec_PROPHECY_STONE_FRAGMENT:
		fallthrough
	case ArtifactSpec_QUANTUM_STONE_FRAGMENT:
		fallthrough
	case ArtifactSpec_TERRA_STONE_FRAGMENT:
		fallthrough
	case ArtifactSpec_LIFE_STONE_FRAGMENT:
		fallthrough
	case ArtifactSpec_CLARITY_STONE_FRAGMENT:
		return "FRAGMENT"
	// Ingredients
	case ArtifactSpec_GOLD_METEORITE:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "TINY"
		case ArtifactSpec_LESSER:
			return "ENRICHED"
		case ArtifactSpec_NORMAL:
			return "SOLID"
		}
	case ArtifactSpec_TAU_CETI_GEODE:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "TAU"
		case ArtifactSpec_LESSER:
			return "GLIMMERING"
		case ArtifactSpec_NORMAL:
			return "RADIANT"
		}
	case ArtifactSpec_SOLAR_TITANIUM:
		switch a.Level {
		case ArtifactSpec_INFERIOR:
			return "ORE"
		case ArtifactSpec_LESSER:
			return "BAR"
		case ArtifactSpec_NORMAL:
			return "GEOGON"
		}
	// Unconfirmed ingredients
	case ArtifactSpec_EXTRATERRESTRIAL_ALUMINUM:
		fallthrough
	case ArtifactSpec_ANCIENT_TUNGSTEN:
		fallthrough
	case ArtifactSpec_SPACE_ROCKS:
		fallthrough
	case ArtifactSpec_ALIEN_WOOD:
		fallthrough
	case ArtifactSpec_CENTAURIAN_STEEL:
		fallthrough
	case ArtifactSpec_ERIDANI_FEATHER:
		fallthrough
	case ArtifactSpec_DRONE_PARTS:
		fallthrough
	case ArtifactSpec_CELESTIAL_BRONZE:
		fallthrough
	case ArtifactSpec_LALANDE_HIDE:
		return "?"
	}
	return "?"
}

func (a *ArtifactSpec) CasedTierName() string {
	return strings.Title(strings.ToLower(a.TierName()))
}

func (r ArtifactSpec_Rarity) Display() string {
	switch r {
	case ArtifactSpec_COMMON:
		return "Common"
	case ArtifactSpec_RARE:
		return "Rare"
	case ArtifactSpec_EPIC:
		return "Epic"
	case ArtifactSpec_LEGENDARY:
		return "Legendary"
	}
	return "Unknown"
}
