package numan

//History represents a stored phone numbers history
type History struct {
	ID   int64 // number entry index
	E164 E164  //an e.164 number
}

//HistoryAPI exposes interface for number history
type HistoryAPI interface {
	//GetHistoryByNumber gets history for a specific phone number
	GetHistoryByNumber(phoneNumber E164) ([]History, error)
	//GetHistoryByUserID gets history for a specific UserID
	GetHistoryByUserID(userID int64) ([]History, error)
}
