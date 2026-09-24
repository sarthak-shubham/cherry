package main

import (
	"log"

	"github.com/sarthak-shubham/cherry/recovery"
	"github.com/sarthak-shubham/cherry/server"
	"github.com/sarthak-shubham/cherry/store"
)

const (
	serverAddress = ":8080"
	walPath       = "data.jsonl"
)

func main() {
	data, err := recovery.Recover(walPath)
	if err != nil {
		log.Fatal(err)
	}

	kvStore := store.NewStore(data, walPath)

	httpServer := server.NewServer(kvStore)

	err = httpServer.Start(serverAddress)
	if err != nil {
		log.Fatal(err)
	}
}
