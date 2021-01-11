package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/fanaticscripter/EggContractor/api"
	"github.com/fanaticscripter/EggContractor/config"
	"github.com/fanaticscripter/EggContractor/db"
)

type localConfig struct {
	Database   config.DatabaseConfig
	Aggregator struct {
		PlayerIdBlacklist []string `mapstructure:"player_id_blacklist"`
		KnownContractIds  []string `mapstructure:"known_contract_ids"`
	}
}

var _config *localConfig

func init() {
	log.SetLevel(log.InfoLevel)
}

func main() {
	var cfgFile string

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: %s <identifier> ...

Scrapes player info to gather contract properties. Players are either directly
specified or scraped from coops.

Each identifier either identifies a player or a coop. A player is identified by
player ID; a coop is identified as contract-id@coop-code.

Config file (defaults to config.toml):

  [database]
  # Path for contracts database. Required.
  path = "data/contracts.db"

  [aggregator]
  # A list of user IDs blacklisted for unreliable contract information, e.g. for
  # cheating on contract expiry timestamps.
  player_id_blacklist = []

  # A complete list of know contract IDs. Leggacy runs SHOULD be duplicated.
  # Used for cross-verification with a third party, e.g.
  # https://docs.google.com/spreadsheets/d/1JE5OlFG7tKfj-yXw-rN9fhLzvdQv0YmB_olrYFAAeKo/
  # from @mikit on Egg, Inc. Discord.
  known_contract_ids = []

Flags:
`, os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVar(&cfgFile, "config", "config.toml", "config file")
	flag.Parse()
	identifiers := flag.Args()

	err := loadConfig(cfgFile)
	if err != nil {
		log.Fatalf("error loading %s: %s", cfgFile, err)
	}

	if _config.Database.Path == "" {
		log.Fatalf("%s: database.path required", cfgFile)
	}
	err = db.InitDB(_config.Database)
	if err != nil {
		log.Fatal(err)
	}

	beforeCount, err := db.GetContractCount()
	if err != nil {
		log.Error(err)
		beforeCount = -1
	}

	playerIdSet := make(map[string]struct{})

	coopRe := regexp.MustCompile(`^([\w-]+)@([\w-]+)`)
	for _, identifier := range identifiers {
		if strings.Contains(identifier, "@") {
			// identifier identifies a coop
			m := coopRe.FindStringSubmatch(identifier)
			if m != nil {
				contractId := m[1]
				coopCode := m[2]
				log.Infof("scraping coop %s#%s for player IDs", contractId, coopCode)
				for _, id := range getPlayerIdsFromCoop(contractId, coopCode) {
					playerIdSet[id] = struct{}{}
				}
			} else {
				log.Errorf("unrecognized coop identifier %#v (should be of the form contract-id@coop-code)", identifier)
			}
		} else {
			// identifier is a player ID
			playerIdSet[identifier] = struct{}{}
		}
	}

	playerIds := make([]string, 0)
LoopPlayerIdSet:
	for id := range playerIdSet {
		for _, bid := range _config.Aggregator.PlayerIdBlacklist {
			if id == bid {
				continue LoopPlayerIdSet
			}
		}
		playerIds = append(playerIds, id)
	}
	sort.Strings(playerIds)

	for i, playerId := range playerIds {
		log.Infof("scraping contracts of player [%d/%d] %s", i+1, len(playerIds), playerId)
		getAndRecordPlayerContracts(playerId)
	}

	afterCount, err := db.GetContractCount()
	if err != nil {
		log.Error(err)
	} else {
		if beforeCount >= 0 {
			log.Infof("%d contracts currently in database, an increase of %d", afterCount, afterCount-beforeCount)
		} else {
			log.Infof("%d contracts currently in database", afterCount)
		}
	}

	printStillMissingContracts()
}

func loadConfig(cfgFile string) error {
	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	_config = &localConfig{}
	if err := viper.Unmarshal(_config); err != nil {
		return err
	}
	return nil
}

func getPlayerIdsFromCoop(contractId, code string) []string {
	status, err := api.RequestCoopStatus(&api.CoopStatusRequestPayload{
		ContractId: contractId,
		Code:       code,
	})
	if err != nil {
		log.Errorf("error scraping coop %s#%s: %s", contractId, code, err)
		return nil
	}
	playerIds := make([]string, 0)
	for _, m := range status.Members {
		playerIds = append(playerIds, m.Id)
	}
	return playerIds
}

func getAndRecordPlayerContracts(playerId string) {
	now := time.Now()
	resp, err := api.RequestFirstContact(&api.FirstContactRequestPayload{
		PlayerId: playerId,
		X3:       1,
	})
	if err != nil {
		log.Error(err)
		return
	}
	if resp.Data == nil || resp.Data.PlayerId == "" {
		log.Errorf("invalid response for player %#v", playerId)
		return
	}
	contracts := resp.Data.AllContractProperties()
	for _, c := range contracts {
		exists, err := db.InsertContract(now, c, true /* checkExistence */)
		if err != nil {
			log.Error(err)
		} else if !exists {
			log.Infof("inserted new contract \"%s\" (%s)", c.Name, c.Id)
		}
	}
}

func printStillMissingContracts() {
	known := make(map[string]int)
	for _, id := range _config.Aggregator.KnownContractIds {
		known[id]++
	}
	have := make(map[string]int)
	contracts, err := db.GetContracts()
	if err != nil {
		log.Error(err)
		return
	}
	for _, c := range contracts {
		have[c.Id]++
	}
	missing := make(map[string]int)
	numMissing := 0
	for id, knownCnt := range known {
		haveCnt := have[id]
		if knownCnt > haveCnt {
			missing[id] = knownCnt - haveCnt
			numMissing += knownCnt - haveCnt
		} else if knownCnt < haveCnt {
			log.Warnf("contract %s: known %d instances, have %d instances", id, knownCnt, haveCnt)
		}
	}
	for id, haveCnt := range have {
		_, ok := known[id]
		if !ok {
			log.Warnf("contract %s: known 0 instances, have %d instances", id, haveCnt)
		}
	}
	if numMissing > 0 {
		log.Infof("%d known contracts still missing", numMissing)
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		for id, missingCnt := range missing {
			knownCnt := known[id]
			fmt.Fprintf(w, "%s:\t%d/%d\n", id, missingCnt, knownCnt)
		}
		w.Flush()
	} else {
		log.Infof("all known contracts are in the database")
	}
}
