package grpc

import (
	context "context"

	"github.com/footfish/numan"
	"github.com/footfish/numan/internal/app"
	"github.com/footfish/numan/internal/datastore"
	grpc "google.golang.org/grpc"
)

//userClientAdapter is used to implement Adapter from User to UserClient.
type userClientAdapter struct {
	user *userClient
}

// NewUserClientAdapter instantiates userClientAdaptor
func NewUserClientAdapter(conn *grpc.ClientConn) numan.UserService {
	c := NewUserClient(conn)
	return &userClientAdapter{c.(*userClient)}
}

// Auth implements UserService.Auth()
func (c *userClientAdapter) Auth(ctx context.Context, username string, password string) (user numan.User, err error) {
	resp, err := c.user.Auth(ctx, &AuthRequest{Username: username, Password: password})
	if err == nil {
		user = numan.User{UID: resp.Uid, Username: resp.Username, PasswordHash: resp.Passwordhash, AccessToken: resp.Token, Role: resp.Role}
	}
	return user, err
}

//userServerAdapter server is used to implement Adapter from UserServer to User.
type userServerAdapter struct {
	user numan.UserService
	UnimplementedUserServer
}

// NewUserServerAdapter creates a new UserServerAdapter
func NewUserServerAdapter(store *datastore.Store) UserServer {
	return &userServerAdapter{user: app.NewUserService(store)}
}

//Auth implements UserServer.Auth()
func (s *userServerAdapter) Auth(ctx context.Context, auth *AuthRequest) (resp *AuthResponse, err error) {
	user, err := s.user.Auth(ctx, auth.Username, auth.Password)
	return &AuthResponse{Uid: user.UID, Username: user.Username, Passwordhash: user.PasswordHash, Token: user.AccessToken}, err
}
