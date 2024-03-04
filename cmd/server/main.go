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

	lobby := shared.NewRoom("Lobby")
	chatCtx := server.NewChatContext(ctx, lobby)

	go lobby.Serve(ctx)

	server.StartWsServer(chatCtx)

	time.Sleep(1 * time.Second)

	log.Println("Bye")
}
