package datastore

import (
	"context"

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

//GetHistoryByNumber implements HistoryService.GetHistoryByNumber()
func (s *historyService) GetHistoryByNumber(ctx context.Context, phoneNumber numan.E164) (history []numan.History, err error) {
	return
}

//GetHistoryByOwnerID implements HistoryService.GetHistoryByUserId()
func (s *historyService) GetHistoryByOwnerID(ctx context.Context, ownerID int64) (history []numan.History, err error) {
	return
}
