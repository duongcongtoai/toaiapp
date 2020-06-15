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
	r         *registry.Registry
	UserGuest *User
	drivers   map[string]AuthDB
	driver    AuthDB
	config    *Configuration
}

type Configuration struct {
	Debug                bool   `yaml:"debug"`
	JWTKeyFile           string `yaml:"jwt_key_file"`
	JWTPubKeyFile        string `yaml:"jwt_pub_key_file"`
	JWTExpirationSeconds int64  `yaml:"jwt_expiration_seconds"`
	AuthDBType           string `yaml:"auth_db_type"`
	AuthDBUrl            string `yaml:"auth_db_url"`
}

func (c *AuthComponent) SetupFromYaml(configFile string, debug bool) error {
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

	return c.SetupStruct(&conf.Configuration, configFile, debug)
}

func (c *AuthComponent) SetupStruct(cfg *Configuration, configFile string, debug bool) error {
	c.config = cfg
	c.config.Debug = debug
	var (
		driver AuthDB
		ok     bool
	)

	if driver, ok = c.drivers[c.config.AuthDBType]; !ok {
		log.Fatalf("Unknown db backend '%s' for auth", c.config.AuthDBType)
	}

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
		fmt.Println(err)
		log.Fatalf("Error: Cannot parse public keyfile %v", c.config.JWTPubKeyFile)
	}

	c.driver = driver
	c.driver.SetConfig(*c.config)
	db, err := driver.DB()
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	if err = db.Initialize(); err != nil {
		log.Fatalf("%v\n", err)
	}
	return nil
}
func (c *AuthComponent) GetConfig() *Configuration {
	return c.config
}

func (c *AuthComponent) RegisterDB(dbName string, db AuthDB) {
	c.drivers[dbName] = db
}

func (c *AuthComponent) GetDriver() AuthDB {
	return c.driver
}

func (c *AuthComponent) SetupEcho(e *echo.Echo) error {
	return c.driver.SetupEcho(e)
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
	Component = &AuthComponent{drivers: make(map[string]AuthDB)}
	registry.Instance().Register(ComponentName, Component)
}
