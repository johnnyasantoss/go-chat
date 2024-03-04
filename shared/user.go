package shared

import "github.com/gorilla/websocket"

type User struct {
	Name  string
	Conn  *websocket.Conn
	Inbox chan string
}

func NewUser(name string, conn *websocket.Conn) *User {
	return &User{
		Name:  name,
		Conn:  conn,
		Inbox: make(chan string, 10),
	}
}
