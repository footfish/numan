package datastore

import (
	"context"

	"github.com/footfish/numan"
)

// numberService implements the UserService interface
type userService struct {
	store Store
}

// NewUserService instantiates a new UserService .
func NewUserService(store *Store) numan.UserService {
	return &userService{
		store: *store,
	}
}

//Auth implements UserService.Auth()
func (s *userService) Auth(ctx context.Context, username string, password string) (userdata numan.User, err error) {
	row := s.store.db.QueryRow("SELECT id, username, passwordhash, role FROM user where username=?", username)
	row.Scan(&userdata.UID, &userdata.Username, &userdata.Password, &userdata.Role)
	return userdata, err
}

//AddUser implements UserService.AddUser()
func (s *userService) AddUser(ctx context.Context, user numan.User) error {
	_, err := s.store.db.Exec("INSERT INTO user(username, passwordhash, role) values(?,?,?)", user.Username, user.Password, user.Role)
	if err != nil {
		return err
	}
	return nil
}
