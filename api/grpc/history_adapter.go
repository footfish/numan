package grpc

import (
	context "context"
	"errors"

	"github.com/footfish/numan"
	"google.golang.org/grpc"
)

//Adaptors are used to facilitate transparent gRPC transport.
//They adapt the service interface to gRPC interface and visa versa.
//ie. Client application (main) -> Service Interface -> ClientAdapter -> grpc transport -> ServiceAdapter -> Service Interface (app)

//historyClientAdapter implements an adapter from HistoryService to historyClient(gGRPC).
type historyClientAdapter struct {
	grpc *historyClient
}

// NewNumberingClientAdapter instantiates NumberingClientAdaptor
func NewHistoryClientAdapter(conn *grpc.ClientConn) numan.HistoryService {
	c := NewHistoryClient(conn)
	return &historyClientAdapter{c.(*historyClient)}
}

//AddHistory  implements HistoryService.AddHistory()
func (c *historyClientAdapter) AddHistory(ctx context.Context, historyEntry numan.History) error {
	return errors.New("Method AddHistory not available via gGRPC")
}

//ListHistoryByNumber implements HistoryService.ListHistoryByNumber()
func (c *historyClientAdapter) ListHistoryByNumber(ctx context.Context, phoneNumber numan.E164) (historyList []numan.History, err error) {
	err = phoneNumber.ValidE164()
	if err != nil {
		return
	}
	listHistoryResponse, err := c.grpc.ListHistoryByNumber(ctx, &ListHistoryByNumberRequest{E164: &E164{Cc: phoneNumber.Cc, Ndc: phoneNumber.Ndc, Sn: phoneNumber.Sn}})
	if err == nil {
		for _, hist := range listHistoryResponse.HistoryEntry {
			historyList = append(historyList, *unMarshalHistory(hist))
		}
	}
	return
}

//ListHistoryByOwnerID implements HistoryService.ListHistoryByUserId()
func (c *historyClientAdapter) ListHistoryByOwnerID(ctx context.Context, ownerID int64) (historyList []numan.History, err error) {
	err = numan.ValidOwnerID(&ownerID)
	if err != nil {
		return
	}
	listHistoryResponse, err := c.grpc.ListHistoryByOID(ctx, &ListHistoryByOIDRequest{OwnerID: ownerID})
	if err == nil {
		for _, hist := range listHistoryResponse.HistoryEntry {
			historyList = append(historyList, *unMarshalHistory(hist))
		}
	}
	return
}

//historyServerAdapter implements an adapter from HistoryServer(gRPC) to HistoryService.
type historyServerAdapter struct {
	service numan.HistoryService
	UnimplementedHistoryServer
}

//ListHistoryByNumber implements HistoryServer.ListHistoryByNumber()
func (h *historyServerAdapter) ListHistoryByNumber(ctx context.Context, in *ListHistoryByNumberRequest) (*ListHistoryResponse, error) {

	historyList, err := h.service.ListHistoryByNumber(ctx, numan.E164{Cc: in.E164.Cc, Ndc: in.E164.Ndc, Sn: in.E164.Sn})
	if err != nil {
		return nil, err
	}

	resp := &ListHistoryResponse{}
	for _, historyEntry := range historyList {
		resp.HistoryEntry = append(resp.HistoryEntry, MarshalHistory(&historyEntry))
	}
	return resp, err
}

//ListHistoryByOID implements HistoryServer.ListHistoryByOID()
func (h *historyServerAdapter) ListHistoryByOID(ctx context.Context, in *ListHistoryByOIDRequest) (*ListHistoryResponse, error) {
	historyList, err := h.service.ListHistoryByOwnerID(ctx, in.OwnerID)
	if err != nil {
		return nil, err
	}

	resp := &ListHistoryResponse{}
	for _, historyEntry := range historyList {
		resp.HistoryEntry = append(resp.HistoryEntry, MarshalHistory(&historyEntry))
	}
	return resp, err
}

func MarshalHistory(h *numan.History) *HistoryEntry {
	if h == nil {
		return &HistoryEntry{}
	}
	return &HistoryEntry{
		Timestamp: h.Timestamp,
		E164:      &E164{Cc: h.E164.Cc, Ndc: h.E164.Ndc, Sn: h.E164.Sn},
		OwnerID:   h.OwnerID,
		Action:    h.Action,
		Notes:     h.Notes,
	}
}

func unMarshalHistory(h *HistoryEntry) *numan.History {
	if h == nil {
		return &numan.History{}
	}
	return &numan.History{
		Timestamp: h.Timestamp,
		E164:      numan.E164{Cc: h.E164.Cc, Ndc: h.E164.Ndc, Sn: h.E164.Sn},
		OwnerID:   h.OwnerID,
		Action:    h.Action,
		Notes:     h.Notes,
	}
}
