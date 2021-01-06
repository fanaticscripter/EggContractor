package cmd

import (
	"errors"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/fanaticscripter/EggContractor/api"
	"github.com/fanaticscripter/EggContractor/coop"
	"github.com/fanaticscripter/EggContractor/db"
	"github.com/fanaticscripter/EggContractor/solo"
	"github.com/fanaticscripter/EggContractor/util"
)

var _refreshCommand = &cobra.Command{
	Use:     "refresh",
	Short:   "Refresh game state and print statuses of active solo contracts & coops",
	Args:    cobra.NoArgs,
	PreRunE: subcommandPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		eventsDone := make(chan bool)
		go func() {
			_, err := refreshEvents()
			if err != nil {
				log.Error(err)
			}
			eventsDone <- true
		}()
		defer func() { <-eventsDone }()

		playerId := _config.Player.Id
		now := time.Now()
		resp, err := api.RequestFirstContact(&api.FirstContactRequestPayload{
			PlayerId: playerId,
			X3:       1,
		})
		if err != nil {
			return err
		}

		nonFatalErrorOccurred := false

		contracts := resp.Data.AllContractProperties()
		for _, c := range contracts {
			err = db.InsertContract(now, c)
			if err != nil {
				log.Error(err)
				nonFatalErrorOccurred = true
			}
		}

		refreshId, err := db.InsertRefresh(now)
		if err != nil {
			log.Error(err)
			nonFatalErrorOccurred = true
		}

		solos := solo.GetActiveSoloContracts(resp)
		unfilteredCoopStatuses := resp.Data.Contracts.ActiveCoopStatuses

		for _, c := range solos {
			c.Display()
			err := db.InsertSoloStatus(now, refreshId, c)
			if err != nil {
				log.Error(err)
				nonFatalErrorOccurred = true
			}
		}

		// Sometimes outdated, completed contracts are still included in
		// ActiveCoopStatuses, possibly due to a serverside bug, so we have to
		// exclude them. See issue #9 (on my private GitLab) for an example of
		// this.
		coopStatuses := make([]*api.CoopStatus, 0)
		for _, c := range unfilteredCoopStatuses {
			if c.DurationUntilCollectionDeadline() < -24*time.Hour {
				// At least one day past collection deadline.
				continue
			}
			coopStatuses = append(coopStatuses, c)
		}

		if len(solos) == 0 && len(coopStatuses) == 0 {
			fmt.Println(util.MsgNoActiveContracts)
		}

		type result struct {
			index     int
			timestamp time.Time
			status    *api.CoopStatus
			err       error
		}

		resultCh := make(chan result, len(coopStatuses))
		done := make(chan struct{})
		go func() {
			defer close(done)
			wg := sync.WaitGroup{}
			for i, coopStatus := range coopStatuses {
				wg.Add(1)
				go func(index int, c *api.CoopStatus) {
					defer wg.Done()
					now := time.Now()
					status, err := api.RequestCoopStatus(&api.CoopStatusRequestPayload{
						ContractId: c.ContractId,
						Code:       c.Code,
						PlayerId:   playerId,
					})
					resultCh <- result{
						index:     index,
						timestamp: now,
						status:    status,
						err:       err,
					}
				}(i, coopStatus)
			}
			wg.Wait()
		}()

		results := make(map[int]result)
	CollectResults:
		for {
			select {
			case res := <-resultCh:
				results[res.index] = res
			case <-done:
				break CollectResults
			}
		}

		for i := range coopStatuses {
			res := results[i]
			if res.err != nil {
				log.Error(res.err)
				nonFatalErrorOccurred = true
				continue
			}
			coop.WrapCoopStatusWithContractList(res.status, contracts).Display(_sortBy.by)
			if err := db.InsertCoopStatus(res.timestamp, refreshId, res.status); err != nil {
				log.Error(err)
				nonFatalErrorOccurred = true
			}
		}

		if nonFatalErrorOccurred {
			return errors.New("certain operations failed")
		}
		return nil
	},
}

func init() {
	_rootCmd.AddCommand(_refreshCommand)
}
