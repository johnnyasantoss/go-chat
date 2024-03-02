package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: true,
	HandshakeTimeout:  3 * time.Second,
	CheckOrigin:       func(r *http.Request) bool { return true },
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func serveWs(w http.ResponseWriter, r *http.Request, server *http.Server) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Failed to upgrade ws connection", err)
		fmt.Fprint(w, "Failed to upgrade connection")

		return
	}

	server.RegisterOnShutdown(func() {
		log.Println("Closing connection on client", conn.RemoteAddr())

		msg := websocket.FormatCloseMessage(websocket.CloseGoingAway, "bye")
		err := conn.WriteControl(websocket.CloseMessage, msg, time.Now().Add(time.Second))

		if err != nil {
			log.Println("Failed to close connection on client", conn.RemoteAddr(), err)
		}
	})
}

func StartWsServer(ctx context.Context) {
	log.Println("Starting WebSocket Server on :1337")

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:              ":1337",
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           mux,
		BaseContext:       func(l net.Listener) context.Context { return context.WithValue(ctx, "test", "test") },
	}

	mux.HandleFunc("/", serveHome)
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r, server)
	})

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	<-ctx.Done()

	log.Println("Server is closing...")
	shutdownCtx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()
	err := server.Shutdown(shutdownCtx)

	if err != nil {
		log.Fatal("Failed to cleaning shutdown server", err)
	}

	<-shutdownCtx.Done()
}
