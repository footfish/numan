package grpc

import (
	"context"
	"time"

	"github.com/footfish/numan"
	"github.com/footfish/numan/internal/app"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

//numanClientAdapter server is used to implement Adapter from Numan to NumanClient.
type numanClientAdapter struct {
	nc   *numanClient
	conn *grpc.ClientConn
}

// NewNumanClientAdapter instantiates NumanClientAdaptor
func NewNumanClientAdapter(address string, creds credentials.TransportCredentials) numan.API {
	// Set up a connection to the server.
	//conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))
	if err != nil {
		panic(err)
	}
	nc := NewNumanClient(conn)
	return &numanClientAdapter{nc.(*numanClient), conn}
}

// Close closes grpc connection
func (n *numanClientAdapter) Close() {
	n.conn.Close()
}

// Add implements NumberAPI.Add()
func (n *numanClientAdapter) Add(number *numan.Number) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := n.nc.Add(ctx, &AddRequest{Number: marshalNumber(number)})
	return err
}

//AddGroup not implemented
func (n *numanClientAdapter) AddGroup() {
}

// List implements NumberAPI.List()
func (n *numanClientAdapter) List(filter *numan.NumberFilter) (numbers []numan.Number, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	numberList, err := n.nc.List(ctx, &ListRequest{NumberFilter: marshalNumberFilter(filter)})
	for _, number := range numberList.Number {
		numbers = append(numbers, *unMarshalNumber(number))
	}
	return
}

// ListUserID implements NumberAPI.ListUserID()
func (n *numanClientAdapter) ListUserID(userID int64) (numbers []numan.Number, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	numberList, err := n.nc.ListUserID(ctx, &ListUserIDRequest{UserID: userID})
	for _, number := range numberList.Number {
		numbers = append(numbers, *unMarshalNumber(number))
	}

	return
}

//Reserve implements NumberAPI.Reserve()
func (n *numanClientAdapter) Reserve(number *numan.E164, userID *int64, untilTS *int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := n.nc.Reserve(ctx, &ReserveRequest{E164: marshalE164(number), UserID: *userID, UntilTS: *untilTS})
	return err
}

//Allocate implements NumberAPI.Allocate()
func (n *numanClientAdapter) Allocate(number *numan.E164, userID *int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := n.nc.Allocate(ctx, &AllocateRequest{E164: marshalE164(number), UserID: *userID})
	return err
}

//DeAllocate implements NumberAPI.DeAllocate()
func (n *numanClientAdapter) DeAllocate(number *numan.E164) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := n.nc.DeAllocate(ctx, &DeAllocateRequest{E164: marshalE164(number)})
	return err
}

//Portout implements NumberAPI.Portout()
func (n *numanClientAdapter) Portout(number *numan.E164, portoutTS *int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := n.nc.Portout(ctx, &PortoutRequest{E164: marshalE164(number), PortoutTS: *portoutTS})
	return err
}

//Portin implements NumberAPI.Portin()
func (n *numanClientAdapter) Portin(number *numan.E164, portinTS *int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := n.nc.Portin(ctx, &PortinRequest{E164: marshalE164(number), PortinTS: *portinTS})
	return err
}

//Delete implements NumberAPI.Delete()
func (n *numanClientAdapter) Delete(number *numan.E164) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := n.nc.Delete(ctx, &DeleteRequest{E164: marshalE164(number)})
	return err
}

//View implements NumberAPI.View()
func (n *numanClientAdapter) View(number *numan.E164) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := n.nc.View(ctx, &ViewRequest{E164: marshalE164(number)})
	return resp.Message, err
}

//Summary implements NumberAPI.Summary()
func (n *numanClientAdapter) Summary() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := n.nc.Summary(ctx, &SummaryRequest{})
	return resp.Message, err

}

//GetHistoryByNumber implements HistoryAPI.GetHistoryByNumber()
func (n *numanClientAdapter) GetHistoryByNumber(phoneNumber numan.E164) (history []numan.History, err error) {
	return
}

//GetHistoryByUserID implements HistoryAPI.GetHistoryByUserId()
func (n *numanClientAdapter) GetHistoryByUserID(userID int64) (history []numan.History, err error) {
	return
}

