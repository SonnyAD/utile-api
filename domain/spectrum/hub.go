package spectrum

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"
	"utile.space/api/domain/valueobjects"
	"utile.space/api/utils"
)

/*type room interface {
	setID(id string)
	setPassword(password string)
	AddUser(u user) error
	RemoveUser(u user) error
	Password() string
	Close()
	IsClosed() bool
	Users() []user
}

type client interface {
	UserID() string
	Send(content []byte)
}

type user interface {
	IsInRoom() bool
	SetRoom(roomId string)
	Room() string
	SetNickname(nickname string)
	Nickname() string
	UserID() string
	setUserID(userID string)
}*/

// Hub maintains the set of active clients with their business entity logic plus the entities associating clients together: Players with Battleships Matches, Participants with Spectrum Rooms, etc.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	users map[string]*User

	mappingUserIDToClient map[string]*Client

	rooms map[string]*Room

	messages chan *valueobjects.Message

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		messages:              make(chan *valueobjects.Message),
		Register:              make(chan *Client),
		unregister:            make(chan *Client),
		clients:               make(map[*Client]bool),
		users:                 make(map[string]*User),
		mappingUserIDToClient: make(map[string]*Client),
		rooms:                 make(map[string]*Room),
	}
}

func (h *Hub) CountOnlineUsers() int {
	return len(h.users)
}

func (h *Hub) CountTotalRooms() int {
	return len(h.rooms)
}

func (h *Hub) CountActiveRooms() int {
	var count = 0
	for _, room := range h.rooms {
		if !room.IsClosed() {
			count = count + 1
		}
	}
	return count
}

func (h *Hub) LinkUserWithClient(userID string, client *Client) {
	if _, ok := h.users[userID]; !ok {
		user := NewUser(userID)
		h.users[userID] = user
	}

	h.mappingUserIDToClient[userID] = client
}

func (h *Hub) NewRoom(userIDs []string) (string, error) {
	roomID := utils.GenerateRandomString(4)

	room := NewRoom(h.users[userIDs[0]], roomID, "Hello")

	for _, userID := range userIDs {
		err := room.AddUser(h.users[userID])
		if err != nil {
			return "", errors.New("unknown problem at room creation")
		}
	}

	h.rooms[roomID] = room

	return roomID, nil
}

func (h *Hub) NewPrivateRoom(userIDs []string) (string, string, error) {
	roomID, err := h.NewRoom(userIDs)
	if err != nil {
		return "", "", errors.New("unknown problem at room creation")
	}
	password := utils.GenerateRandomString(12)

	err = h.rooms[roomID].SetPassword(password)
	if err != nil {
		return "", "", errors.New("unknown problem at room locking")
	}

	return roomID, password, nil
}

func (h *Hub) JoinRoom(roomID string, userID string) error {
	room := h.rooms[roomID]

	if room.IsClosed() {
		return errors.New("room already closed")
	}

	user := h.users[userID]

	if err := room.AddUser(user); err != nil {
		return errors.New("user cannot join room")
	}

	user.SetRoom(roomID)

	userNickname := user.Nickname
	if userNickname == "" {
		userNickname = userID
	}

	h.MessageRoom(roomID, "joined "+userNickname)

	return nil
}

func (h *Hub) JoinPrivateRoom(roomID string, userID string, password string) error {
	room := h.rooms[roomID]

	if password != room.password {
		return errors.New("wrong room or password")
	}

	return h.JoinRoom(roomID, userID)
}

func (h *Hub) Broadcast(senderID string, content string) {
	h.messages <- valueobjects.NewBroadcastMessage(senderID, []byte(content))
}

func (h *Hub) MessageUser(senderID string, recipentID string, content string) {
	h.messages <- valueobjects.NewMessage(senderID, recipentID, []byte(content))
}

func (h *Hub) MessageRoom(roomID string, content string) {
	for _, user := range h.rooms[roomID].participants {
		h.messages <- valueobjects.NewServiceMessage(user.UserID, []byte(content))
	}
}

func (h *Hub) Run(ctx context.Context) {
	log.Debug("Hub runner starting...")
	for {
		select {
		case client := <-h.Register:
			h.clients[client] = true
			log.WithFields(log.Fields{
				"connectionsOpened": len(h.clients),
			}).Debug("New user connected")
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				log.WithFields(log.Fields{
					"player": (*client).UserID(),
				}).Debug("Unregistering client")
				delete(h.clients, client)
				delete(h.mappingUserIDToClient, (*client).UserID())
			}
		case message := <-h.messages:
			if message.IsBroadcastMessage() {
				for client := range h.clients {
					(*client).Send(message.Content())
				}
			} else {
				if client, ok := h.mappingUserIDToClient[message.Recipient()]; ok {
					(*client).Send(message.Content())
				}
			}
		case <-ctx.Done():
			log.Info("Hub runner terminated...")
			return
		}
	}
}
