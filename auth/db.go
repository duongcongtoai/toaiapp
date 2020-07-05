package auth

// DB aka auth repository
type DB interface {
	SetupFromYamlConfig(configFile string) error

	FindUserByID(id float64) (*User, error)
	FindUserByName(name string) (*User, error)
	CreateUser(*User) error

	// Initialize() error
	// DB() (AuthDB, error)
	// FromContext(echo.Context) (AuthDB, error)
}

// GetDB enable access to auth repository
func GetDB() DB {
	return Component.driver
}
