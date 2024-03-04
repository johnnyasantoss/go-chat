package shared

import (
	"context"
	"fmt"
	"log"
)

type Room struct {
	Name         string
	joiningQueue chan *User
	leavingQueue chan *User
	users        []*User
}

func NewRoom(name string) *Room {
	return &Room{
		Name:         name,
		joiningQueue: make(chan *User, 1),
		leavingQueue: make(chan *User, 1),
		users:        make([]*User, 0),
	}
}

func (r *Room) GetUsers() []*User {
	users := make([]*User, len(r.users))

	copy(users, r.users)

	return users
}

func (r *Room) AddUser(user *User) {
	log.Println("Adding user to room", r.Name, user.Name)
	r.joiningQueue <- user
}

func (r *Room) RemoveUser(user *User) {
	log.Println("Removing user from room", r.Name, user.Name)
	r.leavingQueue <- user
}

func (r *Room) IndexOf(user *User) int {
	for i, u := range r.users {
		if u == user {
			return i
		}
	}
	return -1
}

func (r *Room) Broadcast(sender *User, msg []byte) {
	for _, user := range r.GetUsers() {
		if sender == user {
			continue
		}

		msg := fmt.Sprintf("%s: %s", sender.Name, msg)

		user.Inbox <- msg
	}
}

func (r *Room) printCurrentUsers() {
	curUsers := r.GetUsers()
	log.Println("Current users in room:", len(curUsers), "Room:", r.Name)

	for _, user := range curUsers {
		log.Println("User:", user.Name, "Room:", r.Name)
	}
}

func (r *Room) Serve(ctx context.Context) {
	log.Println("Opening room", r.Name)

	for {
		select {
		case user := <-r.joiningQueue:
			r.users = append(r.users, user)
			r.printCurrentUsers()
		case user := <-r.leavingQueue:
			index := r.IndexOf(user)

			if index == -1 {
				continue
			}

			// remove user from Users slice
			copy(r.users[index:], r.users[index+1:])

			r.users = r.users[:len(r.users)-1]
			r.printCurrentUsers()
		case <-ctx.Done():
			log.Println("Closing the room", r.Name)

			close(r.joiningQueue)
			close(r.leavingQueue)

			r.users = nil

			return
		}
	}
}
