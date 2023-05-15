package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Server Server `json:"Server"`
	Nostr  Nostr  `json:"Nostr"`
	//DB     db.Config `json:"DB"`
}

type Server struct {
	Port string `json:"Port"`
	IP   string `json:"IP"`
}

type Nostr struct {
	PublicKey  string `json:"PublicKey"`
	PrivateKey string `json:"PrivateKey"` //TODO:
}

func New(path string) (*Config, error) {
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
