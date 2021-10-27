package grpc

import (
	"context"

	"github.com/footfish/numan"
	"github.com/footfish/numan/internal/service"
	"github.com/footfish/numan/internal/service/datastore"
	"google.golang.org/grpc"
)

//Adaptors are used to facilitate transparent gRPC transport.
//They adapt the service interface to gRPC interface and visa versa.
//ie. Client servicelication (main) -> Service Interface -> ClientAdapter -> grpc transport -> ServiceAdapter -> Service Interface (service)

//userClientAdapter implements an adapter from UserService to UserClient(grpc).
type userClientAdapter struct {
	grpc *userClient
}

// NewUserClientAdapter instantiates userClientAdaptor
func NewUserClientAdapter(conn *grpc.ClientConn) numan.UserService {
	c := NewUserClient(conn)
	return &userClientAdapter{c.(*userClient)}
}

// Auth implements UserService.Auth()
func (c *userClientAdapter) Auth(ctx context.Context, username string, password string) (user numan.User, err error) {
	resp, err := c.grpc.Auth(ctx, &AuthRequest{Username: username, Password: password})
	if err == nil {
		user = numan.User{UID: resp.Uid, Username: resp.Username, PasswordHash: resp.Passwordhash, AccessToken: resp.Token, Role: resp.Role}
	}
	return user, err
}

//userServerAdapter implements an Adapter from UserServer(grpc) to UserService.
type userServerAdapter struct {
	service numan.UserService
	UnimplementedUserServer
}

// NewUserServerAdapter creates a new UserServerAdapter
func NewUserServerAdapter(store *datastore.Store) UserServer {
	return &userServerAdapter{service: service.NewUserService(store)}
}

//Auth implements UserServer.Auth()
func (s *userServerAdapter) Auth(ctx context.Context, auth *AuthRequest) (resp *AuthResponse, err error) {
	user, err := s.service.Auth(ctx, auth.Username, auth.Password)
	return &AuthResponse{Uid: user.UID, Username: user.Username, Passwordhash: user.PasswordHash, Token: user.AccessToken}, err
}
