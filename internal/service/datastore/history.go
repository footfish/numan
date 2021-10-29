package datastore

import (
	"context"
	"errors"
	"time"

	"github.com/footfish/numan"
)

// historyService implements the HistoryService interface
type historyService struct {
	store Store
}

// NewHistoryService instantiates a new HistoryService.
func NewHistoryService(store *Store) numan.HistoryService {
	return &historyService{
		store: *store,
	}
}

//AddHistory  implements HistoryService.AddHistory()
func (s *historyService) AddHistory(ctx context.Context, historyEntry numan.History) error {
	_, err := s.store.db.Exec("INSERT INTO history( cc, ndc, sn, action, timestamp, ownerID, notes) values(?,?,?,?,?,?,?)", historyEntry.E164.Cc, historyEntry.E164.Ndc, historyEntry.E164.Sn, historyEntry.Action, time.Now().Unix(), historyEntry.OwnerID, historyEntry.Notes)
	if err != nil {
		err = errors.New("could not record " + historyEntry.Action + " in history")
	}
	return err
}

//ListHistoryByNumber implements HistoryService.ListHistoryByNumber()
func (s *historyService) ListHistoryByNumber(ctx context.Context, phoneNumber numan.E164) ([]numan.History, error) {
	if phoneNumber.ValidE164() != nil {
		return nil, errors.New("Incorrect number format")
	}

	var result numan.History
	var resultList []numan.History

	rows, err := s.store.db.Query("SELECT timestamp, cc, ndc, sn, ownerID, action, ifnull(notes,'') FROM history where cc=? and ndc=? and sn=? order by timestamp asc", phoneNumber.Cc, phoneNumber.Ndc, phoneNumber.Sn)
	if err != nil {
		return resultList, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&result.Timestamp,
			&result.E164.Cc,
			&result.E164.Ndc,
			&result.E164.Sn,
			&result.OwnerID,
			&result.Action,
			&result.Notes,
		)
		if err != nil {
			return resultList, err
		}
		resultList = append(resultList, result)
	}
	err = rows.Err()
	if err != nil {
		return resultList, err
	}
	return resultList, nil
}

//ListHistoryByOwnerID implements HistoryService.ListHistoryByUserId()
func (s *historyService) ListHistoryByOwnerID(ctx context.Context, ownerID int64) ([]numan.History, error) {
	if numan.ValidOwnerID(&ownerID) != nil {
		return nil, errors.New("Incorrect Owner ID format")
	}

	var result numan.History
	var resultList []numan.History

	rows, err := s.store.db.Query("SELECT timestamp, cc, ndc, sn, ownerID, action, ifnull(notes,'') FROM history where ownerID=? order by timestamp asc", ownerID)
	if err != nil {
		return resultList, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&result.Timestamp,
			&result.E164.Cc,
			&result.E164.Ndc,
			&result.E164.Sn,
			&result.OwnerID,
			&result.Action,
			&result.Notes,
		)
		if err != nil {
			return resultList, err
		}
		resultList = append(resultList, result)
	}
	err = rows.Err()
	if err != nil {
		return resultList, err
	}
	return resultList, nil
}
