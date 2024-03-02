package main

import (
	"context"
	"log"
	"time"

	"johnnyasantos.com/chat/server"
	"johnnyasantos.com/chat/shared"
)

func main() {
	log.SetFlags(log.Lmsgprefix | log.Ldate | log.Ltime)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	shared.HandleSignals(cancel)

	server.StartWsServer(ctx)

	time.Sleep(1 * time.Second)

	log.Println("Bye")
}
