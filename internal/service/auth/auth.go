package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/footfish/numan"
)

//compareUserRole checks given role against the user role extracted from JWT token (in context)
func compareUserRole(requiredRole string, ctx context.Context) error {
	user := &numan.User{}
	if err := user.SetUserFromToken(fmt.Sprintf("%s", ctx.Value("token"))); err != nil { //Get authenticated user data from token
		return errors.New("Unexpected Auth error")
	}
	if user.Role != requiredRole {
		return errors.New("Insufficient user privileges")
	}
	return nil
}
