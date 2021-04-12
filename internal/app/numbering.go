package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/footfish/numan"
	"github.com/footfish/numan/internal/datastore"
)

// numberingService implements the NumberingService interface
type numberingService struct {
	next numan.NumberingService
	auth numan.User
}

// NewNumberService instantiates a new NumberService.
func NewNumberingService(store *datastore.Store) numan.NumberingService {
	return &numberingService{
		next: datastore.NewNumberingService(store),
	}
}

// Add implements NumberingService.Add()
func (s *numberingService) Add(ctx context.Context, number *numan.Numbering) error {
	if number == nil {
		return errors.New("nil pointer")
	}
	if err := number.E164.ValidE164(); err != nil {
		return err
	}
	if len(number.Domain) == 0 || len(number.Carrier) == 0 {
		return errors.New("Carrier & domain required")
	}
	newNumber := numan.Numbering{E164: number.E164, Domain: number.Domain, Carrier: number.Carrier} //clean
	//add to storage
	return s.next.Add(ctx, &newNumber)
}

//AddGroup not implemented
func (s *numberingService) AddGroup(ctx context.Context) {
}

//List implements NumberingService.List()
func (s *numberingService) List(ctx context.Context, filter *numan.NumberFilter) ([]numan.Numbering, error) {
	if filter == nil {
		return nil, errors.New("nil pointer")
	}
	filtered, err := s.next.List(ctx, filter)
	if err != nil {
		return filtered, err
	}
	return filtered, nil
}

//ListUserID implements NumberingService.ListUserID()
func (s *numberingService) ListUserID(ctx context.Context, uid int64) ([]numan.Numbering, error) {
	return s.next.ListUserID(ctx, uid)
}

//Summary implements NumberingService.Summary()
func (s *numberingService) Summary(ctx context.Context) (string, error) {
	//Auth
	if err := s.auth.SetUserFromToken(fmt.Sprintf("%s", ctx.Value("token"))); err != nil { //Get authenticated user data from token
		return "Unexpected Auth error", err
	}
	/*	if s.auth.Role != "admin" {
		return "Insufficient privileges", errors.New("Insufficient privileges")
	} */

	return s.next.Summary(ctx)
}

//Delete implements NumberingService.Delete()
func (s *numberingService) Delete(ctx context.Context, phonenumber *numan.E164) error {
	if phonenumber == nil {
		return errors.New("nil pointer")
	}
	return s.next.Delete(ctx, phonenumber)
}

//View implements NumberingService.View()
func (s *numberingService) View(ctx context.Context, number *numan.E164) (string, error) {
	if number == nil {
		return "Run time error, nil pointer", errors.New("nil pointer")
	}

	return s.next.View(ctx, number)
}

//Reserve implements NumberingService.Reserve()
func (s *numberingService) Reserve(ctx context.Context, number *numan.E164, userID *int64, untilTS *int64) error {
	if number == nil {
		return errors.New("nil pointer")
	}

	if *untilTS < time.Now().Unix() || *untilTS > (time.Now().Unix()+numan.MAXRESERVATIONTIME) {
		return errors.New("Can't reserve number, time out of bounds")
	}
	if err := number.ValidE164(); err != nil {
		return errors.New("Can't reserve number, " + err.Error())
	}
	if err := numan.ValidUserID(userID); err != nil {
		return errors.New("Can't reserve number, " + err.Error())
	}
	return s.next.Reserve(ctx, number, userID, untilTS)
}

//Allocate implements NumberingService.Allocate()
func (s *numberingService) Allocate(ctx context.Context, number *numan.E164, userID *int64) error {
	if number == nil || userID == nil {
		return errors.New("nil pointer")
	}

	if err := number.ValidE164(); err != nil {
		return errors.New("Can't allocate number, " + err.Error())
	}
	if err := numan.ValidUserID(userID); err != nil {
		return errors.New("Can't allocate number, " + err.Error())
	}
	return s.next.Allocate(ctx, number, userID)
}

//DeAllocate implements NumberingService.DeAllocate()
func (s *numberingService) DeAllocate(ctx context.Context, number *numan.E164) error {
	if number == nil {
		return errors.New("nil pointer")
	}

	if err := number.ValidE164(); err != nil {
		return errors.New("Can't allocate number, " + err.Error())
	}
	return s.next.DeAllocate(ctx, number)
}

//Portout implements NumberingService.Portout()
func (s *numberingService) Portout(ctx context.Context, number *numan.E164, PortoutTS *int64) error {
	if number == nil || PortoutTS == nil {
		return errors.New("nil pointer")
	}

	if *PortoutTS < time.Now().Unix()-(365*24*60*60) || *PortoutTS > (time.Now().Unix()+(365*24*60*60)) {
		return errors.New("Can't use date, time out of bounds (+-1 year)")
	}

	if err := number.ValidE164(); err != nil {
		return errors.New("Can't, " + err.Error())
	}
	return s.next.Portout(ctx, number, PortoutTS)
}

//Portin implements NumberingService.Portin()
func (s *numberingService) Portin(ctx context.Context, number *numan.E164, PortinTS *int64) error {
	if number == nil || PortinTS == nil {
		return errors.New("nil pointer")
	}

	if *PortinTS < time.Now().Unix()-(365*24*60*60) || *PortinTS > (time.Now().Unix()+(365*24*60*60)) {
		return errors.New("Can't use date, time out of bounds (+-1 year)")
	}

	if err := number.ValidE164(); err != nil {
		return errors.New("Can't, " + err.Error())
	}
	return s.next.Portin(ctx, number, PortinTS)
}
