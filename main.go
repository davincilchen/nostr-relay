package main

import (
	"nostr-relay/pkg/app/server"
	"nostr-relay/pkg/config"
)

const confPath = "./config.json"

func main() {

	cfg, err := config.New(confPath)
	if err != nil {
		//log.Fatal(err)
		cfg = &config.Config{}
		cfg.Server.Port = "8100"
	}

	svr := server.New(cfg)
	svr.Serve()

}
