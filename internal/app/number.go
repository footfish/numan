package app

import (
	"errors"
	"fmt"
	"time"

	"github.com/footfish/numan"
	"github.com/footfish/numan/internal/storage"
)

// NumberService implements the NumberAPI
type NumberService struct {
	// a database dependency would go here but instead we're going to have a static map
	db *storage.Store
	//n map[int64]Number
}

// NewNumberService instantiates a new NuService.
func NewNumberService(dsn string) *NumberService {
	return &NumberService{
		db: storage.NewStore(dsn),
	}
}

// Add implements NumberAPI.Add()
func (s *NumberService) Add(number *numan.Number) error {
	if err := number.E164.ValidE164(); err != nil {
		return err
	}
	if len(number.Domain) == 0 || len(number.Carrier) == 0 {
		return errors.New("Carrier & domain required")
	}
	newNumber := numan.Number{E164: number.E164, Domain: number.Domain, Carrier: number.Carrier} //clean
	//add to storage
	return s.db.Add(&newNumber)
}

//List implements NumberAPI.List()
func (s *NumberService) List(filter *numan.NumberFilter) ([]numan.Number, error) {
	filtered, err := s.db.List(filter)
	if err != nil {
		return filtered, err
	}
	return filtered, nil
}

//ListFree implements NumberAPI.ListFree()
func (s *NumberService) ListFree(filter *numan.NumberFilter) ([]numan.Number, error) {
	filter.State = 1
	return s.List(filter)
}

//ListUserID implements NumberAPI.ListUserID()
func (s *NumberService) ListUserID(uid int) ([]numan.Number, error) {
	filter := &numan.NumberFilter{UserID: uid}
	return s.List(filter)
}

//Summary implements NumberAPI.Summary()
func (s *NumberService) Summary() (string, error) {
	return s.db.Summary()
}

//Delete implements NumberAPI.Delete()
func (s *NumberService) Delete(phonenumber *numan.E164) error {
	return s.db.Delete(phonenumber)
}

//View implements NumberAPI.View()
func (s *NumberService) View(number *numan.E164) (string, error) {
	var view string
	if result, err := s.List(&numan.NumberFilter{E164: *number}); err != nil {
		return "", err
	} else {
		for _, r := range result {
			view += fmt.Sprintf("#%d) +%v-%v-%v, Domain: %v, Carrier: %v\n", r.ID, r.E164.Cc, r.E164.Ndc, r.E164.Sn, r.Domain, r.Carrier)
			if r.Used { //Used can be reserved or allocated
				if r.Allocated > 0 {
					view += fmt.Sprintf("Allocated to UserID: %v on %v\n", r.UserID, time.Unix(int64(r.Allocated), 0).Format(numan.DATEPRINTFORMAT))
				}
				if r.Reserved > 0 {
					view += fmt.Sprintf("Reserved %v\n", time.Unix(int64(r.Reserved), 0).Format(numan.DATEPRINTFORMAT))
				}
			} else if r.DeAllocated > 0 {
				view += fmt.Sprintf("Last allocated %v\n", time.Unix(int64(r.DeAllocated), 0).Format(numan.DATEPRINTFORMAT))
				if r.PortedOut > 0 {
					view += fmt.Sprintf("Ported out %v\n", time.Unix(int64(r.PortedOut), 0).Format(numan.DATEPRINTFORMAT))
				}
			} else {
				view += fmt.Sprintf("Never allocated\n")
			}
			if r.PortedIn > 0 {
				view += fmt.Sprintf("Ported in %v\n", time.Unix(int64(r.PortedIn), 0).Format(numan.DATEPRINTFORMAT))
			}
		}
	}

	// TODO - History
	return view, nil
}

//Reserve implements NumberAPI.Reserve()
func (s *NumberService) Reserve(number *numan.E164, userID *int, untilTS *int64) error {
	if *untilTS < time.Now().Unix() || *untilTS > (time.Now().Unix()+numan.MAXRESERVATIONTIME) {
		return errors.New("Can't reserve number, time out of bounds")
	}
	if err := number.ValidE164(); err != nil {
		return errors.New("Can't reserve number, " + err.Error())
	}
	if err := numan.ValidUserID(userID); err != nil {
		return errors.New("Can't reserve number, " + err.Error())
	}
	return s.db.Reserve(number, userID, untilTS)
}

//Allocate implements NumberAPI.Allocate()
func (s *NumberService) Allocate(number *numan.E164, userID *int) error {
	if err := number.ValidE164(); err != nil {
		return errors.New("Can't allocate number, " + err.Error())
	}
	if err := numan.ValidUserID(userID); err != nil {
		return errors.New("Can't allocate number, " + err.Error())
	}
	return s.db.Allocate(number, userID)
}

//DeAllocate implements NumberAPI.DeAllocate()
func (s *NumberService) DeAllocate(number *numan.E164) error {
	if err := number.ValidE164(); err != nil {
		return errors.New("Can't allocate number, " + err.Error())
	}
	return s.db.DeAllocate(number)
}

//Portout implements NumberAPI.Portout()
func (s *NumberService) Portout(number *numan.E164, PortoutTS *int64) error {
	if *PortoutTS < time.Now().Unix()-(365*24*60*60) || *PortoutTS > (time.Now().Unix()+(365*24*60*60)) {
		return errors.New("Can't use date, time out of bounds (+-1 year)")
	}

	if err := number.ValidE164(); err != nil {
		return errors.New("Can't, " + err.Error())
	}
	return s.db.Portout(number, PortoutTS)
}

//Portin implements NumberAPI.Portin()
func (s *NumberService) Portin(number *numan.E164, PortinTS *int64) error {
	if *PortinTS < time.Now().Unix()-(365*24*60*60) || *PortinTS > (time.Now().Unix()+(365*24*60*60)) {
		return errors.New("Can't use date, time out of bounds (+-1 year)")
	}

	if err := number.ValidE164(); err != nil {
		return errors.New("Can't, " + err.Error())
	}
	return s.db.Portin(number, PortinTS)
}

//Close closes db connection
func (s *NumberService) Close() {
	s.db.Close()
}
