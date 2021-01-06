package db

import (
	"database/sql"
	"fmt"
	"sort"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"github.com/fanaticscripter/EggContractor/api"
	"github.com/fanaticscripter/EggContractor/coop"
	"github.com/fanaticscripter/EggContractor/solo"
	"github.com/fanaticscripter/EggContractor/util"
)

func InsertContract(now time.Time, c *api.ContractProperties) error {
	action := fmt.Sprintf("insert contract %s into database", c.Id)
	marshalledProps, err := proto.Marshal(c)
	if err != nil {
		return errors.Wrap(err, action)
	}
	expiryYear := c.ExpiryTime().Year()
	return transact(
		action,
		func(tx *sql.Tx) error {
			if _, err := tx.Exec(`INSERT INTO
				contract(text_id, expiry_year, coop_allowed, props, first_seen_timestamp, expiry_timestamp)
				VALUES(?, ?, ?, ?, ?, ?)
				ON CONFLICT(text_id, expiry_year)
				DO UPDATE SET
					coop_allowed = excluded.coop_allowed,
					props = excluded.props,
					expiry_timestamp = excluded.expiry_timestamp`,
				c.Id, expiryYear, c.CoopAllowed, marshalledProps, util.TimeToDouble(now), c.ExpiryTimestamp); err != nil {
				return err
			}
			return nil
		},
	)
}

