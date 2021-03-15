package storage

import "github.com/footfish/numan"

//GetHistoryByNumber implements HistoryAPI.GetHistoryByNumber()
func (s *store) GetHistoryByNumber(phoneNumber numan.E164) (history []numan.History, err error) {
	return
}

//GetHistoryByUserID implements HistoryAPI.GetHistoryByUserId()
func (s *store) GetHistoryByUserID(userID int64) (history []numan.History, err error) {
	return
}
