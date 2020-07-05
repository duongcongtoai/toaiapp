package auth

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"toaiapp/registry"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"gopkg.in/yaml.v2"
)

const (
	ComponentName = "auth"
)

var (
	Component *AuthComponent

	PubKey  *rsa.PublicKey
	SignKey *rsa.PrivateKey
)

type AuthComponent struct {
	r             *registry.Registry
	UserGuest     *User
	drivers       map[string]DB
	driver        DB
	sessionStore  SessionStore
	sessionStores map[string]SessionStore
	config        Configuration
}

type Configuration struct {
	Debug                bool   `yaml:"debug"`
	JWTKeyFile           string `yaml:"jwt_key_file"`
	JWTPubKeyFile        string `yaml:"jwt_pub_key_file"`
	JWTExpirationSeconds int64  `yaml:"jwt_expiration_seconds"`
	DBType               string `yaml:"auth_db_type"`
	SessionStoreType     string `yaml:"session_store_type"`
}

func (c *AuthComponent) SetupFromYaml(configFile string) error {
	type YamlContainer struct {
		Configuration `yaml:"auth"`
	}
	conf := YamlContainer{}

	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(yamlFile, &conf); err != nil {
		return fmt.Errorf("Failed to parse section '%s': %v", ComponentName, err)
	}

	c.config = conf.Configuration

	if err := c.setupDB(configFile); err != nil {
		return err
	}

	if err := c.setupJwt(configFile); err != nil {
		return err
	}

	if err := c.setupSessionStore(configFile); err != nil {
		return err
	}

	return nil

}

func (c *AuthComponent) setupDB(configFile string) error {
	dbType := c.config.DBType
	driver, ok := c.drivers[dbType]
	if !ok {
		log.Fatalf("Unknown db backend '%s' for auth", dbType)
	}

	c.driver = driver
	if err := c.driver.SetupFromYamlConfig(configFile); err != nil {
		log.Fatalf("%v\n", err)
	}
	// if err := c.driver.Initialize(); err != nil {
	// 	log.Fatalf("%v\n", err)
	// }
	return nil
}

func (c *AuthComponent) setupJwt(configFile string) error {
	if !filepath.IsAbs(c.config.JWTKeyFile) {
		c.config.JWTKeyFile = filepath.Join(filepath.Dir(configFile), c.config.JWTKeyFile)
	}

	if !filepath.IsAbs(c.config.JWTPubKeyFile) {
		c.config.JWTPubKeyFile = filepath.Join(filepath.Dir(configFile), c.config.JWTPubKeyFile)
	}
	signBytes, err := ioutil.ReadFile(c.config.JWTKeyFile)
	if err != nil {
		log.Fatalf("Error: Cannot read private keyfile %v", c.config.JWTKeyFile)
	}
	SignKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Fatalf("Error parsing private key file: %v", c.config.JWTKeyFile)
	}

	pubBytes, err := ioutil.ReadFile(c.config.JWTPubKeyFile)
	if err != nil {
		log.Fatalf("Error: Cannot read public keyfile %v", c.config.JWTPubKeyFile)
	}
	PubKey, err = jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	if err != nil {
		log.Fatalf("Error: Cannot parse public keyfile %v", c.config.JWTPubKeyFile)
	}
	return nil
}

func (c *AuthComponent) setupSessionStore(configFile string) error {
	if c.config.SessionStoreType == "" {
		return nil
	}
	storeType := c.config.SessionStoreType
	store, ok := c.sessionStores[storeType]
	if !ok {
		log.Fatalf("Unknown session store '%s' for auth", storeType)
	}

	c.sessionStore = store
	if err := c.sessionStore.SetupFromYamlConfig(configFile); err != nil {

		log.Fatalf("%v\n", err)
	}
	return nil
}

func (c *AuthComponent) GetConfig() Configuration {
	return c.config
}

func (c *AuthComponent) RegisterDB(dbName string, db DB) {
	c.drivers[dbName] = db
}

func (c *AuthComponent) RegisterSessionStore(storeName string, store SessionStore) {
	c.sessionStores[storeName] = store
}

// func (c *AuthComponent) GetDriver() DB {
// 	return c.driver
// }

func (c *AuthComponent) SetupEcho(e *echo.Echo) error {
	return nil
	// return c.driver.SetupEcho(e)
}

func (c *AuthComponent) GetWeight() int {
	return 100
}

func (c *AuthComponent) Shutdown() error {
	return nil
}

func (c *AuthComponent) SetRegistry(r *registry.Registry) {
	c.r = r
}

func (c *AuthComponent) GetName() string {
	return "Auth"
}

func init() {
	Component = &AuthComponent{drivers: make(map[string]DB), sessionStores: make(map[string]SessionStore)}
	registry.Instance().Register(ComponentName, Component)
}
