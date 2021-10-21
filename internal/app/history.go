package app

import (
	"context"

	"github.com/footfish/numan"
	"github.com/footfish/numan/internal/app/datastore"
)

// historyService implements the HistoryService interface
type historyService struct {
	next numan.HistoryService
}

// NewHistoryService instantiates a new HistoryService.
func NewHistoryService(store *datastore.Store) numan.HistoryService {
	return &historyService{
		next: datastore.NewHistoryService(store),
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
