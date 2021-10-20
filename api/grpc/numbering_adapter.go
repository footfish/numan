package grpc

import (
	"context"

	"github.com/footfish/numan"
	"github.com/footfish/numan/internal/app"
	"github.com/footfish/numan/internal/app/datastore"
	"google.golang.org/grpc"
)

//numberingClientAdapter is used to implement Adapter from Numbering to NumberingClient.
type numberingClientAdapter struct {
	numbering *numberingClient
}

// NewNumberingClientAdapter instantiates NumberingClientAdaptor
func NewNumberingClientAdapter(conn *grpc.ClientConn) numan.NumberingService {
	c := NewNumberingClient(conn)
	return &numberingClientAdapter{c.(*numberingClient)}
}

/*
// NewNumberingClientAdapter instantiates NumberingClientAdaptor
func NewNumberingClientAdapter(address string, creds credentials.TransportCredentials) numan.NumberingService {
	// Set up a connection to the server.
	//conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds),
		grpc.WithUnaryInterceptor(authClientInterceptor))
	if err != nil {
		panic(err)
	}
	c := NewNumberingClient(conn)
	return &numberingClientAdapter{c.(*numberingClient)}
}
*/
// Add implements NumberingService.Add()
func (c *numberingClientAdapter) Add(ctx context.Context, number *numan.Numbering) error {
	_, err := c.numbering.Add(ctx, &AddRequest{Number: marshalNumber(number)})
	return err
}

//AddGroup not implemented
func (c *numberingClientAdapter) AddGroup(ctx context.Context) {
}

// List implements NumberingService.List()
func (c *numberingClientAdapter) List(ctx context.Context, filter *numan.NumberFilter) (numbers []numan.Numbering, err error) {
	numberList, err := c.numbering.List(ctx, &ListRequest{NumberFilter: marshalNumberFilter(filter)})
	if err == nil {
		for _, number := range numberList.Number {
			numbers = append(numbers, *unMarshalNumber(number))
		}
	}
	return
}

// ListUserID implements NumberingService.ListUserID()
func (c *numberingClientAdapter) ListUserID(ctx context.Context, userID int64) (numbers []numan.Numbering, err error) {
	numberList, err := c.numbering.ListUserID(ctx, &ListUserIDRequest{UserID: userID})
	if err == nil {
		for _, number := range numberList.Number {
			numbers = append(numbers, *unMarshalNumber(number))
		}
	}
	return
}

//Reserve implements NumberingService.Reserve()
func (c *numberingClientAdapter) Reserve(ctx context.Context, number *numan.E164, userID *int64, untilTS *int64) error {
	_, err := c.numbering.Reserve(ctx, &ReserveRequest{E164: marshalE164(number), UserID: *userID, UntilTS: *untilTS})
	return err
}

//Allocate implements NumberingService.Allocate()
func (c *numberingClientAdapter) Allocate(ctx context.Context, number *numan.E164, userID *int64) error {
	_, err := c.numbering.Allocate(ctx, &AllocateRequest{E164: marshalE164(number), UserID: *userID})
	return err
}

//DeAllocate implements NumberingService.DeAllocate()
func (c *numberingClientAdapter) DeAllocate(ctx context.Context, number *numan.E164) error {
	_, err := c.numbering.DeAllocate(ctx, &DeAllocateRequest{E164: marshalE164(number)})
	return err
}

//Portout implements NumberingService.Portout()
func (c *numberingClientAdapter) Portout(ctx context.Context, number *numan.E164, portoutTS *int64) error {
	_, err := c.numbering.Portout(ctx, &PortoutRequest{E164: marshalE164(number), PortoutTS: *portoutTS})
	return err
}

//Portin implements NumberingService.Portin()
func (c *numberingClientAdapter) Portin(ctx context.Context, number *numan.E164, portinTS *int64) error {
	_, err := c.numbering.Portin(ctx, &PortinRequest{E164: marshalE164(number), PortinTS: *portinTS})
	return err
}

//Delete implements NumberingService.Delete()
func (c *numberingClientAdapter) Delete(ctx context.Context, number *numan.E164) error {
	_, err := c.numbering.Delete(ctx, &DeleteRequest{E164: marshalE164(number)})
	return err
}

