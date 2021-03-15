package app

import (
	"errors"
	"time"

	"github.com/footfish/numan"
)

// Add implements NumberAPI.Add()
func (s *numanService) Add(number *numan.Number) error {
	if number == nil {
		return errors.New("nil pointer")
	}
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

//AddGroup not implemented
func (s *numanService) AddGroup() {
}

//List implements NumberAPI.List()
func (s *numanService) List(filter *numan.NumberFilter) ([]numan.Number, error) {
	if filter == nil {
		return nil, errors.New("nil pointer")
	}
	filtered, err := s.db.List(filter)
	if err != nil {
		return filtered, err
	}
	return filtered, nil
}

//ListUserID implements NumberAPI.ListUserID()
func (s *numanService) ListUserID(uid int64) ([]numan.Number, error) {
	return s.db.ListUserID(uid)
}

//Summary implements NumberAPI.Summary()
func (s *numanService) Summary() (string, error) {
	if s == nil {
		return "Run time error, nil pointer", errors.New("nil pointer")
	}

	return s.db.Summary()
}

//Delete implements NumberAPI.Delete()
func (s *numanService) Delete(phonenumber *numan.E164) error {
	if phonenumber == nil {
		return errors.New("nil pointer")
	}
	return s.db.Delete(phonenumber)
}

//View implements NumberAPI.View()
func (s *numanService) View(number *numan.E164) (string, error) {
	if number == nil {
		return "Run time error, nil pointer", errors.New("nil pointer")
	}

	return s.db.View(number)
}

//Reserve implements NumberAPI.Reserve()
func (s *numanService) Reserve(number *numan.E164, userID *int64, untilTS *int64) error {
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
	return s.db.Reserve(number, userID, untilTS)
}

//Allocate implements NumberAPI.Allocate()
func (s *numanService) Allocate(number *numan.E164, userID *int64) error {
	if number == nil || userID == nil {
		return errors.New("nil pointer")
	}

	if err := number.ValidE164(); err != nil {
		return errors.New("Can't allocate number, " + err.Error())
	}
	if err := numan.ValidUserID(userID); err != nil {
		return errors.New("Can't allocate number, " + err.Error())
	}
	return s.db.Allocate(number, userID)
}

//DeAllocate implements NumberAPI.DeAllocate()
func (s *numanService) DeAllocate(number *numan.E164) error {
	if number == nil {
		return errors.New("nil pointer")
	}

	if err := number.ValidE164(); err != nil {
		return errors.New("Can't allocate number, " + err.Error())
	}
	return s.db.DeAllocate(number)
}

//Portout implements NumberAPI.Portout()
func (s *numanService) Portout(number *numan.E164, PortoutTS *int64) error {
	if number == nil || PortoutTS == nil {
		return errors.New("nil pointer")
	}

	if *PortoutTS < time.Now().Unix()-(365*24*60*60) || *PortoutTS > (time.Now().Unix()+(365*24*60*60)) {
		return errors.New("Can't use date, time out of bounds (+-1 year)")
	}

	if err := number.ValidE164(); err != nil {
		return errors.New("Can't, " + err.Error())
	}
	return s.db.Portout(number, PortoutTS)
}

//Portin implements NumberAPI.Portin()
func (s *numanService) Portin(number *numan.E164, PortinTS *int64) error {
	if number == nil || PortinTS == nil {
		return errors.New("nil pointer")
	}

	if *PortinTS < time.Now().Unix()-(365*24*60*60) || *PortinTS > (time.Now().Unix()+(365*24*60*60)) {
		return errors.New("Can't use date, time out of bounds (+-1 year)")
	}

	if err := number.ValidE164(); err != nil {
		return errors.New("Can't, " + err.Error())
	}
	return s.db.Portin(number, PortinTS)
}
