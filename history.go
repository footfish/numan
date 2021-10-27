package numan

import "context"

//History represents a stored phone numbers history
type History struct {
	Timestamp int64
	E164      E164   //an e.164 number
	OwnerID   int64  //who the number was allocated to
	Action    string //what command action is logged
	Notes     string //additional notes
}

//HistoryService exposes interface for number history
type HistoryService interface {
	//AddHistory adds history for a specific phone number
	AddHistory(ctx context.Context, historyEntry History) error
	//ListHistoryByNumber gets history log for a specific phone number
	ListHistoryByNumber(ctx context.Context, phoneNumber E164) ([]History, error)
	//ListHistoryByOwnerID gets history log for a specific OwnerID
	ListHistoryByOwnerID(ctx context.Context, ownerID int64) ([]History, error)
}