//View implements NumberingService.View()
func (c *numberingClientAdapter) View(ctx context.Context, number *numan.E164) (string, error) {
	resp, err := c.numbering.View(ctx, &ViewRequest{E164: marshalE164(number)})
	if err != nil {
		return err.Error(), err
	}
	return resp.Message, err
}

//Summary implements NumberingService.Summary()
func (c *numberingClientAdapter) Summary(ctx context.Context) (string, error) {
	resp, err := c.numbering.Summary(ctx, &SummaryRequest{})
	if err != nil {
		return err.Error(), err
	}
	return resp.Message, err

}

//GetHistoryByNumber implements HistoryAPI.GetHistoryByNumber()
func (c *numberingClientAdapter) GetHistoryByNumber(ctx context.Context, phoneNumber numan.E164) (history []numan.History, err error) {
	return
}

//GetHistoryByUserID implements HistoryAPI.GetHistoryByUserId()
func (c *numberingClientAdapter) GetHistoryByUserID(ctx context.Context, userID int64) (history []numan.History, err error) {
	return
}

//numberingServerAdapter server is used to implement Adapter from NumberingServer to Numan.
type numberingServerAdapter struct {
	//numbering *app.NumanService
	numbering numan.NumberingService
	UnimplementedNumberingServer
}

// NewNumberingServerAdapter creates a new NumberingServerAdapter
func NewNumberingServerAdapter(store *datastore.Store) NumberingServer {
	return &numberingServerAdapter{numbering: app.NewNumberingService(store)}
}

//Add implements NumberingServer.Add()
func (s *numberingServerAdapter) Add(ctx context.Context, in *AddRequest) (*AddResponse, error) {
	err := s.numbering.Add(ctx, unMarshalNumber(in.Number))
	return &AddResponse{}, err
}

//List implements NumberingServer.List()
func (s *numberingServerAdapter) List(ctx context.Context, in *ListRequest) (*ListResponse, error) {
	numberFilter := unMarshalNumberFilter(in.NumberFilter)
	numberList, err := s.numbering.List(ctx, numberFilter)
	if err != nil {
		return nil, err
	}

	resp := &ListResponse{}
	for _, number := range numberList {
		resp.Number = append(resp.Number, marshalNumber(&number))
	}
	return resp, err
}

//ListUserID implements NumberingServer.ListUserID()
func (s *numberingServerAdapter) ListUserID(ctx context.Context, in *ListUserIDRequest) (*ListUserIDResponse, error) {

	numberList, err := s.numbering.ListUserID(ctx, in.UserID)
	if err != nil {
		return nil, err
	}

	resp := &ListUserIDResponse{}
	for _, number := range numberList {
		resp.Number = append(resp.Number, marshalNumber(&number))
	}
	return resp, err
}

//Reserve implements NumberingServer.Reserve()
func (s *numberingServerAdapter) Reserve(ctx context.Context, in *ReserveRequest) (*ReserveResponse, error) {
	err := s.numbering.Reserve(ctx, unMarshalE164(in.E164), &in.UserID, &in.UntilTS)
	return &ReserveResponse{}, err
}

//Allocate  implements NumberingServer.Reserve()
func (s *numberingServerAdapter) Allocate(ctx context.Context, in *AllocateRequest) (*AllocateResponse, error) {
	err := s.numbering.Allocate(ctx, unMarshalE164(in.E164), &in.UserID)
	return &AllocateResponse{}, err
}

//DeAllocate  implements NumberingServer.DeAllocate()
func (s *numberingServerAdapter) DeAllocate(ctx context.Context, in *DeAllocateRequest) (*DeAllocateResponse, error) {
	err := s.numbering.DeAllocate(ctx, unMarshalE164(in.E164))
	return &DeAllocateResponse{}, err
}

//Portout  implements NumberingServer.Portout()
func (s *numberingServerAdapter) Portout(ctx context.Context, in *PortoutRequest) (*PortoutResponse, error) {
	err := s.numbering.Portout(ctx, unMarshalE164(in.E164), &in.PortoutTS)
	return &PortoutResponse{}, err
}

