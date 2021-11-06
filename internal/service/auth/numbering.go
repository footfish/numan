package auth

import (
	"context"

	"github.com/footfish/numan"
	"github.com/footfish/numan/internal/service/datastore"
)

// numberingService implements the NumberingService interface
type numberingService struct {
	next numan.NumberingService
}

// NewNumberService instantiates a new NumberService.
func NewNumberingService(store *datastore.Store) numan.NumberingService {
	return &numberingService{
		next: datastore.NewNumberingService(store),
	}
}

// Add implements NumberingService.Add()
func (s *numberingService) Add(ctx context.Context, number *numan.Numbering) error {
	if err := checkUserRole(numan.RoleUser, ctx); err != nil {
		return err
	}
	return s.next.Add(ctx, number) //storage
}

//AddGroup not implemented
func (s *numberingService) AddGroup(ctx context.Context) {
}

//List implements NumberingService.List()
func (s *numberingService) List(ctx context.Context, filter *numan.NumberFilter) ([]numan.Numbering, error) {
	if err := checkUserRole(numan.RoleUser, ctx); err != nil {
		return []numan.Numbering{}, err
	}
	return s.next.List(ctx, filter)
}

//ListOwnerID implements NumberingService.ListOwnerID()
func (s *numberingService) ListOwnerID(ctx context.Context, oid int64) ([]numan.Numbering, error) {
	if err := checkUserRole(numan.RoleUser, ctx); err != nil {
		return []numan.Numbering{}, err
	}
	return s.next.ListOwnerID(ctx, oid)
}

//Summary implements NumberingService.Summary()
func (s *numberingService) Summary(ctx context.Context) (string, error) {
	if err := checkUserRole(numan.RoleUser, ctx); err != nil {
		return err.Error(), err
	}
	return s.next.Summary(ctx)
}

//Delete implements NumberingService.Delete()
func (s *numberingService) Delete(ctx context.Context, phonenumber *numan.E164) error {
	if err := checkUserRole(numan.RoleUser, ctx); err != nil {
		return err
	}
	return s.next.Delete(ctx, phonenumber)
}

//View implements NumberingService.View()
func (s *numberingService) View(ctx context.Context, number *numan.E164) (string, error) {
	if err := checkUserRole(numan.RoleUser, ctx); err != nil {
		return err.Error(), err
	}
	return s.next.View(ctx, number)
}

//Reserve implements NumberingService.Reserve()
func (s *numberingService) Reserve(ctx context.Context, number *numan.E164, ownerID *int64, untilTS *int64) error {
	if err := checkUserRole(numan.RoleUser, ctx); err != nil {
		return err
	}
	return s.next.Reserve(ctx, number, ownerID, untilTS)
}

//Allocate implements NumberingService.Allocate()
func (s *numberingService) Allocate(ctx context.Context, number *numan.E164, ownerID *int64) error {
	if err := checkUserRole(numan.RoleUser, ctx); err != nil {
		return err
	}
	return s.next.Allocate(ctx, number, ownerID)
}

//DeAllocate implements NumberingService.DeAllocate()
func (s *numberingService) DeAllocate(ctx context.Context, number *numan.E164, ownerID *int64) error {
	if err := checkUserRole(numan.RoleUser, ctx); err != nil {
		return err
	}
	return s.next.DeAllocate(ctx, number, ownerID)
}

//Portout implements NumberingService.Portout()
func (s *numberingService) Portout(ctx context.Context, number *numan.E164, PortoutTS *int64) error {
	if err := checkUserRole(numan.RoleUser, ctx); err != nil {
		return err
	}
	return s.next.Portout(ctx, number, PortoutTS)
}

//Portin implements NumberingService.Portin()
func (s *numberingService) Portin(ctx context.Context, number *numan.E164, PortinTS *int64) error {
	if err := checkUserRole(numan.RoleUser, ctx); err != nil {
		return err
	}
	return s.next.Portin(ctx, number, PortinTS)
}
