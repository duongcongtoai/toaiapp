package auth

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
)

//AuthorizationWrapper provide user parameter in callback function
func AuthorizationWrapper(callable func(c echo.Context, u *User) error, permission string) func(echo.Context) error {
	return func(c echo.Context) error {
		user, ok := c.Get(userContextKey).(*User)
		if !ok {
			return errors.New("No user found in context")
		}

		if permission != "" {
			if !user.HasPermission(permission) {
				return fmt.Errorf("User has no permission for %s", permission)
			}
		}
		return callable(c, user)
	}
}
