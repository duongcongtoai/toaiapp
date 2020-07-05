package psql

import (
	"encoding/gob"
	"io/ioutil"
	"toaiapp/auth"

	"github.com/antonlindstrom/pgstore"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"gopkg.in/yaml.v2"
)

type PGresSessionStore struct {
	config Configuration
	store  *pgstore.PGStore
}

type Configuration struct {
	DBURL      string `yaml:"db_url"`
	KeyFile    string `yaml:"session_key_file"`
	PubKeyFile string `yaml:"session_pub_key_file"`
}

func (p *PGresSessionStore) Store() sessions.Store {
	return p.store
}

func (p *PGresSessionStore) SetupFromYamlConfig(configFile string) error {
	configByte, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	type ConfigurationContainer struct {
		Configuration `yaml:"auth_session"`
	}
	var cfgContainer ConfigurationContainer
	if err := yaml.Unmarshal(configByte, &cfgContainer); err != nil {
		return err
	}
	p.config = cfgContainer.Configuration
	// if !filepath.IsAbs(p.config.KeyFile) {
	// 	p.config.KeyFile = filepath.Join(filepath.Dir(configFile), p.config.KeyFile)
	// }

	// if !filepath.IsAbs(p.config.PubKeyFile) {
	// 	p.config.PubKeyFile = filepath.Join(filepath.Dir(configFile), p.config.PubKeyFile)
	// }
	// signBytes, err := ioutil.ReadFile(p.config.KeyFile)
	// if err != nil {
	// 	log.Fatalf("Error: Cannot read private keyfile %v", p.config.KeyFile)
	// }
	// // SignKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	// // if err != nil {
	// // 	log.Fatalf("Error parsing private key file: %v", p.config.KeyFile)
	// // }

	// pubBytes, err := ioutil.ReadFile(p.config.PubKeyFile)
	// if err != nil {
	// 	log.Fatalf("Error: Cannot read public keyfile %v", p.config.PubKeyFile)
	// }
	// PubKey, err = jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	// if err != nil {
	// 	log.Fatalf("Error: Cannot parse public keyfile %v", p.config.PubKeyFile)
	// }
	store, err := pgstore.NewPGStore(p.config.DBURL, securecookie.GenerateRandomKey(16), securecookie.GenerateRandomKey(16))
	if err != nil {
		return err
	}
	p.store = store
	gob.Register(auth.User{})
	return nil
}

// func toByteSlice(in []string) [][]byte {
// 	out := make([][]byte, len(in))
// 	for i := range in {
// 		out[i] = []byte(in[i])
// 	}
// 	return out
// }

// func (p *PGresSessionStore) SetupEcho(e *echo.Echo) error {
// 	return nil
// 	// e.Use(psqlSessionMiddlewareFunc(p.config.))
// }
// func (p *PGresSessionStore) SetConfig(c auth.Configuration) {

// }
// func (p *PGresSessionStore) Initialize() error {

// }
// func (p *PGresSessionStore) Store() (sessions.Store, error) {

// }

// func (p *PGresSessionStore) StoreFromContext(c echo.Context) (sessions.Store, error) {

// }

func init() {
	auth.Component.RegisterSessionStore("postgresql", &PGresSessionStore{})
}
