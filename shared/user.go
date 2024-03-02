package shared

import "github.com/gorilla/websocket"

type User struct {
	name string
	conn *websocket.Conn
}
