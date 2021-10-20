package app

import (
	"context"
	"errors"
	"strings"

	"github.com/footfish/numan"
	"github.com/footfish/numan/internal/app/datastore"
	"golang.org/x/crypto/bcrypt"
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
func (s *userService) Auth(ctx context.Context, username string, password string) (user numan.User, err error) {
	//TODO better sanity check user/pass
	if len(password) < 8 && len(password) > 32 {
		return user, errors.New("Bad password format")
	}
	password = strings.ToLower(password)
	username = strings.ToLower(username)

	if user, err = s.next.Auth(ctx, username, password); err == nil {
		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		if err != nil {
			return user, errors.New("Username/password mismatch")
		}
		user.SetNewAccessToken()
	}
	return
}

/*
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return userdata, fmt.Errorf("cannot hash password: %w", err)
	}
*/
