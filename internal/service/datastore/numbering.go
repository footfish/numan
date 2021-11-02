package datastore

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	//"time"
	"github.com/footfish/numan"
	// register sqlite driver
	_ "modernc.org/sqlite"
)

// numberingService implements the NumberingService interface
type numberingService struct {
	store Store
}

// NewNumberingService instantiates a NumberingService .
func NewNumberingService(store *Store) numan.NumberingService {
	return &numberingService{
		store: *store,
	}
}

// Add implements NumberingService.Add()
func (s *numberingService) Add(ctx context.Context, number *numan.Numbering) error {
	_, err := s.store.db.Exec("INSERT INTO number(cc, ndc, sn, domain, carrier) values(?,?,?,?,?)", number.E164.Cc, number.E164.Ndc, number.E164.Sn, number.Domain, number.Carrier)
	if err != nil {
		return err
	}
	return nil
}

// AddGroup not implemented
func (s *numberingService) AddGroup(ctx context.Context) {
}

//List implements NumberingService.List()
func (s *numberingService) List(ctx context.Context, filter *numan.NumberFilter) ([]numan.Numbering, error) {
	//build WHERE args from filter
	where, args := []string{"1 = 1"}, []interface{}{}
	if v := filter.E164.Cc; len(v) != 0 {
		where, args = append(where, "cc = ?"), append(args, v)
	}
	if v := filter.E164.Ndc; len(v) != 0 {
		where, args = append(where, "ndc = ?"), append(args, v)
	}
	if v := filter.E164.Sn; len(v) != 0 {
		where, args = append(where, "sn like ?"), append(args, v+"%")
	}
	if v := filter.ID; v != 0 {
		where, args = append(where, "id = ?"), append(args, v)
	}
	if v := filter.OwnerID; v != 0 {
		where, args = append(where, "ownerID = ?"), append(args, v)
	}
	if v := filter.State; v != 0 {
		where, args = append(where, "used = ?"), append(args, ((v-1) != 0)) // convert State->Used(bool)
	}
	if v := filter.Domain; len(v) != 0 {
		where, args = append(where, "domain = ?"), append(args, v)
	}

	var result numan.Numbering
	var resultList []numan.Numbering

	rows, err := s.store.db.Query("SELECT * FROM number where "+strings.Join(where, " AND "), args...)
	if err != nil {
		return resultList, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&result.ID,
			&result.E164.Cc,
			&result.E164.Ndc,
			&result.E164.Sn,
			&result.Used,
			&result.Domain,
			&result.Carrier,
			&result.OwnerID,
			&result.Allocated,
			&result.Reserved,
			&result.DeAllocated,
			&result.PortedIn,
			&result.PortedOut)
		if err != nil {
			return resultList, err
		}
		resultList = append(resultList, result)
	}
	err = rows.Err()
	if err != nil {
		return resultList, err

	}
	return resultList, nil
}

//ListOwnerID implements NumberingService.ListOwnerID()
func (s *numberingService) ListOwnerID(ctx context.Context, oid int64) ([]numan.Numbering, error) {
	filter := &numan.NumberFilter{OwnerID: oid}
	return s.List(ctx, filter)
}

//Summary implements NumberingService.Summary()
func (s *numberingService) Summary(ctx context.Context) (string, error) {
	summary := fmt.Sprintf("%-15v %5v %5v %5v %5v %5v\n", "Domain", "CC", "NDC", "Used", "Free", "Total")
	rows, err := s.store.db.Query("SELECT domain, cc, ndc, sum(used) as used, count(*)-sum(used) as free,  count(*) as total from number group by domain,cc,ndc; ")
	if err != nil {
		return summary, err
	}
	defer rows.Close()
	for rows.Next() {
		var row struct {
			domain string
			cc     string
			ndc    string
			used   int
			free   int
			total  int
		}
		err = rows.Scan(
			&row.domain,
			&row.cc,
			&row.ndc,
			&row.used,
			&row.free,
			&row.total,
		)
		if err != nil {
			return summary, err
		}
		summary += fmt.Sprintf("%-15v %5v %5v %5v %5v %5v\n", row.domain, row.cc, row.ndc, row.used, row.free, row.total)
	}
	if err != nil {
		return summary, err
	}
	return summary, nil
}

//Delete implements NumberingService.Delete()
func (s *numberingService) Delete(ctx context.Context, phonenumber *numan.E164) error {
	row, err := s.store.db.Exec("DELETE from number where used == 0 and cc=? and ndc=? and sn=?", phonenumber.Cc, phonenumber.Ndc, phonenumber.Sn)
	if err != nil {
		return err
	}
	if n, _ := row.RowsAffected(); n == 0 { //ok for sqlite. RowsAffected may not be supported with other drivers.
		return errors.New("Unable to delete, check the number")
	}
	return nil
}