// expiryYear is optional. When left unspecified, the latest iteration is
// retrieved. If no matching contract is found in the database, the return value
// is (nil, nil), i.e., error is not set in that case.
func GetContract(id string, expiryYear int) (*api.ContractProperties, error) {
	action := fmt.Sprintf("query contract %s", id)
	var foundRow bool
	var props []byte
	err := transact(
		action,
		func(tx *sql.Tx) error {
			var row *sql.Row
			if expiryYear == 0 {
				row = tx.QueryRow("SELECT props FROM contract WHERE text_id = ? ORDER BY expiry_year DESC LIMIT 1",
					id)
			} else {
				row = tx.QueryRow("SELECT props FROM contract WHERE text_id = ? AND expiry_year = ?",
					id, expiryYear)
			}
			err := row.Scan(&props)
			switch {
			case err == sql.ErrNoRows:
				// No such contract
			case err != nil:
				return err
			default:
				foundRow = true
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	if !foundRow {
		return nil, nil
	}
	contract := &api.ContractProperties{}
	err = proto.Unmarshal(props, contract)
	if err != nil {
		return nil, errors.Wrap(err, action)
	}
	return contract, nil
}

func GetContracts() ([]*api.ContractProperties, error) {
	contracts := make([]*api.ContractProperties, 0)
	action := "retrieve contracts"
	err := transact(
		action,
		func(tx *sql.Tx) error {
			rows, err := tx.Query("SELECT props FROM contract")
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var marshalledProps []byte
				if err := rows.Scan(&marshalledProps); err != nil {
					return err
				}
				contract := &api.ContractProperties{}
				if err := proto.Unmarshal(marshalledProps, contract); err != nil {
					return err
				}
				contracts = append(contracts, contract)
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	sort.SliceStable(contracts, func(i, j int) bool {
		return contracts[i].ExpiryTime().After(contracts[j].ExpiryTime())
	})
	return contracts, nil
}

func GetCoopContracts() ([]*api.ContractProperties, error) {
	contracts, err := GetContracts()
	if err != nil {
		return nil, err
	}
	coopContracts := make([]*api.ContractProperties, 0)
	for _, c := range contracts {
		if c.CoopAllowed {
			coopContracts = append(coopContracts, c)
		}
	}
	return coopContracts, nil
}

func InsertRefresh(timestamp time.Time) (refreshId int64, err error) {
	err = transact(
		"insert to table refresh",
		func(tx *sql.Tx) error {
			result, err := tx.Exec("INSERT INTO refresh(timestamp) VALUES (?)",
				util.TimeToDouble(timestamp))
			if err != nil {
				return err
			}
			refreshId, err = result.LastInsertId()
			if err != nil {
				return err
			}
			return nil
		},
	)
	if err != nil {
		return 0, err
	}
	return
}

// refreshId is optional, and will be inserted as NULL if zero.
func InsertCoopStatus(timestamp time.Time, refreshId int64, c *api.CoopStatus) error {
	action := fmt.Sprintf("insert coop status for (%s, %s) into database", c.ContractId, c.Code)
	marshalledStatus, err := proto.Marshal(c)
	if err != nil {
		return errors.Wrap(err, action)
	}
	return transact(
		action,
		func(tx *sql.Tx) error {
			var contractId int
			err := tx.QueryRow("SELECT id FROM contract WHERE text_id = ? ORDER BY expiry_year DESC LIMIT 1",
				c.ContractId).Scan(&contractId)
			switch {
			case err == sql.ErrNoRows:
				return errors.Errorf("no contract with id %#v in database", c.ContractId)
			case err != nil:
				return err
			}

			if _, err := tx.Exec("INSERT INTO coop(contract_id, code) VALUES (?, ?) ON CONFLICT DO NOTHING",
				contractId, c.Code); err != nil {
				return err
			}

			if _, err := tx.Exec(`INSERT INTO coop_status(coop_id, refresh_id, timestamp, status) VALUES (
				(SELECT id FROM coop WHERE contract_id = ? AND code = ?), ?, ?, ?)`,
				contractId, c.Code,
				toNullInt64(refreshId),
				util.TimeToDouble(timestamp),
				marshalledStatus); err != nil {
				return err
			}

			return nil
		},
	)
}

// refreshId is optional, and will be inserted as NULL if zero.
func InsertSoloStatus(timestamp time.Time, refreshId int64, c *solo.SoloContract) error {
	action := fmt.Sprintf("insert solo status for %s into database", c.GetId())
	marshalledStatus, err := c.Marshal()
	if err != nil {
		return errors.Wrap(err, action)
	}
	return transact(
		action,
		func(tx *sql.Tx) error {
			var contractId int
			err := tx.QueryRow("SELECT id FROM contract WHERE text_id = ? ORDER BY expiry_year DESC LIMIT 1",
				c.GetId()).Scan(&contractId)
			switch {
			case err == sql.ErrNoRows:
				return errors.Errorf("no contract with id %#v in database", c.GetId())
			case err != nil:
				return err
			}

			if _, err := tx.Exec("INSERT INTO solo_status(contract_id, refresh_id, timestamp, status) VALUES (?, ?, ?, ?)",
				contractId,
				toNullInt64(refreshId),
				util.TimeToDouble(timestamp),
				marshalledStatus); err != nil {
				return err
			}

			return nil
		},
	)
}

// GetSoloAndCoopStatusesFromRefresh returns contract statuses from the last
// refresh by the specified cutoff time. If no such refresh can be found in the
// database, all return values will be set to the corresponding zero values
// (including nil err).
func GetSoloAndCoopStatusesFromRefresh(byThisTime time.Time) (
	timestamp time.Time,
	solos []*solo.SoloContract,
	coops []*coop.CoopStatus,
	err error,
) {
	action := "retrieve latest coop and solo statuses"
	coops = make([]*coop.CoopStatus, 0)
	solos = make([]*solo.SoloContract, 0)
	err = transact(
		action,
		func(tx *sql.Tx) error {
			var refreshId int64
			var epochTimestamp float64
			err := tx.QueryRow(`SELECT id, timestamp FROM refresh WHERE timestamp <= ? ORDER BY timestamp DESC LIMIT 1`,
				util.TimeToDouble(byThisTime)).Scan(&refreshId, &epochTimestamp)
			switch {
			case err == sql.ErrNoRows:
				// No refresh found before the specified timestamp.
				return nil
			case err != nil:
				return err
			}
			timestamp = util.DoubleToTime(epochTimestamp)

			rows, err := tx.Query(`SELECT status FROM solo_status WHERE refresh_id = ?`, refreshId)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var marshalledSoloStatus []byte
				if err := rows.Scan(&marshalledSoloStatus); err != nil {
					return err
				}
				soloStatus, err := solo.UnmarshalSoloContract(marshalledSoloStatus)
				if err != nil {
					return err
				}
				solos = append(solos, soloStatus)
			}

			rows, err = tx.Query(`SELECT coop_status.status, contract.props FROM coop_status
				INNER JOIN coop ON coop_status.coop_id = coop.id
				INNER JOIN contract ON coop.contract_id = contract.id
				WHERE coop_status.refresh_id = ?`, refreshId)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var marshalledCoopStatus []byte
				var marshalledContractProps []byte
				if err := rows.Scan(&marshalledCoopStatus, &marshalledContractProps); err != nil {
					return err
				}
				coopStatus := &api.CoopStatus{}
				contractProps := &api.ContractProperties{}
				if err := proto.Unmarshal(marshalledCoopStatus, coopStatus); err != nil {
					return err
				}
				if err := proto.Unmarshal(marshalledContractProps, contractProps); err != nil {
					return err
				}
				coops = append(coops, &coop.CoopStatus{
					CoopStatus: coopStatus,
					Contract:   contractProps,
				})
			}

			return nil
		},
	)
	return
}

func InsertPeeked(p *Peeked) error {
	action := fmt.Sprintf("insert peeked coop (%s, %s)", p.ContractId, p.Code)
	return transact(
		action,
		func(tx *sql.Tx) error {
			_, err := tx.Exec(`INSERT INTO peeked(
					contract_id, code, last_peeked, has_completed, openings,
					eggs_laid, eggs_per_hour, required_eggs_per_hour, time_left,
					max_eb_percentage, mean_eb_percentage
				) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
				ON CONFLICT(contract_id, code) DO UPDATE SET
					last_peeked = excluded.last_peeked,
					has_completed = excluded.has_completed,
					openings = excluded.openings,
					eggs_laid = excluded.eggs_laid,
					eggs_per_hour = excluded.eggs_per_hour,
					required_eggs_per_hour = excluded.required_eggs_per_hour,
					time_left = excluded.time_left,
					max_eb_percentage = excluded.max_eb_percentage,
					mean_eb_percentage = excluded.mean_eb_percentage`,
				p.ContractId, p.Code, util.TimeToDouble(p.LastPeekedTime), p.HasCompleted, p.Openings,
				p.EggsLaid, p.EggsPerHour, p.RequiredEggsPerHour, p.TimeLeft.Seconds(),
				p.MaxEarningBonusPercentage, p.MeanEarningBonusPercentage)
			if err != nil {
				return err
			}
			return nil
		},
	)
}

// GetPeeked retrieves the list of peeked coops since the specified time.
// contractId is optional; when nonempty, only return peeked coops for that
// particular contract. Returned coops are ordered reverse chronologically.
func GetPeeked(contractId string, since time.Time) ([]*Peeked, error) {
	action := "retrieve peeked coops"
	peekedList := make([]*Peeked, 0)
	err := transact(
		action,
		func(tx *sql.Tx) error {
			baseQuery := `SELECT contract_id, code, last_peeked, has_completed, openings,
				eggs_laid, eggs_per_hour, required_eggs_per_hour, time_left,
				max_eb_percentage, mean_eb_percentage FROM peeked`
			var query string
			params := make([]interface{}, 0)
			if contractId == "" {
				query = baseQuery + ` WHERE last_peeked >= ? ORDER BY last_peeked DESC`
				params = append(params, util.TimeToDouble(since))
			} else {
				query = baseQuery + ` WHERE last_peeked >= ? AND contract_id = ? ORDER BY last_peeked DESC`
				params = append(params, util.TimeToDouble(since), contractId)
			}
			rows, err := tx.Query(query, params...)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var lastPeeked, timeLeft float64
				peeked := &Peeked{}
				if err := rows.Scan(
					&peeked.ContractId, &peeked.Code, &lastPeeked, &peeked.HasCompleted, &peeked.Openings,
					&peeked.EggsLaid, &peeked.EggsPerHour, &peeked.RequiredEggsPerHour, &timeLeft,
					&peeked.MaxEarningBonusPercentage, &peeked.MeanEarningBonusPercentage); err != nil {
					return err
				}
				peeked.LastPeekedTime = util.DoubleToTime(lastPeeked)
				peeked.TimeLeft = util.DoubleToDuration(timeLeft)
				peekedList = append(peekedList, peeked)
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return peekedList, nil
}

// GetPeekedGroupedByContract retrieves the list of peeked coops since the
// specified time, grouped by contract ID. Returned groups are ordered reverse
// chronologically (order is preserved in contractIds).
func GetPeekedGroupedByContract(since time.Time) (contractIds []string, groups map[string][]*Peeked, err error) {
	peekedList, err := GetPeeked("", since)
	if err != nil {
		return nil, nil, err
	}
	contractIds = make([]string, 0)
	groups = make(map[string][]*Peeked)
	for _, p := range peekedList {
		var group []*Peeked
		group, ok := groups[p.ContractId]
		if !ok {
			contractIds = append(contractIds, p.ContractId)
			group = make([]*Peeked, 0)
		}
		groups[p.ContractId] = append(group, p)
	}
	return contractIds, groups, nil
}

func InsertEvent(seen time.Time, e *api.Event) error {
	action := fmt.Sprintf("insert event %s (%s)", e.Id, e.Message)
	seenTimestamp := util.TimeToDouble(seen)
	expiryTimestamp := seenTimestamp + e.SecondsRemaining
	return transact(
		action,
		func(tx *sql.Tx) error {
			_, err := tx.Exec(`INSERT INTO event(
				id, event_type, multiplier, message, first_seen_timestamp, last_seen_timestamp, expiry_timestamp
			) VALUES(?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(id) DO UPDATE SET
				event_type = excluded.event_type,
				multiplier = excluded.multiplier,
				message = excluded.message,
				first_seen_timestamp = min(first_seen_timestamp, excluded.first_seen_timestamp),
				last_seen_timestamp = max(last_seen_timestamp, excluded.last_seen_timestamp),
				expiry_timestamp = excluded.expiry_timestamp`,
				e.Id, e.EventType, e.Multiplier, e.Message, seenTimestamp, seenTimestamp, expiryTimestamp)
			if err != nil {
				return err
			}
			return nil
		})
}

func GetEvents() (events []*Event, err error) {
	action := "retrieve recorded events"
	events = make([]*Event, 0)
	err = transact(
		action,
		func(tx *sql.Tx) error {
			rows, err := tx.Query(`SELECT
				id, event_type, multiplier, message, first_seen_timestamp, last_seen_timestamp, expiry_timestamp
				FROM event
				ORDER BY first_seen_timestamp DESC`)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var firstSeenTimestamp, lastSeenTimestamp, expiryTimestamp float64
				e := &Event{}
				if err := rows.Scan(&e.Id, &e.EventType, &e.Multiplier, &e.Message,
					&firstSeenTimestamp, &lastSeenTimestamp, &expiryTimestamp); err != nil {
					return err
				}
				e.FirstSeenTime = util.DoubleToTime(firstSeenTimestamp)
				e.LastSeenTime = util.DoubleToTime(lastSeenTimestamp)
				e.ExpiryTime = util.DoubleToTime(expiryTimestamp)
				events = append(events, e)
			}
			return nil
		},
	)
	return
}

func transact(description string, txFunc func(*sql.Tx) error) (err error) {
	tx, err := _db.Begin()
	if err != nil {
		return errors.Wrap(err, description)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
			err = errors.Wrap(err, description)
		} else {
			err = tx.Commit()
			if err != nil {
				err = errors.Wrap(err, description)
			}
		}
	}()
	err = txFunc(tx)
	return err
}

func toNullInt64(n int64) sql.NullInt64 {
	if n == 0 {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: n,
		Valid: true,
	}
}
