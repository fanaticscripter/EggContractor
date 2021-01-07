package coop

import "time"

type CoopMemberActivity struct {
	PlayerId         string
	PlayerName       string
	LastUpdateTime   time.Time
	OfflineTime      time.Duration
	EggsPerHourSince float64
	// NoActivityRecorded indicates whether any activity has been recorded for a
	// player at all; the player could have stayed at the exact same eggs laid
	// since we start recording a coop, in which case NoActivityRecorded would
	// be true.
	NoActivityRecorded bool
}
