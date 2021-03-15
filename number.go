package numan

import (
	"errors"
	"regexp"
)

const (
	// MAXRESERVATIONTIME the maximum time a number can be reserved (pre-allocated)
	MAXRESERVATIONTIME = 108000 //in secs (30 mins)
	//DATEPRINTFORMAT the format used for display of timestamps
	DATEPRINTFORMAT = "02/01/2006"
	//QUARANTINE period in seconds. Numbers can't be reserved/allocated during quarantie.
	QUARANTINE = 13 * 31 * 24 * 60 * 60 //  (13 months approx)
)

//Number represents a stored phone number entry
type Number struct {
	ID          int64  // number entry index
	E164        E164   //an e.164 number
	Used        bool   // indicates number is actively being used or not (reserved, allocated).
	Domain      string // which domain is using the number (which domain can allocate)
	Carrier     string // who is the block owner
	UserID      int64  // which client a/c is using
	Allocated   int64  // timestamp of when the number was allocated OR 0 if unused
	Reserved    int64  // timestamp if the number is reserved OR 0
	DeAllocated int64  // timestamp when number was last cancelled (use for quarantine) OR 0
	PortedIn    int64  // timestamp number was ported in OR 0
	PortedOut   int64  // timestamp number was ported out  OR 0
}

//E164 represents a  phone number in e164 format
type E164 struct {
	Sn  string // subscriber number (5 to 13 digits)
	Ndc string // network code, with leading zero or '1' (2 to 5 digits)
	Cc  string // country code (1 to 3 digits), no leading zero
}

//NumberFilter represents a stored phone number lookup filter
type NumberFilter struct {
	ID          int64  // number entry index (0 unused)
	E164        E164   // an e.164 number
	State       byte   // 0 - ignore, 1 - free, 2 - used
	Domain      string // which domain is using the number (which domain can allocate)
	Carrier     string // who is the block owner
	UserID      int64  // which client a/c is using
	Allocated   bool   // if the number was ordered
	Reserved    bool   // if the number is reserved
	DeAllocated bool   // if number was last cancelled (use for quarantine)
	PortedIn    bool   // if number was ported in
	PortedOut   bool   // if number was ported out
}

//NumberAPI exposes interface for managing numbers
type NumberAPI interface {
	//Adds a new unused number to database.
	//params E164, Domain & Carrier must be included, others (supplied or not) are initialised.
	Add(number *Number) error
	//AddGroup adds a series of new unused numbers
	AddGroup()
	//List returns a filtered list of numbers
	List(filter *NumberFilter) ([]Number, error)
	//ListUserID gets list of numbers attached to specific UserID
	ListUserID(userID int64) ([]Number, error)
	//Reserve locks a number to a UserID until untilTS (unix timestamp)
	Reserve(number *E164, userID *int64, untilTS *int64) error
	//Allocate marks a number 'used' by a User
	Allocate(number *E164, userID *int64) error
	//DeAllocate number from User (number goes to quarantine)
	DeAllocate(number *E164) error
	//Portout sets a port out date (just a log, doesn't care about state or do anything else)
	Portout(number *E164, PortoutTS *int64) error
	//Portin sets a port in date (just a log, doesn't care about state or do anything else)
	Portin(number *E164, PortinTS *int64) error
	//Delete - number no longer used, removed from number db, must be unused (history kept).
	Delete(number *E164) error
	//View formatted table of details for a specific number (with history).
	View(number *E164) (string, error)
	//Summary formatted table of usage stats
	Summary() (string, error)
}

//ValidE164 validates an phonenumber is E164
func (phoneNumber E164) ValidE164() error {
	if ok, _ := regexp.MatchString(`^[1-9][0-9]{0,2}$`, phoneNumber.Cc); !ok {
		return errors.New("Invalid country code in phone number")
	}
	if ok, _ := regexp.MatchString(`^[01][1-9][0-9]{0,3}$`, phoneNumber.Ndc); !ok {
		return errors.New("Invalid destination code in phone number")
	}
	if ok, _ := regexp.MatchString(`^[0-9]{5,13}$`, phoneNumber.Sn); !ok {
		return errors.New("Invalid subscriber number in phone number")
	}
	return nil
}

//ValidUserID validates userid format.
func ValidUserID(userID *int64) error {
	if *userID == 0 {
		return errors.New("User id invalid")
	}
	return nil
}
