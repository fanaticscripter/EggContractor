package web

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/fanaticscripter/EggContractor/api"
	"github.com/fanaticscripter/EggContractor/coop"
	"github.com/fanaticscripter/EggContractor/db"
)

var _peekedCh chan *db.Peeked

type peekPayload struct {
	Errors   []error
	Warnings []string

	RetrievalTime time.Time
	Coop          *coop.CoopStatus

	Peeker *peekerPayload
}

func init() {
	_peekedCh = make(chan *db.Peeked, 4)
}

// GET /peek/:contractId/:code/
func peekHandler(c echo.Context) error {
	contractId := strings.ToLower(c.Param("contractId"))
	code := strings.ToLower(c.Param("code"))
	payload := getPeekPayload(contractId, code)
	return c.Render(http.StatusOK, "peek.html", payload)
}

func getPeekPayload(contractId string, code string) *peekPayload {
	errs := make([]error, 0)
	warnings := make([]string, 0)
	now := time.Now()
	status, err := api.RequestCoopStatus(&api.CoopStatusRequestPayload{
		ContractId: contractId,
		Code:       code,
	})
	if err != nil {
		errs = append(errs, errors.Wrapf(
			err, "retrieve coop status for %#v contract with code %#v",
			contractId, code))
		return &peekPayload{
			Errors: errs,
			Peeker: &peekerPayload{},
		}
	}
	wrapped, err := db.WrapCoopStatusWithDB(status)
	if err != nil {
		errs = append(errs, err)
	} else if wrapped.Contract == nil {
		warnings = append(warnings, fmt.Sprintf(
			"contract %s not found in database, try using the refresh subcommand of EggContractor to populate the contract table",
			contractId))
	}

	peeker, err := newPeekerPayload(contractId, code)
	if err != nil {
		errs = append(errs, err)
	}

	_peekedCh <- db.NewPeeked(wrapped, now)

	return &peekPayload{
		Errors:        errs,
		Warnings:      warnings,
		RetrievalTime: now,
		Coop:          wrapped,
		Peeker:        peeker,
	}
}

func dbPeekedWorker() {
	for peeked := range _peekedCh {
		if err := db.InsertPeeked(peeked); err != nil {
			log.Error(err)
		}
	}
}