//View implements NumberingService.View()
func (s *numberingService) View(ctx context.Context, number *numan.E164) (view string, err error) {
	result, err := s.List(ctx, &numan.NumberFilter{E164: *number})
	if err != nil {
		return "", err
	}

	for _, r := range result {
		view += fmt.Sprintf("#%d) +%v-%v-%v, Domain: %v, Carrier: %v\n", r.ID, r.E164.Cc, r.E164.Ndc, r.E164.Sn, r.Domain, r.Carrier)
		if r.Used { //Used can be reserved or allocated
			if r.Allocated > 0 {
				view += fmt.Sprintf("Allocated to OwnerID: %v on %v\n", r.OwnerID, time.Unix(int64(r.Allocated), 0).Format(numan.DATEPRINTFORMAT))
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

	return view, nil
}

//Reserve implements NumberingService.Reserve()
//Mark 'used' & set ownerID & reserved date.
//Numbers must be out of quarantine
func (s *numberingService) Reserve(ctx context.Context, number *numan.E164, ownerID *int64, untilTS *int64) error {
	row, err := s.store.db.Exec("UPDATE number set used=1, deallocated=0, reserved=?, ownerID=? where reserved==0 and used==0 and ownerID==0 and cc=? and ndc=? and sn=? and deallocated<?", *untilTS, *ownerID, number.Cc, number.Ndc, number.Sn, time.Now().Unix()-numan.QUARANTINE)
	if err != nil {
		return err
	}
	if n, _ := row.RowsAffected(); n == 0 { //ok for sqlite. RowsAffected may not be supported with other drivers.
		return errors.New("Unable to reserve number (check number, already reserved?)")
	}

	return nil
}

//Allocate implements NumberingService.Allocate()
//Mark 'used' & set ownerID & allocation date. Reset reservation & de-allocation flag
//Numbers must be out of quarantine
func (s *numberingService) Allocate(ctx context.Context, number *numan.E164, ownerID *int64) error {
	row, err := s.store.db.Exec("UPDATE number set used=1, deallocated=0, reserved=0, allocated=?, ownerID=? where used==0 and ownerID==0 and cc=? and ndc=? and sn=? and deallocated<?", time.Now().Unix(), *ownerID, number.Cc, number.Ndc, number.Sn, time.Now().Unix()-numan.QUARANTINE)
	if err != nil {
		return err
	}
	if n, _ := row.RowsAffected(); n == 0 { //ok for sqlite. RowsAffected may not be supported with other drivers.
		return errors.New("Unable to allocate number (check number, already allocated?)")
	}
	return nil
}

//DeAllocate implements NumberingService.DeAllocate()
//Mark 'unused' & set de-allocation date (quarantine). Resets  ownerID, reservation & allocation dateflag.
func (s *numberingService) DeAllocate(ctx context.Context, number *numan.E164, ownerID *int64) error {
	row, err := s.store.db.Exec("UPDATE number set used=0, deallocated=?, reserved=0, allocated=0, ownerID=0 where used==1 and cc=? and ndc=? and sn=? and ownerID=? and deallocated=0", time.Now().Unix(), number.Cc, number.Ndc, number.Sn, ownerID)
	if err != nil {
		return err
	}
	if n, _ := row.RowsAffected(); n == 0 { //ok for sqlite. RowsAffected may not be supported with other drivers.
		return errors.New("Unable to de-allocate number (Wrong owner?, already de-allocated?)")
	}
	return nil
}

//Portout implements NumberingService.Portout()
func (s *numberingService) Portout(ctx context.Context, number *numan.E164, PortoutTS *int64) error {
	row, err := s.store.db.Exec("UPDATE number set portedOut=? where  cc=? and ndc=? and sn=?", *PortoutTS, number.Cc, number.Ndc, number.Sn)
	if err != nil {
		return err
	}
	if n, _ := row.RowsAffected(); n == 0 { //ok for sqlite. RowsAffected may not be supported with other drivers.
		return errors.New("Unable to set ported out date. db update failed, check number ")
	}
	return nil
}

//Portin implements NumberingService.Portin()
func (s *numberingService) Portin(ctx context.Context, number *numan.E164, PortinTS *int64) error {
	row, err := s.store.db.Exec("UPDATE number set portedIn=? where  cc=? and ndc=? and sn=?", *PortinTS, number.Cc, number.Ndc, number.Sn)
	if err != nil {
		return err
	}
	if n, _ := row.RowsAffected(); n == 0 { //ok for sqlite. RowsAffected may not be supported with other drivers.
		return errors.New("Unable to set ported in date. db update failed, check number ")
	}
	return nil
}
