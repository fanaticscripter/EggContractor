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
	"github.com/fanaticscripter/EggContractor/notify"
	"github.com/fanaticscripter/EggContractor/solo"
	"github.com/fanaticscripter/EggContractor/util"
)

var _notifications chan notify.Notification

var _refreshCommand = &cobra.Command{
	Use:     "refresh",
	Short:   "Refresh game state and print statuses of active solo contracts & coops",
	Args:    cobra.NoArgs,
	PreRunE: subcommandPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		notificationsDone := make(chan bool)
		_notifications = make(chan notify.Notification, 4)
		go func() {
			notify.NotificationWorker(_config.Notification, _notifications)
			notificationsDone <- true
		}()
		defer func() {
			close(_notifications)
			<-notificationsDone
		}()

		periodicalsDone := make(chan bool)
		go func() {
			_, _, err := refreshPeriodicals()
			if err != nil {
				log.Error(err)
			}
			periodicalsDone <- true
		}()
		defer func() { <-periodicalsDone }()

		playerId := _config.Player.Id
		now := time.Now()
		resp, err := api.RequestFirstContact(&api.FirstContactRequestPayload{
			PlayerId: playerId,
			X3:       1,
		})
		if err != nil {
			return err
		}
		if resp.Data == nil || resp.Data.PlayerId == "" {
			return fmt.Errorf("invalid /first_contact response for player %#v: %+v", playerId, resp)
		}

		nonFatalErrorOccurred := false

		contracts := resp.Data.AllContractProperties()
		for _, c := range contracts {
			exists, err := db.InsertContract(now, c, true /* checkExistence */)
			if err != nil {
				log.Error(err)
				nonFatalErrorOccurred = true
			} else if !exists {
				notifyNewContract(c)
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
			wrapped := coop.WrapCoopStatusWithContractList(res.status, contracts)
			activities, err := db.GetCoopMemberActivityStats(wrapped, res.timestamp)
			if err != nil {
				log.Error(err)
				nonFatalErrorOccurred = true
				wrapped.Display(_sortBy.by, nil)
			} else {
				wrapped.Display(_sortBy.by, activities)
			}
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

func refreshPeriodicals() (activeEvents []*api.Event, activeContracts []*api.ContractProperties, err error) {
	now := time.Now()
	p, err := api.RequestPeriodicals(&api.GetPeriodicalsRequestPayload{
		PlayerId:     _config.Player.Id,
		X2:           1,
		EarningBonus: 1e12, // Use a reasonably large EB just in case
	})
	if err != nil {
		return
	}
	activeEvents = p.Events.Events
	activeContracts = p.Contracts.Contracts
	seen := now
	if p.Contracts.ResponseTimestamp != 0 {
		seen = util.DoubleToTime(p.Contracts.ResponseTimestamp)
	}
	for _, e := range activeEvents {
		if err := db.InsertEvent(seen, e); err != nil {
			log.Error(err)
		}
	}
	for _, c := range activeContracts {
		exists, err := db.InsertContract(seen, c, true /* checkExistence */)
		if err != nil {
			log.Error(err)
		} else if !exists {
			notifyNewContract(c)
		}
	}
	return
}

func notifyNewContract(c *api.ContractProperties) {
	if _notifications == nil {
		// Notification system not initialized.
		return
	}
	if time.Now().After(c.ExpiryTime()) || c.Id == "first-contract" {
		return
	}
	n, err := notify.NewContractNotification(c)
	if err != nil {
		log.Error(err)
	} else {
		_notifications <- n
	}
}
