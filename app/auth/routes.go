package auth

import (
	"toaiapp/auth"
	"net/http"

	"github.com/labstack/echo/v4"
)

func apiRegisterRoutes(e *echo.Echo) {
	gr := e.Group("/api/v1/user")

	gr.POST("/login", postLogin)
	gr.Use(auth.MiddlewareJWTAuth)
	gr.GET("/", auth.AuthorizationWrapper(getUser, auth.UserGet))
}

type ResultGet struct {
	User auth.User
}

func getUser(c echo.Context, u *auth.User) error {
	return c.JSON(http.StatusOK, ResultGet{*u})
	// return nil

}

func postLogin(c echo.Context) error {
	type userData = struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	input := &userData{}
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid input"})
	}
	if input.Username == "" || input.Password == "" {
		return c.JSON(http.StatusBadRequest,
			map[string]string{"message": "Missing information"})
	}

	db, err := auth.Component.GetDriver().FromContext(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}
	user, err := db.FindUserByName(input.Username)
	if err != nil {
		return c.JSON(http.StatusNonAuthoritativeInfo,
			map[string]string{"message": err.Error()})
	}
	if err = user.Authenticate(input.Password); err != nil {
		return c.JSON(http.StatusNonAuthoritativeInfo,
			map[string]string{"message": "wrong username or password"})
	}
	return sendToken(c, user)
}

func sendToken(c echo.Context, u *auth.User) error {
	token, err := u.GenerateToken()
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			map[string]string{"message": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
