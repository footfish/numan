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
		user = numan.User{UID: resp.Uid, Username: resp.Username, Password: resp.Passwordhash, AccessToken: resp.Token, Role: resp.Role}
	}
	return user, err
}

//AddUser implements UserService.AddUser()
func (c *userClientAdapter) AddUser(ctx context.Context, user numan.User) (err error) {
	_, err = c.grpc.AddUser(ctx, &AddUserRequest{Username: user.Username, Password: user.Password, Role: user.Role})
	return err
}

//DeleteUser  implements UserService.DeleteUser
func (c *userClientAdapter) DeleteUser(ctx context.Context, username string) (err error) {
	_, err = c.grpc.DeleteUser(ctx, &DeleteUserRequest{Username: username})
	return err
}

//ListUsers  implements UserService.ListUsers
func (c *userClientAdapter) ListUsers(ctx context.Context, userfilter string) (userlist []numan.User, err error) {
	listUsersRequest := ListUsersRequest{Userfilter: userfilter}
	listUsersResponse, err := c.grpc.ListUsers(ctx, &listUsersRequest)
	if err == nil {
		for _, user := range listUsersResponse.Userlist {
			userlist = append(userlist, numan.User{Username: user.Username, Role: user.Role})
		}
	}
	return
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
	return &AuthResponse{Uid: user.UID, Username: user.Username, Passwordhash: user.Password, Token: user.AccessToken}, err
}

//AddUser implements UserServer.AddUser()
func (s *userServerAdapter) AddUser(ctx context.Context, in *AddUserRequest) (resp *AddUserResponse, err error) {
	return &AddUserResponse{}, s.service.AddUser(ctx, numan.User{Username: in.Username, Password: in.Password, Role: in.Role})
}

//ListsUsers implements UserServer.ListsUsers()
func (s *userServerAdapter) ListUsers(ctx context.Context, in *ListUsersRequest) (*ListUsersResponse, error) {
	userList, err := s.service.ListUsers(ctx, in.Userfilter)
	if err != nil {
		return nil, err
	}

	var resp ListUsersResponse
	for _, userEntry := range userList {
		resp.Userlist = append(resp.Userlist, &UserEntry{Username: userEntry.Username, Role: userEntry.Role})
	}

	return &resp, err
}

//DeleteUser implements UserServer.DeleteUser()
func (s *userServerAdapter) DeleteUser(ctx context.Context, in *DeleteUserRequest) (resp *DeleteUserResponse, err error) {
	return &DeleteUserResponse{}, s.service.DeleteUser(ctx, in.Username)
}
