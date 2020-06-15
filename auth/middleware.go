package auth

import (
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

var (
	//Todo: init this

	MiddlewareJWTAuth = echo.MiddlewareFunc(middlewareJWTAuth(true))
)

const (
	bearer = "Bearer"
)

func middlewareJWTAuth(returnJSON bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get(echo.HeaderAuthorization)

			l := len(bearer)

			if len(auth) < l+1 || auth[:l] != bearer {
				c.Set("user", Component.UserGuest)
				return next(c)
			}

			tdata, err := TokenValidate(auth[l+1:])
			if err != nil {
				return err
			}

			db, err := Component.GetDriver().FromContext(c)
			if err != nil {
				return errors.New("Error connecting to db")
			}

			user, err := db.FindUserByID(tdata.ID)

			if err != nil {
				return err
			}
			c.Set("user", user)
			c.Set("claims", tdata)
			return next(c)

		}
	}
}

type TokenData struct {
	ID  float64
	Iat float64
	Exp float64
}

func TokenValidate(token string) (*TokenData, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return PubKey, nil
	})
	if err != nil || !t.Valid {
		switch err.(type) {
		case *jwt.ValidationError:
			vErr := err.(*jwt.ValidationError)
			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				return nil, errors.New("Token expired")
			default:
				return nil, errors.New("Token invalid")
			}
		}
	}

	tokenData := t.Claims.(jwt.MapClaims)
	return &TokenData{tokenData["ID"].(float64), tokenData["iat"].(float64), tokenData["exp"].(float64)}, nil
}
