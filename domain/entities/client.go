package entities

import (
	"bytes"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var lastClientId = 0

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	id int

	PlayerID string

	CurrentMatchID string
}

func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	lastClientId = lastClientId + 1
	return &Client{
		Hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
		id:   lastClientId,
	}
}

func (c *Client) SetPlayerID(playerID string) {
	c.PlayerID = playerID
}

func (c *Client) SetPlayerInMatchID(matchID string) {
	c.CurrentMatchID = matchID
}

func (c *Client) IsInMatch() bool {
	return c.CurrentMatchID != ""
}

// ReadPump pumps messages from the websocket connection to the hub.
//
// The application runs ReadPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		r := regexp.MustCompile(`^(signin|startgame|join|shoot|giveup|commit|emoji|miss|hit|prove|proof|lose)(\s+([0-9a-f-]*))?(\s+([0-9]+,[0-9]+))?(\s+([\x{1F600}-\x{1F6FF}|[\x{2600}-\x{26FF}]|[\x{1FAE3}]|[\x{1F92F}]|[\x{1FAE1}]|[\x{1F6DF}]))?$`)
		subMatch := r.FindStringSubmatch(string(message))

		ack := []byte("ack")
		nack := []byte("nack")
		/*hit := []byte("hit")
		miss := []byte("miss")*/

		if subMatch != nil {
			if subMatch[1] == "signin" {
				c.SetPlayerID(subMatch[3])
				c.Hub.mappingPlayerIDToClient[c.PlayerID] = c
				c.send <- ack
			} else if subMatch[1] == "startgame" {
				matchID := c.Hub.NewMatch(c.PlayerID)
				c.CurrentMatchID = matchID
				c.send <- ack
			} else if subMatch[1] == "join" {
				matchID, err := c.Hub.QuickMatch(c.PlayerID)
				if err == nil {
					c.CurrentMatchID = matchID
					c.send <- ack
				} else {
					c.send <- nack
				}
			} else if subMatch[1] == "commit" {
				index, _ := strconv.Atoi(strings.Split(subMatch[5], ",")[0])
				c.send <- ack
				c.Hub.MessageOpponent(c.PlayerID, c.CurrentMatchID, string(message))
				c.Hub.PlayerCommit(c.CurrentMatchID, c.PlayerID, subMatch[3], index)
			} else if subMatch[1] == "shoot" {
				c.Hub.MessageOpponent(c.PlayerID, c.CurrentMatchID, "shot "+subMatch[5])
			} else if subMatch[1] == "emoji" {
				c.Hub.MessageOpponent(c.PlayerID, c.CurrentMatchID, "receive "+subMatch[7])
			} else if subMatch[1] == "giveup" {
				c.Hub.EndMatch(c.PlayerID, "gaveup")
			} else if subMatch[1] == "lose" {
				c.Hub.EndMatch(c.PlayerID, "lose")
			} else if subMatch[1] == "hit" {
				c.Hub.MessageOpponent(c.PlayerID, c.CurrentMatchID, string(message))
				c.send <- []byte("turn")
			} else if subMatch[1] == "miss" {
				c.Hub.MessageOpponent(c.PlayerID, c.CurrentMatchID, string(message))
				c.send <- []byte("turn")
			} else if subMatch[1] == "prove" {
				c.Hub.MessageOpponent(c.PlayerID, c.CurrentMatchID, string(message))
			} else if subMatch[1] == "proof" {
				c.Hub.MessageOpponent(c.PlayerID, c.CurrentMatchID, string(message))
			}

			if err != nil {
				log.Println("write:", err)
				continue
			}
		}
	}
}

// WritePump pumps messages from the hub to the websocket connection.
//
// A goroutine running WritePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
