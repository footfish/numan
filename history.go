package numan

import "context"

//History represents a stored phone numbers history
type History struct {
	ID   int64 // number entry index
	E164 E164  //an e.164 number
}

//HistoryService exposes interface for number history
type HistoryService interface {
	//GetHistoryByNumber gets history for a specific phone number
	GetHistoryByNumber(ctx context.Context, phoneNumber E164) ([]History, error)
	//GetHistoryByUserID gets history for a specific UserID
	GetHistoryByUserID(ctx context.Context, userID int64) ([]History, error)
}
