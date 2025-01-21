package spectrum

import (
	"errors"
)

type Room struct {
	id           string
	topic        string
	password     string
	closed       bool
	admins       []string
	participants []*User
}

func (r *Room) Join(newUser *User) error {
	if r.closed {
		return errors.New("room closed")
	}
	newUser.SetRoom(r.password)
	return nil
}

func NewRoom(creator *User, id string, topic string) *Room {
	creator.SetRoom(id)
	return &Room{
		topic:        topic,
		closed:       false,
		admins:       []string{creator.UserID},
		participants: make([]*User, 0, 10),
	}
}

func RemoveByPointer(slice []*User, ptr *User) []*User {
	for i, item := range slice {
		if item == ptr {
			// Remove the element by concatenating the slice
			return append(slice[:i], slice[i+1:]...)
		}
	}
	// Return the original slice if the element wasn't found
	return slice
}

func (r *Room) Leave(user *User) error {
	if r.closed {
		return errors.New("room closed")
	}
	r.participants = RemoveByPointer(r.participants, user)
	user.SetRoom("")
	return nil
}

func (r *Room) RoomID() string {
	return r.id
}

func (r *Room) SetPassword(password string) error {
	if r.closed {
		return errors.New("room closed")
	}
	r.password = password
	return nil
}

func (r *Room) AddUser(user *User) error {
	if r.closed {
		return errors.New("room closed")
	}
	r.participants = append(r.participants, user)
	return nil
}

func (r *Room) SetAdmin(user *User) error {
	if r.closed {
		return errors.New("room closed")
	}
	r.admins = append(r.admins, user.UserID)
	return nil
}

func (r *Room) Close() error {
	if r.closed {
		return errors.New("room closed")
	}
	r.closed = true
	return nil
}

func (r *Room) IsClosed() bool {
	return r.closed
}
