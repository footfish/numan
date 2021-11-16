package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/footfish/numan"
	"github.com/footfish/numan/internal/service/auth"
	"github.com/footfish/numan/internal/service/datastore"
)

//userService implements the UserService interface
type userService struct {
	next numan.UserService
}

// NewUserService instantiates a new UserService.
func NewUserService(store *datastore.Store) numan.UserService {
	return &userService{
		next: auth.NewUserService(store),
	}
}

//Auth implements UserService.Auth()
func (s *userService) Auth(ctx context.Context, username string, password string) (numan.User, error) {
	//sanity checks
	enteredUser := numan.User{Username: strings.ToLower(username), Password: strings.ToLower(password)}
	if !enteredUser.ValidRawPassword() {
		return enteredUser, errors.New("Invalid password")
	}
	if !enteredUser.ValidUsername() {
		return enteredUser, errors.New("Invalid Username")
	}
	//Fetch User from store (password ignored)
	storedUser, err := s.next.Auth(ctx, enteredUser.Username, enteredUser.Password)
	//Authenticate
	if err == nil {
		err = storedUser.ComparePassword(enteredUser.Password)
		if err != nil {
			return enteredUser, errors.New("Username/password mismatch")
		}
		storedUser.SetNewAccessToken()
	}
	//Note: if public login should be obfiscating error here
	return storedUser, err
}

//AddUser implements UserService.AddUser()
func (s *userService) AddUser(ctx context.Context, user numan.User) (err error) {
	//sanity checks
	if !(user.Role == numan.RoleAdmin || user.Role == numan.RoleUser) {
		return errors.New("incompatible role")
	}
	if !user.ValidUsername() {
		return errors.New("bad username")
	}
	//Hash password if needed
	if !user.PasswordIsHashed() {
		if !user.ValidRawPassword() {
			return errors.New("bad password")
		}
		if err = user.HashPassword(); err != nil {
			return fmt.Errorf("can't hash password: %w", err)
		}
	}
	//store
	return s.next.AddUser(ctx, user)
}

//DeleteUser  implements UserService.DeleteUser
func (s *userService) DeleteUser(ctx context.Context, username string) error {
	u := numan.User{Username: username}
	if !u.ValidUsername() {
		return errors.New("Invalid Username")
	}
	return s.next.DeleteUser(ctx, username)
}

//ListUsers  implements UserService.DeleteUser
func (s *userService) ListUsers(ctx context.Context, userfilter string) ([]numan.User, error) {
	return s.next.ListUsers(ctx, userfilter)
}
