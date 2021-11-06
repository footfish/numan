package auth

import (
	"context"

	"github.com/footfish/numan"
	"github.com/footfish/numan/internal/service/datastore"
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

//AddHistory  implements HistoryService.AddHistory()
func (s *historyService) AddHistory(ctx context.Context, historyEntry numan.History) error {
	if err := checkUserRole(numan.RoleUser, ctx); err != nil {
		return err
	}
	return s.next.AddHistory(ctx, historyEntry)
}

//ListHistoryByNumber implements HistoryService.ListHistoryByNumber()
func (s *historyService) ListHistoryByNumber(ctx context.Context, phoneNumber numan.E164) (history []numan.History, err error) {
	if err := checkUserRole(numan.RoleUser, ctx); err != nil {
		return history, err
	}
	return s.next.ListHistoryByNumber(ctx, phoneNumber)
}

//ListHistoryByOwnerID implements HistoryService.ListHistoryByUserId()
func (s *historyService) ListHistoryByOwnerID(ctx context.Context, ownerID int64) (history []numan.History, err error) {
	if err := checkUserRole(numan.RoleUser, ctx); err != nil {
		return history, err
	}
	return s.next.ListHistoryByOwnerID(ctx, ownerID)
}
