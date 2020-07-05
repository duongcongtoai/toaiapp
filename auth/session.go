package auth

import "github.com/gorilla/sessions"

type SessionStore interface {
	// SetupEcho(e *echo.Echo) error
	SetupFromYamlConfig(configFile string) error
	Store() sessions.Store

	// Initialize() error
	// Store() (sessions.Store, error)
	// StoreFromContext(echo.Context) (sessions.Store, error)
}

// func SessionStoreFromContext(c echo.Context) (sessions.Store, error) {
// 	return Component.sessionStore.StoreFromContext(c)
// }

func GetSessionStore() sessions.Store {
	return Component.sessionStore.Store()
}
