package storage

import (
	"errors"
	"fmt"
	"strings"
	"time"

	//"time"
	"github.com/footfish/numan"
	// register sqlite driver
	_ "modernc.org/sqlite"
)

// Add implements NumberAPI.Add()
func (s *Store) Add(number *numan.Number) error {
	_, err := s.db.Exec("INSERT INTO number(cc, ndc, sn, domain, carrier) values(?,?,?,?,?)", number.E164.Cc, number.E164.Ndc, number.E164.Sn, number.Domain, number.Carrier)
	if err != nil {
		return err
	}
	return nil
}

//List implements NumberAPI.List()
func (s *Store) List(filter *numan.NumberFilter) ([]numan.Number, error) {
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
	if v := filter.State; v != 0 {
		where, args = append(where, "used = ?"), append(args, ((v-1) != 0)) // convert State->Used(bool)
	}
	if v := filter.Domain; len(v) != 0 {
		where, args = append(where, "domain = ?"), append(args, v)
	}

	var result numan.Number
	var resultList []numan.Number

	rows, err := s.db.Query("SELECT * FROM number where "+strings.Join(where, " AND "), args...)
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
			&result.UserID,
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

//Summary implements NumberAPI.Summary()
func (s *Store) Summary() (string, error) {
	summary := fmt.Sprintf("%-15v %5v %5v %5v %5v %5v\n", "Domain", "CC", "NDC", "Used", "Free", "Total")
	rows, err := s.db.Query("SELECT domain, cc, ndc, sum(used) as used, count(*)-sum(used) as free,  count(*) as total from number group by domain,cc,ndc; ")
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

//Delete implements NumberAPI.Delete()
func (s *Store) Delete(phonenumber *numan.E164) error {
	row, err := s.db.Exec("DELETE from number where used == 0 and cc=? and ndc=? and sn=?", phonenumber.Cc, phonenumber.Ndc, phonenumber.Sn)
	if err != nil {
		return err
	}
	if n, _ := row.RowsAffected(); n == 0 { //ok for sqlite. RowsAffected may not be supported with other drivers.
		return errors.New("Unable to delete, check the number")
	}
	//TODO - log to history
	return nil
}

//Reserve implements NumberAPI.Reserve()
//Mark 'used' & set userID & reserved date.
//Numbers must be out of quarantine
func (s *Store) Reserve(number *numan.E164, userID *int, untilTS *int64) error {
	row, err := s.db.Exec("UPDATE number set used=1, deallocated=0, reserved=?, userID=? where reserved==0 and used==0 and userID==0 and cc=? and ndc=? and sn=? and deallocated<?", *untilTS, *userID, number.Cc, number.Ndc, number.Sn, time.Now().Unix()-numan.QUARANTINE)
	if err != nil {
		return err
	}
	if n, _ := row.RowsAffected(); n == 0 { //ok for sqlite. RowsAffected may not be supported with other drivers.
		return errors.New("Unable to reserve number (check number, already reserved?)")
	}
	//TODO - log to history
	return nil
}

//Allocate implements NumberAPI.Allocate()
//Mark 'used' & set userID & allocation date. Reset reservation & de-allocation flag
//Numbers must be out of quarantine
func (s *Store) Allocate(number *numan.E164, userID *int) error {
	row, err := s.db.Exec("UPDATE number set used=1, deallocated=0, reserved=0, allocated=?, userID=? where used==0 and userID==0 and cc=? and ndc=? and sn=? and deallocated<?", time.Now().Unix(), *userID, number.Cc, number.Ndc, number.Sn, time.Now().Unix()-numan.QUARANTINE)
	if err != nil {
		return err
	}
	if n, _ := row.RowsAffected(); n == 0 { //ok for sqlite. RowsAffected may not be supported with other drivers.
		return errors.New("Unable to allocate number (check number, already allocated?)")
	}
	//TODO - log to history
	return nil
}

//DeAllocate implements NumberAPI.DeAllocate()
//Mark 'unused' & set de-allocation date (quarantine). Resets  userID, reservation & allocation dateflag.
func (s *Store) DeAllocate(number *numan.E164) error {
	row, err := s.db.Exec("UPDATE number set used=0, deallocated=?, reserved=0, allocated=0, userID=0 where used==1 and cc=? and ndc=? and sn=? and deallocated=0", time.Now().Unix(), number.Cc, number.Ndc, number.Sn)
	if err != nil {
		return err
	}
	if n, _ := row.RowsAffected(); n == 0 { //ok for sqlite. RowsAffected may not be supported with other drivers.
		return errors.New("Unable to de-allocate number (check number, already de-allocated?)")
	}
	//TODO - log to history
	return nil
}

//Portout implements NumberAPI.Portout()
func (s *Store) Portout(number *numan.E164, PortoutTS *int64) error {
	row, err := s.db.Exec("UPDATE number set portedOut=? where  cc=? and ndc=? and sn=?", *PortoutTS, number.Cc, number.Ndc, number.Sn)
	if err != nil {
		return err
	}
	if n, _ := row.RowsAffected(); n == 0 { //ok for sqlite. RowsAffected may not be supported with other drivers.
		return errors.New("Unable to set ported out date. db update failed, check number ")
	}
	//TODO - log to history
	return nil
}

//Portin implements NumberAPI.Portin()
func (s *Store) Portin(number *numan.E164, PortinTS *int64) error {
	row, err := s.db.Exec("UPDATE number set portedIn=? where  cc=? and ndc=? and sn=?", *PortinTS, number.Cc, number.Ndc, number.Sn)
	if err != nil {
		return err
	}
	if n, _ := row.RowsAffected(); n == 0 { //ok for sqlite. RowsAffected may not be supported with other drivers.
		return errors.New("Unable to set ported in date. db update failed, check number ")
	}
	//TODO - log to history
	return nil
}
