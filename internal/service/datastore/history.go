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
func (s *historyService) ListHistoryByNumber(ctx context.Context, phoneNumber numan.E164) (history []numan.History, err error) {
	return
}

//ListHistoryByOwnerID implements HistoryService.ListHistoryByUserId()
func (s *historyService) ListHistoryByOwnerID(ctx context.Context, ownerID int64) (history []numan.History, err error) {
	return
}
