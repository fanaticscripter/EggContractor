package web

import (
	"math"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/fanaticscripter/EggContractor/db"
	"github.com/fanaticscripter/EggContractor/util"
)

type event struct {
	Id             string  `json:"id"`
	Type           string  `json:"type"`
	Multiplier     float64 `json:"multiplier"`
	Message        string  `json:"message"`
	StartTimestamp float64 `json:"startTimestamp"`
	EndTimestamp   float64 `json:"endTimestamp"`
}

// GET /api/events/
func apiEventsHandler(c echo.Context) error {
	dbEvents, err := db.GetEvents()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	events := []event{}
	for i := len(dbEvents) - 1; i >= 0; i-- {
		events = append(events, newEvent(dbEvents[i]))
	}
	return c.JSON(http.StatusOK, events)
}

func newEvent(e *db.Event) event {
	startTimestamp := util.TimeToDouble(e.FirstSeenTime)
	endTimestamp := util.TimeToDouble(e.ExpiryTime)
	// Adjust start timestamp: if the recorded duration is less that 2 minutes
	// from a whole number of hours, assume the actual duration is that whole
	// number of hours. (Events are typically 25hrs, recorded as 24:58 or 24:59
	// since I poll every two minutes.)
	durationSeconds := endTimestamp - startTimestamp
	hours, fraction := math.Modf(durationSeconds / 3600)
	if fraction*3600 > 3480 {
		durationSeconds = (hours + 1) * 3600
		startTimestamp = endTimestamp - durationSeconds
	}
	return event{
		Id:             e.Id,
		Type:           e.EventType,
		Multiplier:     e.Multiplier,
		Message:        e.Message,
		StartTimestamp: startTimestamp,
		EndTimestamp:   endTimestamp,
	}
}
