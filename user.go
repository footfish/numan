package numan

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

//add a user
//find a user
//auth a user
//validate a token (internal)
const (
	secretKey     = "secret"
	tokenDuration = 15 * time.Minute
	//tokenDuration  = 1 * time.Minute //TODO testing
	AuthTokenField     = "token" //field name to use in ctx and meta data for storing auth token
	RoleUser           = "user"
	RoleAdmin          = "admin"
	PatternUser        = "^[1-9a-z]{3,13}$"
	PatternRawPassword = "^[1-9a-zA-Z._%-]{5,13}$"
)

type User struct {
	UID         int64
	Username    string
	Password    string //can be hashed or raw
	Role        string
	AccessToken string
}

//UserService exposes interface for managing users
type UserService interface {
	//Auth authenticates a user and returns a copy of user data with JWT token
	Auth(ctx context.Context, Username string, Password string) (user User, err error)
	//AddUser adds a new user
	AddUser(ctx context.Context, user User) error
	//DeleteUser removes a new user
	DeleteUser(ctx context.Context, username string) error
	//ListUsers returns an array of matching users.
	ListUsers(ctx context.Context, userfilter string) ([]User, error)
	//SetPassword changes a users password
	SetPassword(ctx context.Context, username string, newPassword string) error
}

//userClaims is JWT claims object
type userClaims struct {
	jwt.StandardClaims
	UID      int64  `json:"uid"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

//SetGeneratedToken creates a JWT access token in UserAuth struct
func (u *User) SetNewAccessToken() (err error) {
	// Set claims
	claims := userClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenDuration).Unix(),
		},
		UID:      u.UID,
		Username: u.Username,
		Role:     u.Role,
	} // create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign and store the complete encoded access token as a string
	u.AccessToken, err = token.SignedString([]byte(secretKey))
	return
}

//SetUserWithToken verifies & reads claims into userAuth from raw accessToken
func (u *User) SetUserFromToken(accessToken string) (err error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&userClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("Auth error: unexpected token signing method")
			}
			return []byte(secretKey), nil
		},
	)

	if err != nil {
		return fmt.Errorf("Auth error: invalid token: %w", err)
	}

	claims, ok := token.Claims.(*userClaims)
	if !ok {
		return fmt.Errorf("Auth error: invalid token claims")
	}
	u.UID = claims.UID
	u.Username = claims.Username
	u.Role = claims.Role
	return nil
}

//AuthRefreshRequired returns true if UserAuth.AccessToken expired/invalid.
//for client use, token is parsed unverified.
func (u *User) AuthRefreshRequired() bool {
	var p jwt.Parser
	token, _, err := p.ParseUnverified(u.AccessToken, &userClaims{})
	if err != nil {
		return true
	}

	claims, ok := token.Claims.(*userClaims)
	if !ok {
		return true
	}

	if claims.StandardClaims.ExpiresAt < time.Now().Unix() {
		return true
	}
	return false
}

//ValidUsername checks if the format of User.Username is valid
func (u *User) ValidUsername() bool {
	res, _ := regexp.MatchString(PatternUser, u.Username)
	return res
}

//ValidPasssword checks if User.Password is valid RAW password
func (u *User) ValidRawPassword() bool {
	res, _ := regexp.MatchString(PatternRawPassword, u.Password)
	return res
}

//PasswordIsHashed checks if User.Password is hashed
func (u *User) PasswordIsHashed() bool {
	if _, err := bcrypt.Cost([]byte(u.Password)); err != nil {
		return false
	}
	return true
}

//HashPassword will hash User.Password
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err == nil {
		u.Password = string(hashedPassword)
	}
	return err
}

func (u *User) ComparePassword(password string) error {
	if u.PasswordIsHashed() {
		return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	}
	if !(u.Password == password) {
		return errors.New("password mismatch")
	}
	return nil
}
