package auth

import (
	"context"

	"github.com/footfish/numan"
	"github.com/footfish/numan/internal/service/datastore"
)

//userService implements the UserService interface
type userService struct {
	next numan.UserService
}

// NewUserService instantiates a new UserService.
func NewUserService(store *datastore.Store) numan.UserService {
	return &userService{
		next: datastore.NewUserService(store),
	}
}

//Auth implements UserService.Auth()
func (s *userService) Auth(ctx context.Context, username string, password string) (numan.User, error) {
	return s.next.Auth(ctx, username, password)
}

//AddUser implements UserService.AddUser()
func (s *userService) AddUser(ctx context.Context, user numan.User) (err error) {
	if err := checkUserRole(numan.RoleAdmin, ctx); err != nil {
		return err
	}
	return s.next.AddUser(ctx, user)
}
