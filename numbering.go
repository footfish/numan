package numan

import (
	"context"
	"errors"
	"regexp"
)

const (
	// MAXRESERVATIONTIME the maximum time a number can be reserved (pre-allocated)
	MAXRESERVATIONTIME = 108000 //in secs (30 mins)
	//DATEPRINTFORMAT the format used for display of timestamps
	DATEPRINTFORMAT = "02/01/2006"
	//TIMESTAMPPRINTFORMAT the format used for display of timestamps
	TIMESTAMPPRINTFORMAT = "02/01/06 15:04"
	//QUARANTINE period in seconds. Numbers can't be reserved/allocated during quarantie.
	QUARANTINE = 13 * 31 * 24 * 60 * 60 //  (13 months approx)
)

//Numbering represents a stored phone number entry
type Numbering struct {
	ID          int64  // number entry index
	E164        E164   //an e.164 number
	Used        bool   // indicates number is actively being used or not (reserved, allocated).
	Domain      string // which domain is using the number (which domain can allocate)
	Carrier     string // who is the block owner
	OwnerID     int64  // which client/customer currently 'owns' the number
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
	OwnerID     int64  // which client a/c is using
	Allocated   bool   // if the number was ordered
	Reserved    bool   // if the number is reserved
	DeAllocated bool   // if number was last cancelled (use for quarantine)
	PortedIn    bool   // if number was ported in
	PortedOut   bool   // if number was ported out
}

//NumberingService exposes interface for managing numbers
type NumberingService interface {
	//Adds a new unused number to database.
	//params E164, Domain & Carrier must be included, others (supplied or not) are initialised.
	Add(ctx context.Context, number *Numbering) error
	//AddGroup adds a series of new unused numbers
	AddGroup(ctx context.Context)
	//List returns a filtered list of numbers
	List(ctx context.Context, filter *NumberFilter) ([]Numbering, error)
	//ListOwnerID gets list of numbers attached to specific OwnerID
	ListOwnerID(ctx context.Context, ownerID int64) ([]Numbering, error)
	//Reserve locks a number to a OwnerID until untilTS (unix timestamp)
	Reserve(ctx context.Context, number *E164, ownerID *int64, untilTS *int64) error
	//Allocate marks a number 'used' by a User
	Allocate(ctx context.Context, number *E164, ownerID *int64) error
	//DeAllocate number from User (number goes to quarantine)
	DeAllocate(ctx context.Context, number *E164) error
	//Portout sets a port out date (just a log, doesn't care about state or do anything else)
	Portout(ctx context.Context, number *E164, PortoutTS *int64) error
	//Portin sets a port in date (just a log, doesn't care about state or do anything else)
	Portin(ctx context.Context, number *E164, PortinTS *int64) error
	//Delete - number no longer used, removed from number db, must be unused (history kept).
	Delete(ctx context.Context, number *E164) error
	//View formatted table of details for a specific number (with history).
	View(ctx context.Context, number *E164) (string, error)
	//Summary formatted table of usage stats
	Summary(ctx context.Context) (string, error)
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

//ValidOwnerID validates ownerID format.
func ValidOwnerID(ownerID *int64) error {
	if *ownerID == 0 {
		return errors.New("User id invalid")
	}
	return nil
}
