package web

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/fanaticscripter/EggContractor/db"
)

type eventsPayload struct {
	Errors   []error
	Warnings []string

	Events []*db.Event
}

// GET /events/
func eventsHandler(c echo.Context) error {
	payload := getEventsPayload()
	return c.Render(http.StatusOK, "events.html", payload)
}

func getEventsPayload() *eventsPayload {
	errs := make([]error, 0)
	warnings := make([]string, 0)

	events, err := db.GetEvents()
	if err != nil {
		errs = append(errs, err)
		return &eventsPayload{
			Errors: errs,
		}
	}
	if len(events) == 0 {
		warnings = append(warnings,
			"no event found in the database, try using the refresh subcommand of EggContractor")
	}
	return &eventsPayload{
		Errors:   errs,
		Warnings: warnings,
		Events:   events,
	}
}
