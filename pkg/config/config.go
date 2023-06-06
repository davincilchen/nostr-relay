package config

import (
	"encoding/json"
	"io/ioutil"
	"nostr-relay/pkg/db"

	"github.com/caarlos0/env"
	_ "github.com/joho/godotenv/autoload" //support .env && autoload
)

type Config struct {
	Server Server    `json:"Server"`
	DB     db.Config `json:"DB"`
}

type Server struct {
	Port string `json:"Port" env:"SERVER_PORT" envDefault:"8100"`
}

var config *Config

func GetConfig() *Config {
	return config
}

func New(path string) (*Config, error) {

	cfg, err := new(path)

	if cfg != nil {
		config = cfg
	}

	return cfg, err

}

func new(path string) (*Config, error) {

	cfg, err := NewFromFile(path)
	if err == nil { //priority first
		return cfg, nil
	}

	cfg, _ = NewFromEnv()

	return cfg, nil
}

func NewFromFile(path string) (*Config, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(buf, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func NewFromEnv() (*Config, error) {

	var config Config

	config.Server = *GetServerConfig()
	config.DB = *GetDBConfig()

	return &config, nil
}

func GetServerConfig() *Server {
	cfg := &Server{}
	env.Parse(cfg)
	return cfg
}

func GetDBConfig() *db.Config {
	cfg := &db.Config{}
	env.Parse(cfg)
	return cfg
}
