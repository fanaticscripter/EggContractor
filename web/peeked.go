package web

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/fanaticscripter/EggContractor/db"
)

const _peekedThreshold = 7 * 24 * time.Hour

type peekedPayload struct {
	Errors   []error
	Warnings []string

	PeekedContracts []*peekedContract

	Peeker *peekerPayload
}

type peekedContract struct {
	ContractId string
	Coops      []*db.Peeked
}

// GET /peeked/
func peekedHandler(c echo.Context) error {
	payload := getPeekedPayload()
	return c.Render(http.StatusOK, "peeked.html", payload)
}

func getPeekedPayload() *peekedPayload {
	errs := make([]error, 0)
	warnings := make([]string, 0)

	contractIds, groups, err := db.GetPeekedGroupedByContract(time.Now().Add(-_peekedThreshold))
	if err != nil {
		errs = append(errs, err)
		return &peekedPayload{
			Errors: errs,
			Peeker: &peekerPayload{},
		}
	}
	peekedContracts := make([]*peekedContract, 0)
	for _, contractId := range contractIds {
		peekedContracts = append(peekedContracts, &peekedContract{
			ContractId: contractId,
			Coops:      groups[contractId],
		})
	}

	contractId := ""
	if len(contractIds) > 0 {
		contractId = contractIds[0]
	}
	peeker, err := newPeekerPayload(contractId, "")
	if err != nil {
		errs = append(errs, err)
	}

	return &peekedPayload{
		Errors:          errs,
		Warnings:        warnings,
		PeekedContracts: peekedContracts,
		Peeker:          peeker,
	}
}
