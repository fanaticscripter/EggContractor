package pb

import (
	"time"

	"github.com/fanaticscripter/EggContractor/util"
)

func (c *SoloContract) GetDurationUntilProductionDeadline() time.Duration {
	return util.DoubleToDuration(c.SecondsUntilProductionDeadline)
}

func (c *SoloContract) GetDurationUntilCollectionDeadline() time.Duration {
	return util.DoubleToDuration(c.SecondsUntilCollectionDeadline)
}

func (c *SoloContract) GetLastRefreshedTime() time.Time {
	return util.DoubleToTime(c.LastRefreshedTimestamp)
}
