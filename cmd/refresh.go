package cmd

import (
	"errors"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/fanaticscripter/EggContractor/api"
	"github.com/fanaticscripter/EggContractor/config"
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

		now := time.Now()
		nonFatalErrorOccurred := false

		saves, errored := retrieveSaves()
		if len(saves) == 0 {
			return errors.New("failed to retrieve any save, cannot proceed")
		}
		if errored {
			nonFatalErrorOccurred = true
		}

		contracts, errored := aggregateContractsFromSaves(now, saves)
		if errored {
			nonFatalErrorOccurred = true
		}

		refreshId, err := db.InsertRefresh(now)
		if err != nil {
			log.Error(err)
			nonFatalErrorOccurred = true
		}

		solos, errored := processSolosFromSaves(now, refreshId, saves)
		if errored {
			nonFatalErrorOccurred = true
		}

		coops, errored := processCoopsFromSaves(refreshId, saves, contracts)
		if errored {
			nonFatalErrorOccurred = true
		}

		if len(solos) == 0 && len(coops) == 0 {
			fmt.Println(util.MsgNoActiveContracts)
		}

		if nonFatalErrorOccurred {
			return errors.New("certain operations failed")
		}
		return nil
	},
}

func init() {
	// TODO: option to suppress status display
	_rootCmd.AddCommand(_refreshCommand)
}

func refreshPeriodicals() (activeEvents []*api.Event, activeContracts []*api.ContractProperties, err error) {
	now := time.Now()
	p, err := api.RequestPeriodicals(&api.GetPeriodicalsRequestPayload{
		UserId:   _config.Players[0].Id,
		SoulEggs: 1e12, // Use a reasonably large SE count just in case
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

func retrieveSaves() (saves []*api.FirstContact_Payload, errored bool) {
	type result struct {
		index int
		save  *api.FirstContact_Payload
	}

	players := _config.Players
	var wg sync.WaitGroup
	resultCh := make(chan result, len(players))
	for i, player := range _config.Players {
		wg.Add(1)
		go func(index int, player config.PlayerConfig) {
			defer wg.Done()
			resp, err := api.RequestFirstContact(&api.FirstContactRequestPayload{
				EiUserId: player.Id,
				DeviceId: player.DeviceId,
			})
			if err != nil {
				log.Errorf("/first_contact error for player %s: %s", player.Id, err)
				return
			}
			if resp.Data == nil || resp.Data.EiUserId == "" {
				log.Errorf("invalid /first_contact response for player %s: %+v", player.Id, resp)
				return
			}
			resultCh <- result{
				index: index,
				save:  resp.Data,
			}
		}(i, player)
	}
	wg.Wait()
	close(resultCh)

	results := make(map[int]result)
	for res := range resultCh {
		results[res.index] = res
	}
	if len(results) < len(players) {
		errored = true
	}
	for i := range _config.Players {
		res, ok := results[i]
		if ok {
			saves = append(saves, res.save)
		}
	}
	return
}

func aggregateContractsFromSaves(refreshTime time.Time, saves []*api.FirstContact_Payload) (
	contracts []*api.ContractProperties,
	errored bool,
) {
	sigs := make(map[db.ContractSignature]struct{})
	for _, save := range saves {
		for _, c := range save.AllContractProperties() {
			sig := db.GetContractSignature(c)
			_, exists := sigs[sig]
			if !exists {
				contracts = append(contracts, c)
				sigs[sig] = struct{}{}
			}
		}
	}

	existingContracts, err := db.GetContracts()
	if err != nil {
		log.Error(err)
		errored = true
		return
	}
	existingSigs := make(map[db.ContractSignature]struct{})
	for _, c := range existingContracts {
		existingSigs[db.GetContractSignature(c)] = struct{}{}
	}

	for _, c := range contracts {
		sig := db.GetContractSignature(c)
		_, exists := existingSigs[sig]
		if exists {
			continue
		}
		exists, err := db.InsertContract(refreshTime, c, true /* checkExistence */)
		if err != nil {
			log.Error(err)
			errored = true
			continue
		}
		if !exists {
			notifyNewContract(c)
		}
	}
	return
}

func processSolosFromSaves(refreshTime time.Time, refreshId int64, saves []*api.FirstContact_Payload) (
	solos []*solo.SoloContract,
	errored bool,
) {
	for _, save := range saves {
		solos = append(solos, solo.GetActiveSoloContracts(save)...)
	}
	for _, c := range solos {
		c.Display(refreshTime, _config.MultiPlayerMode())
		err := db.InsertSoloStatus(refreshTime, refreshId, c)
		if err != nil {
			log.Error(err)
			errored = true
		}
	}
	return
}

func processCoopsFromSaves(refreshId int64, saves []*api.FirstContact_Payload, contractList []*api.ContractProperties) (
	coops []*coop.CoopStatus,
	errored bool,
) {
	type coopSignature struct {
		contractId string
		code       string
	}
	savedStatuses := make([]*api.CoopStatus, 0)
	sigs := make(map[coopSignature]struct{})
	for _, save := range saves {
		// Note that although ActiveCoopStatuses already contains coop statuses,
		// they are snapshots made at the point of the save, hence usually
		// outdated. We still have to fetch a fresh copy of each one.
		for _, s := range save.Contracts.ActiveCoopStatuses {
			sig := coopSignature{
				contractId: s.ContractId,
				code:       s.Code,
			}
			if _, exists := sigs[sig]; exists {
				continue
			}
			// Sometimes outdated, completed contracts are still included in
			// ActiveCoopStatuses, possibly due to a serverside bug, so we have
			// to exclude them. See issue #9 (on my private GitLab) for an
			// example of this.
			if s.DurationUntilCollectionDeadline() < -24*time.Hour {
				continue
			}
			sigs[sig] = struct{}{}
			savedStatuses = append(savedStatuses, s)
		}
	}

	type result struct {
		index     int
		timestamp time.Time
		status    *api.CoopStatus
	}

	resultCh := make(chan result, len(savedStatuses))
	wg := sync.WaitGroup{}
	for i, s := range savedStatuses {
		wg.Add(1)
		go func(index int, s *api.CoopStatus) {
			defer wg.Done()
			now := time.Now()
			status, err := api.RequestCoopStatus(&api.CoopStatusRequestPayload{
				ContractId: s.ContractId,
				Code:       s.Code,
			})
			if err != nil {
				// One possibility here is to just use the outdated status
				// instead. But all things considered, it's probably better
				// to not report than report a misleading outdated version.
				log.Errorf("error retrieving status for coop (%s, %s): %s", s.ContractId, s.Code, err)
				return
			}
			resultCh <- result{
				index:     index,
				timestamp: now,
				status:    status,
			}
		}(i, s)
	}
	wg.Wait()
	close(resultCh)

	results := make(map[int]result)
	for res := range resultCh {
		results[res.index] = res
	}
	if len(results) < len(savedStatuses) {
		errored = true
	}

	for i := range savedStatuses {
		res := results[i]
		c := coop.WrapCoopStatusWithContractList(res.status, contractList)
		coops = append(coops, c)
		activities, err := db.GetCoopMemberActivityStats(c, res.timestamp)
		if err != nil {
			log.Error(err)
			errored = true
			c.Display(_sortBy.by, nil)
		} else {
			c.Display(_sortBy.by, activities)
		}
		if err := db.InsertCoopStatus(res.timestamp, refreshId, res.status); err != nil {
			log.Error(err)
			errored = true
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
