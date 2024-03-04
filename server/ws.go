package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"johnnyasantos.com/chat/shared"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: true,
	HandshakeTimeout:  3 * time.Second,
	CheckOrigin:       func(r *http.Request) bool { return true },
}

func serveWs(resW http.ResponseWriter, req *http.Request, server *http.Server) {
	conn, err := upgrader.Upgrade(resW, req, nil)
	if err != nil {
		log.Println("Failed to upgrade ws connection", err)
		http.Error(resW, "Failed to upgrade ws connection", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	userName := "TODO"
	userNameHeader := req.Header["X-UserName"]
	if len(userNameHeader) > 0 {
		userName = userNameHeader[0]
	}

	user := shared.NewUser(userName, conn)

	ctx, cancel := context.WithCancel(req.Context())
	defer cancel()
	server.RegisterOnShutdown(cancel)
	closedByUser := false

	chatCtx, ok := ChatContextFromContext(ctx)

	if !ok {
		log.Fatalf("Invalid context on ws upgrade: %T\n", req.Context())
		http.Error(resW, "Internal error", http.StatusInternalServerError)
		return
	}

	room := chatCtx.Room

	connCloseHandler := conn.CloseHandler()
	conn.SetCloseHandler(func(code int, text string) error {
		closedByUser = true
		cancel()
		return connCloseHandler(code, text)
	})

	room.AddUser(user)
	defer room.RemoveUser(user)

	userMessages := readUserMessages(ctx, conn)

userLoop:
	for {
		select {
		case <-ctx.Done():
			break userLoop
		case message, ok := <-userMessages:
			if !ok {
				break userLoop
			}

			log.Printf("Message: %s: %s\n", user.Name, message)

			room.Broadcast(user, message)
		case msg := <-user.Inbox:
			if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
				log.Println("Failed to send message to user", user.Name)
				break userLoop
			}
		}
	}

	log.Println("Closing connection on client", conn.RemoteAddr())

	if closedByUser {
		return
	}

	msg := websocket.FormatCloseMessage(websocket.CloseGoingAway, "bye")

	if err := conn.WriteControl(websocket.CloseMessage, msg, time.Now().Add(time.Second)); err != nil {
		log.Println("Failed to close connection on client", conn.RemoteAddr(), err)
	}
}

func readUserMessages(ctx context.Context, conn *websocket.Conn) <-chan []byte {
	userMessages := make(chan []byte, 1)

	go func() {
		for {
			_, message, err := conn.ReadMessage()

			if ctx.Err() != nil || err != nil {
				close(userMessages)
				break
			}

			userMessages <- message
		}
	}()

	return userMessages
}

func StartWsServer(ctx context.Context) {
	log.Println("Starting WebSocket Server on :1337")

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:              ":1337",
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           mux,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/ws":
			serveWs(w, r, server)
		default:
			http.Error(w, "Use the chat client", http.StatusBadRequest)
		}
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
