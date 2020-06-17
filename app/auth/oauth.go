package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	stderr "errors"
	"toaiapp/auth"

	"github.com/go-session/session"
	"github.com/labstack/echo/v4"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
)

func registerOauthRoutes(e *echo.Echo) {
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	manager.MustTokenStorage(store.NewMemoryTokenStore())
	clientStore := store.NewClientStore()
	//Hard code client id
	clientStore.Set("client_app_id", &models.Client{
		ID:     "client_app_id",
		Secret: "client_secret",
		Domain: "http://localhost:8084",
	})
	manager.MapClientStorage(clientStore)

	srv := server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)

	srv.SetUserAuthorizationHandler(userAuthorizeHandlerFunc(srv, clientStore))
	manager.SetRefreshTokenCfg(manage.DefaultRefreshTokenCfg)
	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Printf("Internal Error:%v", err)
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Printf("Response Error: %v", re.Error)
	})
	oauth := e.Group("/oauth")

	oauth.POST("/get_token", tokenFunc(srv))
	oauth.GET("/authorize", authorizeFunc(srv))
	oauth.GET("/login", loginGetFunc())
	oauth.POST("/login", loginPost)
}

func loginPost(c echo.Context) error {

	store, err := session.Start(nil, c.Response(), c.Request())

	type userData = struct {
		Username string `form:"username"`
		Password string `form:"password"`
	}

	input := &userData{}
	if err := c.Bind(input); err != nil {
		return c.HTML(http.StatusBadRequest, "invalid input")
	}
	if input.Username == "" || input.Password == "" {

		return c.HTML(http.StatusBadRequest, "missing information")
	}

	db, err := auth.Component.GetDriver().FromContext(c)
	if err != nil {

		return c.HTML(http.StatusBadRequest, fmt.Sprintf("%v\n", err.Error()))
	}
	user, err := db.FindUserByName(input.Username)
	if err != nil {
		return c.HTML(http.StatusBadRequest, fmt.Sprintf("%v\n", err.Error()))
	}
	if err = user.Authenticate(input.Password); err != nil {
		return c.HTML(http.StatusBadRequest, "wrong username or password")
	}
	store.Set("userid", strconv.Itoa(int(user.ID)))
	store.Save()
	return c.HTML(http.StatusFound, fmt.Sprintf("Session is ready, try oauth2! for user %s", user.Name))
}

func userAuthorizeHandlerFunc(server *server.Server, store *store.ClientStore) server.UserAuthorizationHandler {
	return func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
		store, err := session.Start(nil, w, r)
		if err != nil {
			return "", err
		}
		uid, ok := store.Get("userid")
		if !ok {
			return "", stderr.New("No session found")
		}
		store.Save()
		return uid.(string), nil
	}
}

func loginGetFunc() echo.HandlerFunc {
	httpHandlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		outputHTML(w, r, "static/login.html")
	})
	return echo.WrapHandler(httpHandlerFunc)
}

func outputHTML(w http.ResponseWriter, r *http.Request, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()
	fi, _ := file.Stat()
	http.ServeContent(w, r, file.Name(), fi.ModTime(), file)
}

func authorizeFunc(srv *server.Server) echo.HandlerFunc {
	httpHandlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		err := srv.HandleAuthorizeRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})
	return echo.WrapHandler(httpHandlerFunc)
}

func tokenFunc(srv *server.Server) echo.HandlerFunc {
	httpHandlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := srv.HandleTokenRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})
	return echo.WrapHandler(httpHandlerFunc)
}

// func credentialFunc(clientStore *store.ClientStore) func(echo.Context, *auth.User) error {
// 	return func(c echo.Context, user *auth.User) error {
// 		clientID := uuid.New().String()[:8]
// 		clientSecret := uuid.New().String()[:8]
// 		err := clientStore.Set(clientID, &models.Client{
// 			ID:     clientID,
// 			Secret: clientSecret,
// 			Domain: "http://localhost:8082",
// 			UserID: strconv.Itoa(int(user.ID)),
// 		})
// 		if err != nil {

// 			fmt.Println(err.Error())
// 		}
// 		return c.JSON(http.StatusOK, map[string]string{"CLIENT_ID": clientID, "CLIENT_SECRET": clientSecret})
// 	}
// }
