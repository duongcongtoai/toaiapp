package mysql

import (
	"fmt"
	"io/ioutil"

	"toaiapp/auth"
	pgcom "toaiapp/db/postgresql"

	"github.com/jinzhu/gorm"
	"gopkg.in/yaml.v2"
)

var (
	driverName = "auth_db_pg"
)

type pgDB struct {
	db     *gorm.DB
	config Configuration
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

// func (m *pgDB) SetConfig(config []byte) {
// 	m.config = cfg
// }

func (m *pgDB) CreateUser(u *auth.User) error {
	return m.db.Create(u).Error
}

type Configuration struct {
	DBURL string `yaml:"auth_db_url"`
}

func (m *pgDB) SetupFromYamlConfig(configFile string) error {
	configByte, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	type ConfigContainer struct {
		Configuration `yaml:"auth_db_postgres"`
	}
	var cfgContainer ConfigContainer
	if err := yaml.Unmarshal(configByte, &cfgContainer); err != nil {
		return err
	}

	db := pgcom.ConnectDB(cfgContainer.Configuration.DBURL)
	m.db = db

	m.config = cfgContainer.Configuration
	if db := m.db.AutoMigrate(&auth.User{}); db.Error != nil {
		return db.Error
	}
	return nil
}

func init() {
	//Component has mysql driver now
	auth.Component.RegisterDB("postgresql", &pgDB{})
}
