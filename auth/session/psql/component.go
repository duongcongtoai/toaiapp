package psql

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
	DBURL        string `yaml:"db_url"`
	HashKeyFile  string `yaml:"session_hash_key_file"`
	BlockKeyFile string `yaml:"session_block_key_file"`
}

func (p *PGresSessionStore) Store() sessions.Store {
	return p.store
}

func (p *PGresSessionStore) prepareKeyPairs(configFile string) error {
	if !filepath.IsAbs(p.config.HashKeyFile) {
		p.config.HashKeyFile = filepath.Join(filepath.Dir(configFile), p.config.HashKeyFile)
	}
	if _, err := os.Stat(p.config.HashKeyFile); os.IsNotExist(err) {
		hashKeyByte := securecookie.GenerateRandomKey(16)
		if err := ioutil.WriteFile(p.config.HashKeyFile, hashKeyByte, 0644); err != nil {
			return fmt.Errorf("Error writing hash key file %s", p.config.HashKeyFile)
		}
	}
	if !filepath.IsAbs(p.config.BlockKeyFile) {
		p.config.BlockKeyFile = filepath.Join(filepath.Dir(configFile), p.config.BlockKeyFile)
	}
	if _, err := os.Stat(p.config.BlockKeyFile); os.IsNotExist(err) {
		blockKeyByte := securecookie.GenerateRandomKey(16)
		if err := ioutil.WriteFile(p.config.BlockKeyFile, blockKeyByte, 0644); err != nil {
			return fmt.Errorf("Error writing hash key file %s", p.config.BlockKeyFile)
		}
	}
	return nil
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
	if err := p.prepareKeyPairs(configFile); err != nil {
		return err
	}
	hashByte, err := ioutil.ReadFile(p.config.HashKeyFile)
	if err != nil {
		log.Fatalf("Cannot read hash key file %s", p.config.HashKeyFile)
	}
	blockByte, err := ioutil.ReadFile(p.config.BlockKeyFile)
	if err != nil {
		log.Fatalf("Cannot read block key file %s", p.config.BlockKeyFile)
	}
	store, err := pgstore.NewPGStore(p.config.DBURL, hashByte, blockByte)
	if err != nil {
		return err
	}
	p.store = store
	gob.Register(auth.User{})
	return nil
}

func init() {
	auth.Component.RegisterSessionStore("postgresql", &PGresSessionStore{})
}