//numanServerAdapter server is used to implement Adapter from NumanServer to Numan.
type numanServerAdapter struct {
	//nu *app.NumanService
	nu numan.API
	UnimplementedNumanServer
}

// NewGrpcServer creates a new grpc.Server and NumanServerAdapter
func NewGrpcServer(dsn string, creds credentials.TransportCredentials) (*grpc.Server, NumanServer) {
	return grpc.NewServer(grpc.Creds(creds)), &numanServerAdapter{nu: app.NewNumanService(dsn)}
}

//CloseServerAdapter shuts db connection
func CloseServerAdapter(s NumanServer) {
	s.(*numanServerAdapter).nu.Close()
}

//Add implements NumanServer.Add()
func (s *numanServerAdapter) Add(ctx context.Context, in *AddRequest) (*AddResponse, error) {
	err := s.nu.Add(unMarshalNumber(in.Number))
	return &AddResponse{}, err
}

//List implements NumanServer.List()
func (s *numanServerAdapter) List(ctx context.Context, in *ListRequest) (*ListResponse, error) {
	numberFilter := unMarshalNumberFilter(in.NumberFilter)
	numberList, err := s.nu.List(numberFilter)
	if err != nil {
		return nil, err
	}

	resp := &ListResponse{}
	for _, number := range numberList {
		resp.Number = append(resp.Number, marshalNumber(&number))
	}
	return resp, err
}

//ListUserID implements NumanServer.ListUserID()
func (s *numanServerAdapter) ListUserID(ctx context.Context, in *ListUserIDRequest) (*ListUserIDResponse, error) {

	numberList, err := s.nu.ListUserID(in.UserID)
	if err != nil {
		return nil, err
	}

	resp := &ListUserIDResponse{}
	for _, number := range numberList {
		resp.Number = append(resp.Number, marshalNumber(&number))
	}
	return resp, err
}

//Reserve implements NumanServer.Reserve()
func (s *numanServerAdapter) Reserve(ctx context.Context, in *ReserveRequest) (*ReserveResponse, error) {
	err := s.nu.Reserve(unMarshalE164(in.E164), &in.UserID, &in.UntilTS)
	return &ReserveResponse{}, err
}

//Allocate  implements NumanServer.Reserve()
func (s *numanServerAdapter) Allocate(ctx context.Context, in *AllocateRequest) (*AllocateResponse, error) {
	err := s.nu.Allocate(unMarshalE164(in.E164), &in.UserID)
	return &AllocateResponse{}, err
}

//DeAllocate  implements NumanServer.DeAllocate()
func (s *numanServerAdapter) DeAllocate(ctx context.Context, in *DeAllocateRequest) (*DeAllocateResponse, error) {
	err := s.nu.DeAllocate(unMarshalE164(in.E164))
	return &DeAllocateResponse{}, err
}

//Portout  implements NumanServer.Portout()
func (s *numanServerAdapter) Portout(ctx context.Context, in *PortoutRequest) (*PortoutResponse, error) {
	err := s.nu.Portout(unMarshalE164(in.E164), &in.PortoutTS)
	return &PortoutResponse{}, err
}

//Portin  implements NumanServer.Portin()
func (s *numanServerAdapter) Portin(ctx context.Context, in *PortinRequest) (*PortinResponse, error) {
	err := s.nu.Portin(unMarshalE164(in.E164), &in.PortinTS)
	return &PortinResponse{}, err
}

//Delete  implements NumanServer.Delete()
func (s *numanServerAdapter) Delete(ctx context.Context, in *DeleteRequest) (resp *DeleteResponse, err error) {
	err = s.nu.Delete(unMarshalE164(in.E164))
	return &DeleteResponse{}, err
}

//View  implements NumanServer.View()
func (s *numanServerAdapter) View(ctx context.Context, in *ViewRequest) (*ViewResponse, error) {
	message, err := s.nu.View(unMarshalE164(in.E164))
	return &ViewResponse{Message: message}, err
}

//Summary implements NumanServer.Summary()
func (s *numanServerAdapter) Summary(ctx context.Context, in *SummaryRequest) (*SummaryResponse, error) {
	message, err := s.nu.Summary()
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

func marshalNumber(n *numan.Number) *Number {
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

func unMarshalNumber(n *Number) *numan.Number {
	if n == nil || n.E164 == nil {
		return &numan.Number{}
	}
	return &numan.Number{ID: n.Id,
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
