package service

import (
	"context"
	"errors"
	"time"

	"github.com/footfish/numan"
	"github.com/footfish/numan/internal/service/auth"
	"github.com/footfish/numan/internal/service/datastore"
)

// numberingService implements the NumberingService interface
type numberingService struct {
	next numan.NumberingService
	hist numan.HistoryService //used for logging
}

// NewNumberService instantiates a new NumberService.
func NewNumberingService(store *datastore.Store) numan.NumberingService {
	return &numberingService{
		next: auth.NewNumberingService(store),
		hist: NewHistoryService(store),
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

	err := s.next.Add(ctx, &newNumber) //storage
	if err == nil {                    //log history
		err = s.hist.AddHistory(ctx, numan.History{E164: newNumber.E164, Action: "added", Notes: "Domain:" + number.Domain + ", Carrier:" + number.Carrier})
	}
	return err
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

//ListOwnerID implements NumberingService.ListOwnerID()
func (s *numberingService) ListOwnerID(ctx context.Context, oid int64) ([]numan.Numbering, error) {
	return s.next.ListOwnerID(ctx, oid)
}

//Summary implements NumberingService.Summary()
func (s *numberingService) Summary(ctx context.Context) (string, error) {
	return s.next.Summary(ctx)
}

//Delete implements NumberingService.Delete()
func (s *numberingService) Delete(ctx context.Context, phonenumber *numan.E164) error {
	if phonenumber == nil {
		return errors.New("nil pointer")
	}
	err := s.next.Delete(ctx, phonenumber)
	if err == nil { //log history
		err = s.hist.AddHistory(ctx, numan.History{E164: *phonenumber, Action: "deleted"})
	}
	return err
}

//View implements NumberingService.View()
func (s *numberingService) View(ctx context.Context, number *numan.E164) (string, error) {
	if number == nil {
		return "Run time error, nil pointer", errors.New("nil pointer")
	}

	return s.next.View(ctx, number)
}

//Reserve implements NumberingService.Reserve()
func (s *numberingService) Reserve(ctx context.Context, number *numan.E164, ownerID *int64, untilTS *int64) error {
	if number == nil {
		return errors.New("nil pointer")
	}

	if *untilTS < time.Now().Unix() || *untilTS > (time.Now().Unix()+numan.MAXRESERVATIONTIME) {
		return errors.New("Can't reserve number, time out of bounds")
	}
	if err := number.ValidE164(); err != nil {
		return errors.New("Can't reserve number, " + err.Error())
	}
	if err := numan.ValidOwnerID(ownerID); err != nil {
		return errors.New("Can't reserve number, " + err.Error())
	}
	return s.next.Reserve(ctx, number, ownerID, untilTS)
}

//Allocate implements NumberingService.Allocate()
func (s *numberingService) Allocate(ctx context.Context, number *numan.E164, ownerID *int64) error {
	if number == nil || ownerID == nil {
		return errors.New("nil pointer")
	}

	if err := number.ValidE164(); err != nil {
		return errors.New("Can't allocate number, " + err.Error())
	}
	if err := numan.ValidOwnerID(ownerID); err != nil {
		return errors.New("Can't allocate number, " + err.Error())
	}

	err := s.next.Allocate(ctx, number, ownerID)
	if err == nil { //log history
		err = s.hist.AddHistory(ctx, numan.History{E164: *number, Action: "allocated", OwnerID: *ownerID})
	}
	return err
}

//DeAllocate implements NumberingService.DeAllocate()
func (s *numberingService) DeAllocate(ctx context.Context, number *numan.E164, ownerID *int64) error {
	if number == nil || ownerID == nil {
		return errors.New("nil pointer")
	}

	if err := number.ValidE164(); err != nil {
		return errors.New("Can't deallocate number, " + err.Error())
	}
	err := s.next.DeAllocate(ctx, number, ownerID)
	if err == nil { //log history
		err = s.hist.AddHistory(ctx, numan.History{E164: *number, Action: "deallocated", OwnerID: *ownerID})
	}
	return err
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
	err := s.next.Portout(ctx, number, PortoutTS)
	if err == nil { //log history
		err = s.hist.AddHistory(ctx, numan.History{E164: *number, Action: "port-out", Notes: "Scheduled: " + time.Unix(*PortoutTS, 0).Format(numan.TIMESTAMPPRINTFORMAT)})
	}
	return err
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
	err := s.next.Portin(ctx, number, PortinTS)
	if err == nil { //log history
		err = s.hist.AddHistory(ctx, numan.History{E164: *number, Action: "port-in", Notes: "Scheduled: " + time.Unix(*PortinTS, 0).Format(numan.TIMESTAMPPRINTFORMAT)})
	}
	return err
}
