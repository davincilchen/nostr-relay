package main

import (
	"nostr-relay/pkg/app/server"
	"nostr-relay/pkg/config"
)

const confPath = "./config.json"

func main() {

	//cfg, err := config.New(confPath)
	cfg, _ := config.New(confPath)

	svr := server.New(cfg)
	svr.Serve()

}
