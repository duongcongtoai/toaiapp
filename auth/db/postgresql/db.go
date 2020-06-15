package mysql

import (
	"errors"
	"fmt"

	"toaiapp/auth"
	pgcom "toaiapp/db/postgresql"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo/v4"
)

var (
	driverName = "auth_db_pg"
)

type pgDB struct {
	db     *gorm.DB
	config auth.Configuration
}

func (m *pgDB) SetupEcho(e *echo.Echo) error {
	if auth.Component.GetConfig().AuthDBType == "postgresql" {
		url := m.config.AuthDBUrl
		if url == "" {
			return errors.New("Empty url configuration for component auth")
		}
		e.Use(pgcom.PostgresMiddleware(driverName, pgcom.ConnectDB(url)))
	}
	return nil
}

func (m *pgDB) FindUserByID(id float64) (*auth.User, error) {
	var user auth.User
	if m.db.First(&user, uint(id)).RecordNotFound() {
		return nil, fmt.Errorf("User not found with id %f", id)
	}
	return &user, nil
}

func (m *pgDB) FindUserByName(name string) (*auth.User, error) {
	var user auth.User
	if m.db.Where(&auth.User{Name: name}).First(&user).RecordNotFound() {
		return nil, fmt.Errorf("User not found with name %s", name)
	}
	return &user, nil
}

func (m *pgDB) SetConfig(cfg auth.Configuration) {
	m.config = cfg
}

func (m *pgDB) CreateUser(u *auth.User) error {
	return m.db.Create(u).Error
}

func (m *pgDB) DB() (auth.AuthDB, error) {

	url := m.config.AuthDBUrl
	if url == "" {
		return nil, errors.New("Empty url configuration for component auth")
	}
	db := pgcom.ConnectDB(url)
	cloned := &pgDB{config: m.config, db: db}
	return cloned, nil
}
func (m *pgDB) FromContext(c echo.Context) (auth.AuthDB, error) {
	if db, ok := c.Get(driverName).(*gorm.DB); ok {
		return &pgDB{db: db, config: m.config}, nil
	}
	return nil, errors.New("Postgresql Database is not properly configured")
}

func (m *pgDB) Initialize() error {
	if db := m.db.AutoMigrate(&auth.User{}); db.Error != nil {
		return db.Error
	}
	return nil
}

func init() {
	//Component has mysql driver now
	auth.Component.RegisterDB("postgresql", &pgDB{})
}
