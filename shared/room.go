package shared

type Room struct {
	Name    string
	Joining chan *User
}

func (*Room) GetUsers() []User {
	// TODO: impl
	return make([]User, 0)
}
