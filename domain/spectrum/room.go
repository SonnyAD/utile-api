package spectrum

import (
	"errors"
	"slices"
)

var generalColors = []string{
	"aeaeae", // Neutral gray
	"ff5555", // Bright red
	"cd5334", // Burnt oran*-ge
	"ff9955", // Vibrant orange
	"ffe680", // Soft yellow
	"aade87", // Light green
	"9fd8cb", // Pale teal
	"aaeeff", // Light cyan
	"c6afe9", // Soft lavender
	"985f6f", // Muted mauve
}

type Room struct {
	id           string
	topic        string
	password     string
	closed       bool
	admins       []string
	participants map[string]*User
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
		participants: make(map[string]*User),
	}
}

func (r *Room) Leave(user *User) error {
	if r.closed {
		return errors.New("room closed")
	}

	for i, item := range r.participants {
		if item == user {
			delete(r.participants, i)
			break
		}
	}

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

func (r *Room) AddUser(color string, user *User) error {
	if r.closed {
		return errors.New("room closed")
	}

	if !slices.Contains(generalColors, color) {
		return errors.New("unknown color")
	}

	if _, alreadyPresent := r.participants[color]; alreadyPresent {
		return errors.New("color already taken")
	}

	r.participants[color] = user
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
