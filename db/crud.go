package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"github.com/fanaticscripter/EggContractor/api"
	"github.com/fanaticscripter/EggContractor/coop"
	"github.com/fanaticscripter/EggContractor/solo"
	"github.com/fanaticscripter/EggContractor/util"
)

// ContractSignature is used to determine contract uniqueness.
type ContractSignature struct {
	Id         string
	ExpiryYear int
}

func GetContractSignature(c *api.ContractProperties) ContractSignature {
	return ContractSignature{
		Id:         c.Id,
		ExpiryYear: c.ExpiryTime().Year(),
	}
}

// InsertContract upserts the contract into the database. If checkExistence is
// true, perform an additional query beforehand to determine whether the
// contract already exists in the database, and if so, set exists to true in the
// return values. If checkExistence is false, then exists in the return values
// is meaningless.
func InsertContract(now time.Time, c *api.ContractProperties, checkExistence bool) (exists bool, err error) {
	action := fmt.Sprintf("insert contract %s into database", c.Id)
	marshalledProps, err := proto.Marshal(c)
	if err != nil {
		err = errors.Wrap(err, action)
		return
	}
	expiryYear := c.ExpiryTime().Year()
	err = transact(
		action,
		func(tx *sql.Tx) error {
			if checkExistence {
				row := tx.QueryRow(`SELECT id FROM contract WHERE text_id = ? AND expiry_year = ?`, c.Id, expiryYear)
				var rowid int64
				err := row.Scan(&rowid)
				switch {
				case err == sql.ErrNoRows:
					// Contract doesn't exist
				case err != nil:
					return err
				default:
					exists = true
				}
			}
			if _, err := tx.Exec(`INSERT INTO
				contract(text_id, expiry_year, coop_allowed, props, first_seen_timestamp, expiry_timestamp)
				VALUES(?, ?, ?, ?, ?, ?)
				ON CONFLICT(text_id, expiry_year)
				DO UPDATE SET
					coop_allowed = excluded.coop_allowed,
					props = excluded.props,
					expiry_timestamp = excluded.expiry_timestamp
				WHERE excluded.expiry_timestamp < expiry_timestamp`,
				c.Id, expiryYear, c.CoopAllowed, marshalledProps, util.TimeToDouble(now), c.ExpiryTimestamp); err != nil {
				return err
			}
			return nil
		},
	)
	return
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
			rows, err := tx.Query(`SELECT props FROM contract
				ORDER BY first_seen_timestamp DESC NULLS LAST, expiry_timestamp DESC;`)
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

func GetContractCount() (int, error) {
	var count int
	err := _db.QueryRow("SELECT COUNT(*) FROM contract").Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "retrieve contract count")
	}
	return count, nil
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

func GetCoopMemberActivityStats(c *coop.CoopStatus, refreshTime time.Time) (
	activities map[string]*coop.CoopMemberActivity,
	err error,
) {
	action := fmt.Sprintf("retrieve member last seen times for coop (%s, %s)", c.ContractId, c.Code)
	numMembers := len(c.Members)
	unaccountedForMembers := make(map[string]*api.CoopStatus_Member, numMembers)
	for _, m := range c.Members {
		unaccountedForMembers[m.Id] = m
	}
	activities = make(map[string]*coop.CoopMemberActivity, numMembers)
	err = transact(
		action,
		func(tx *sql.Tx) error {
			row := tx.QueryRow(`SELECT coop.id FROM coop
				INNER JOIN contract ON coop.contract_id = contract.id
				WHERE contract.text_id = ? AND coop.code = ?`, c.ContractId, c.Code)
			var coopId int64
			err := row.Scan(&coopId)
			switch {
			case err == sql.ErrNoRows:
				// First time seeing this coop.
				return nil
			case err != nil:
				return err
			}

			pageSize := 60
			offset := 0
			lastStatusUpdateTime := refreshTime
			for len(unaccountedForMembers) > 0 {
				rows, err := tx.Query(`SELECT timestamp, status FROM coop_status
					WHERE coop_status.coop_id = ? AND timestamp < ?
					ORDER BY timestamp DESC LIMIT ? OFFSET ?`,
					coopId, util.TimeToDouble(refreshTime), pageSize, offset)
				if err != nil {
					return err
				}
				defer rows.Close()
				numRows := 0
				for rows.Next() {
					numRows++
					var timestamp float64
					var marshalledStatus []byte
					if err := rows.Scan(&timestamp, &marshalledStatus); err != nil {
						return err
					}
					status := &api.CoopStatus{}
					if err := proto.Unmarshal(marshalledStatus, status); err != nil {
						return err
					}
					for _, m := range status.Members {
						mm, exists := unaccountedForMembers[m.Id]
						if !exists {
							continue
						}
						if mm.EggsLaid != m.EggsLaid {
							activities[m.Id] = &coop.CoopMemberActivity{
								PlayerId:         mm.Id,
								PlayerName:       mm.Name,
								LastUpdateTime:   lastStatusUpdateTime,
								OfflineTime:      refreshTime.Sub(lastStatusUpdateTime),
								EggsPerHourSince: mm.EggsPerHour(),
							}
							delete(unaccountedForMembers, m.Id)
						}
					}
					lastStatusUpdateTime = util.DoubleToTime(timestamp)
				}
				if numRows < pageSize {
					break
				}
				offset += pageSize
			}
			for _, m := range unaccountedForMembers {
				activities[m.Id] = &coop.CoopMemberActivity{
					PlayerId:           m.Id,
					PlayerName:         m.Name,
					LastUpdateTime:     lastStatusUpdateTime,
					OfflineTime:        refreshTime.Sub(lastStatusUpdateTime),
					EggsPerHourSince:   m.EggsPerHour(),
					NoActivityRecorded: true,
				}
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
