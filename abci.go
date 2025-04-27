package main

import (
	"log"
	"net"

	"github.com/tendermint/abci/types"
	"github.com/tendermint/abci/server"
)

type Application struct {
	types.BaseApplication
}

func main() {
	app := &Application{}

	// Set up ABCI server
	srv, err := server.NewServer("tcp://127.0.0.1:26658", "socket", app)
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}

	srv.Start()
	defer srv.Stop()

	log.Println("ABCI server running on tcp://127.0.0.1:26658")

	// Block forever
	select {}
}
