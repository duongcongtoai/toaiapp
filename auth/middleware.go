package auth

import (
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"

	"github.com/labstack/echo/v4"
)

const (
	bearer = "Bearer"
)

func MiddlewareSessionAuth() echo.MiddlewareFunc {
	store := GetSessionStore()
	db := GetDB()
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session, err := store.Get(c.Request(), "userid")
			if err != nil {
				return err
			}
			uid := session.Values[useridSessionKey]
			userID, ok := uid.(float64)
			if !ok {
				return fmt.Errorf("Session value of %s is not valid", "userid")
			}
			user, err := db.FindUserByID(userID)
			if err != nil {
				return err
			}
			c.Set(userContextKey, user)
			return nil
		}
	}
}
func MiddlewareJWTAuth() echo.MiddlewareFunc {
	db := GetDB()
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			auth := c.Request().Header.Get(echo.HeaderAuthorization)

			l := len(bearer)

			if len(auth) < l+1 || auth[:l] != bearer {
				c.Set(userContextKey, Component.UserGuest)
				return next(c)
			}

			tdata, err := TokenValidate(auth[l+1:])
			if err != nil {
				return err
			}

			user, err := db.FindUserByID(tdata.ID)

			if err != nil {
				return err
			}
			c.Set(userContextKey, user)
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
