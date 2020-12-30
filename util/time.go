package util

import (
	"fmt"
	"math"
	"time"
)

const InfDuration = time.Duration(math.MaxInt64)

// DoubleToTime converts an epoch timestamp in seconds to Time.
// Overflow is not considered.
func DoubleToTime(t float64) time.Time {
	if t == 0 {
		return time.Time{}
	}
	sec, frac := math.Modf(t)
	return time.Unix(int64(sec), int64(frac*1e9))
}

func TimeToDouble(t time.Time) float64 {
	return float64(t.UnixNano()) / 1e9
}

// DoubleToDuration converts a float64 number of seconds to Duration.
// +/-InfDuration is returned in the event of overflow.
func DoubleToDuration(d float64) time.Duration {
	if d < 0 {
		return -DoubleToDuration(-d)
	}
	if d <= float64(math.MaxInt64)/1e9 {
		return time.Duration(d * 1e9)
	} else {
		return InfDuration
	}
}

func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func FormatDatetime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func FormatDuration(d time.Duration) string {
	if d < 0 {
		return "-" + FormatDuration(-d)
	}
	if d == InfDuration {
		return "forever"
	}
	dd := d / (24 * time.Hour)
	d -= dd * 24 * time.Hour
	hh := d / time.Hour
	d -= hh * time.Hour
	mm := d / time.Minute
	return fmt.Sprintf("%dd%dh%dm", dd, hh, mm)
}

func FormatDurationNonNegative(d time.Duration) string {
	if d < 0 {
		return FormatDuration(0)
	}
	return FormatDuration(d)
}
