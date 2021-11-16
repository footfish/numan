package datastore

import (
	"context"
	"errors"

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

//DeleteUser  implements UserService.DeleteUser
func (s *userService) DeleteUser(ctx context.Context, username string) error {
	row, err := s.store.db.Exec("DELETE from user WHERE username=?", username)
	if err != nil {
		return err
	}
	if n, _ := row.RowsAffected(); n == 0 { //ok for sqlite. RowsAffected may not be supported with other drivers.
		return errors.New("Unable to delete, check the username exists")
	}
	return nil
}

//ListUsers  implements UserService.DeleteUser
func (s *userService) ListUsers(ctx context.Context, userfilter string) (userList []numan.User, err error) {
	var result numan.User
	var resultList []numan.User
	userfilter = userfilter + "%"
	rows, err := s.store.db.Query("SELECT username, role FROM user where username like ? ", userfilter)
	if err != nil {
		return resultList, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&result.Username,
			&result.Role,
		)
		if err != nil {
			return resultList, err
		}
		resultList = append(resultList, result)
	}
	err = rows.Err()
	if err != nil {
		return resultList, err
	}
	return resultList, nil
}
