package main

import (
    "log"
    "github.com/tendermint/tendermint/abci/server"
    abci "github.com/tendermint/tendermint/abci/types"
)

type MyApplication struct {
    abci.BaseApplication
}

func main() {
    app := &MyApplication{}

    // Create a socket server listening at tcp://127.0.0.1:26658
    s := server.NewSocketServer("tcp://127.0.0.1:26658", app)
    err := s.Start()
    if err != nil {
        log.Fatalf("Failed to start ABCI server: %v", err)
    }

    log.Println("ABCI server running at tcp://127.0.0.1:26658")

    // Block forever
    select {}
}
