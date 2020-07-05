package mysql

// import (
// 	"errors"
// 	"fmt"

// 	"toaiapp/auth"
// 	mysqlcom "toaiapp/db/mysql"

// 	"github.com/jinzhu/gorm"
// 	_ "github.com/jinzhu/gorm/dialects/mysql"
// 	"github.com/labstack/echo/v4"
// )

// var (
// 	driverName = "auth_db_mysql"
// )

// type mysqlDB struct {
// 	db     *gorm.DB
// 	config auth.Configuration
// }

// func (m *mysqlDB) SetupEcho(e *echo.Echo) error {

// 	if auth.Component.GetConfig().AuthDBType == "mysql" {
// 		url := m.config.AuthDBUrl
// 		if url == "" {
// 			return errors.New("Empty url configuration for component auth")
// 		}
// 		e.Use(mysqlcom.MysqlMiddleware(driverName, mysqlcom.ConnectDB(url)))
// 	}
// 	return nil
// }

// func (m *mysqlDB) FindUserByID(id float64) (*auth.User, error) {
// 	var user auth.User
// 	if m.db.First(&user, id).RecordNotFound() {
// 		return nil, fmt.Errorf("User not found with id %f", id)
// 	}
// 	return &user, nil
// }

// func (m *mysqlDB) FindUserByName(name string) (*auth.User, error) {
// 	var user auth.User
// 	if m.db.Where(&auth.User{Name: name}).First(&user).RecordNotFound() {
// 		return nil, fmt.Errorf("User not found with id %s", name)
// 	}
// 	return &user, nil
// }

// func (m *mysqlDB) SetConfig(cfg auth.Configuration) {
// 	m.config = cfg
// }

// func (m *mysqlDB) CreateUser(u *auth.User) error {
// 	return m.db.Create(u).Error
// }

// func (m *mysqlDB) DB() (auth.AuthDB, error) {

// 	url := m.config.AuthDBUrl
// 	if url == "" {
// 		return nil, errors.New("Empty url configuration for component auth")
// 	}
// 	db := mysqlcom.ConnectDB(url)
// 	cloned := &mysqlDB{config: m.config, db: db}
// 	return cloned, nil
// }
// func (m *mysqlDB) FromContext(c echo.Context) (auth.AuthDB, error) {
// 	if db, ok := c.Get(driverName).(*gorm.DB); ok {
// 		return &mysqlDB{db: db, config: m.config}, nil
// 	}
// 	return nil, errors.New("Mysql Database is not properly configured")
// }

// func (m *mysqlDB) Initialize() error {
// 	if db := m.db.AutoMigrate(&auth.User{}); db.Error != nil {
// 		return db.Error
// 	}
// 	return nil
// }

// func init() {
// 	//Component has mysql driver now
// 	auth.Component.RegisterDB("mysql", &mysqlDB{})
// }
