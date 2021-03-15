package app

import "github.com/footfish/numan"

//GetHistoryByNumber implements HistoryAPI.GetHistoryByNumber()
func (s *numanService) GetHistoryByNumber(phoneNumber numan.E164) (history []numan.History, err error) {
	return
}

//GetHistoryByUserID implements HistoryAPI.GetHistoryByUserId()
func (s *numanService) GetHistoryByUserID(userID int64) (history []numan.History, err error) {
	return
}