//Portin  implements NumberingServer.Portin()
func (s *numberingServerAdapter) Portin(ctx context.Context, in *PortinRequest) (*PortinResponse, error) {
	err := s.numbering.Portin(ctx, unMarshalE164(in.E164), &in.PortinTS)
	return &PortinResponse{}, err
}

//Delete  implements NumberingServer.Delete()
func (s *numberingServerAdapter) Delete(ctx context.Context, in *DeleteRequest) (resp *DeleteResponse, err error) {
	err = s.numbering.Delete(ctx, unMarshalE164(in.E164))
	return &DeleteResponse{}, err
}

//View  implements NumberingServer.View()
func (s *numberingServerAdapter) View(ctx context.Context, in *ViewRequest) (*ViewResponse, error) {
	message, err := s.numbering.View(ctx, unMarshalE164(in.E164))
	return &ViewResponse{Message: message}, err
}

//Summary implements NumberingServer.Summary()
func (s *numberingServerAdapter) Summary(ctx context.Context, in *SummaryRequest) (*SummaryResponse, error) {
	message, err := s.numbering.Summary(ctx)
	return &SummaryResponse{Message: message}, err
}

func marshalNumberFilter(n *numan.NumberFilter) *NumberFilter {
	if n == nil {
		return &NumberFilter{}
	}
	return &NumberFilter{Id: int64(n.ID),
		E164:        &E164{Cc: n.E164.Cc, Ndc: n.E164.Ndc, Sn: n.E164.Sn},
		State:       int32(n.State),
		Domain:      n.Domain,
		Carrier:     n.Carrier,
		UserID:      int64(n.UserID),
		Allocated:   n.Allocated,
		Reserved:    n.Reserved,
		DeAllocated: n.DeAllocated,
		PortedIn:    n.PortedIn,
		PortedOut:   n.PortedOut,
	}
}

func unMarshalNumberFilter(n *NumberFilter) *numan.NumberFilter {
	if n == nil {
		return &numan.NumberFilter{}
	}
	numberFilter := &numan.NumberFilter{ID: n.Id,
		State:       byte(n.State),
		Domain:      n.Domain,
		Carrier:     n.Carrier,
		UserID:      n.UserID,
		Allocated:   n.Allocated,
		Reserved:    n.Reserved,
		DeAllocated: n.DeAllocated,
		PortedIn:    n.PortedIn,
		PortedOut:   n.PortedOut,
	}
	if n.E164 != nil {
		numberFilter.E164 = numan.E164{Cc: n.E164.Cc, Ndc: n.E164.Ndc, Sn: n.E164.Sn}
	}
	return numberFilter
}

func marshalE164(n *numan.E164) *E164 {
	if n == nil {
		return &E164{}
	}
	return &E164{Cc: n.Cc, Ndc: n.Ndc, Sn: n.Sn}
}

func unMarshalE164(n *E164) *numan.E164 {
	if n == nil {
		return &numan.E164{}
	}
	return &numan.E164{Cc: n.Cc, Ndc: n.Ndc, Sn: n.Sn}
}

func marshalNumber(n *numan.Numbering) *Number {
	if n == nil {
		return &Number{}
	}
	return &Number{Id: int64(n.ID),
		E164:        &E164{Cc: n.E164.Cc, Ndc: n.E164.Ndc, Sn: n.E164.Sn},
		Used:        n.Used,
		Domain:      n.Domain,
		Carrier:     n.Carrier,
		UserID:      n.UserID,
		Allocated:   n.Allocated,
		DeAllocated: n.DeAllocated,
		PortedIn:    n.PortedIn,
		PortedOut:   n.PortedOut,
	}
}

func unMarshalNumber(n *Number) *numan.Numbering {
	if n == nil || n.E164 == nil {
		return &numan.Numbering{}
	}
	return &numan.Numbering{ID: n.Id,
		E164:        numan.E164{Cc: n.E164.Cc, Ndc: n.E164.Ndc, Sn: n.E164.Sn},
		Used:        n.Used,
		Domain:      n.Domain,
		Carrier:     n.Carrier,
		UserID:      n.UserID,
		Allocated:   n.Allocated,
		DeAllocated: n.DeAllocated,
		PortedIn:    n.PortedIn,
		PortedOut:   n.PortedOut,
	}
}
