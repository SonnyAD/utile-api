// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package entities

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"utile.space/api/domain/valueobjects"
)

// Hub maintains the set of active clients/matches and broadcasts messages to the clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	mappingPlayerIDToClient map[string]*Client

	matches map[string]*Match

	messages chan *valueobjects.Message

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		messages:                make(chan *valueobjects.Message),
		Register:                make(chan *Client),
		unregister:              make(chan *Client),
		clients:                 make(map[*Client]bool),
		mappingPlayerIDToClient: make(map[string]*Client),
		matches:                 make(map[string]*Match),
	}
}

func (h *Hub) CountOnlinePlayers() int {
	return len(h.clients)
}

func (h *Hub) CountTotalMatches() int {
	return len(h.matches)
}

func (h *Hub) CountPendingMatches() int {
	var count = 0
	for _, match := range h.matches {
		if match.IsPendingPlayer() {
			count = count + 1
		}
	}
	return count
}

func (h *Hub) CountOngoingMatches() int {
	var count = 0
	for _, match := range h.matches {
		if !match.IsPendingPlayer() && !match.matchOver {
			count = count + 1
		}
	}
	return count
}

func (h *Hub) CountFinishedMatches() int {
	var count = 0
	for _, match := range h.matches {
		if match.matchOver {
			count = count + 1
		}
	}
	return count
}

func (h *Hub) RecordPlayerIDClientMapping(playerId string, client *Client) {
	h.mappingPlayerIDToClient[playerId] = client
}

func (h *Hub) NewMatch(player1 string) string {
	matchID := uuid.NewString()
	match := NewPendingMatch(player1)

	h.matches[matchID] = match

	return matchID
}

func (h *Hub) JoinMatch(matchID string, player2 string) {
	if err := h.matches[matchID].Player2Join(player2); err != nil {
		fmt.Println(err)
	}

	h.MessagePlayer(h.matches[matchID].players[0].playerID, matchID, "joined")
	h.MessagePlayer(h.matches[matchID].players[1].playerID, matchID, "youjoined")
}

func (h *Hub) MessagePlayer(senderID string, recipentID string, content string) {
	h.messages <- valueobjects.NewMessage(senderID, recipentID, []byte(content))
}

func (h *Hub) MessagePlayersInMatch(matchID string, content string) {
	for _, player := range h.matches[matchID].players {
		h.messages <- valueobjects.NewServiceMessage(player.playerID, []byte(content))
	}
}

func (h *Hub) MessageOpponent(playerID string, matchID string, message string) error {
	var opponentID string

	if h.matches[matchID] == nil {
		return errors.New("match not found")
	}

	if h.matches[matchID].players[1] != nil {
		if playerID == h.matches[matchID].players[0].playerID {
			opponentID = h.matches[matchID].players[1].playerID
		} else {
			opponentID = h.matches[matchID].players[0].playerID
		}
	} else {
		return errors.New("no opponent yet")
	}

	if h.mappingPlayerIDToClient[opponentID] == nil {
		return errors.New("player not found")
	}

	h.MessagePlayer(playerID, opponentID, message)

	return nil
}

func (h *Hub) EndMatch(player string, reason string) {
	matchID := h.mappingPlayerIDToClient[player].CurrentMatchID
	h.matches[matchID].matchOver = true

	// The first player who gives up lose
	h.matches[matchID].player1Turn = (player == h.matches[matchID].players[0].playerID)

	h.mappingPlayerIDToClient[player].CurrentMatchID = ""

	h.MessageOpponent(player, matchID, reason)
}

func (h *Hub) PlayerCommit(matchID string, playerID string, commit string, index int) {
	err := h.matches[matchID].PlayerCommit(playerID, commit, index)
	if err == nil && h.matches[matchID].players[0].IsPlayerReady() && h.matches[matchID].players[1].IsPlayerReady() {
		h.MessagePlayersInMatch(matchID, "battlestart")

		time.Sleep(time.Second)
		if h.matches[matchID].player1Turn {
			h.MessagePlayer("", h.matches[matchID].players[0].playerID, "turn")
		} else {
			h.MessagePlayer("", h.matches[matchID].players[1].playerID, "turn")
		}
	}
}

func (h *Hub) QuickMatch(player2 string) (string, error) {
	for matchID, match := range h.matches {
		if match.IsPendingPlayer() {
			match.Player2Join(player2)
			h.MessagePlayer("", match.players[0].playerID, "joined")
			h.MessagePlayer("", match.players[1].playerID, "youjoined")
			return matchID, nil
		}
	}

	fmt.Println("no match found")

	return "", errors.New("no match found")
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case client := <-h.Register:
			h.clients[client] = true
			log.WithFields(log.Fields{
				"onlinePlayerCount": len(h.clients),
			}).Debug("New player connected")
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.mappingPlayerIDToClient, client.PlayerID)
				close(client.Send)
			}
		case message := <-h.messages:
			if message.IsBroadcastMessage() {
				for client := range h.clients {
					select {
					case client.Send <- message.Content():
					default:
						delete(h.clients, client)
						delete(h.mappingPlayerIDToClient, client.PlayerID)
						close(client.Send)
					}
				}
			} else {
				if client, ok := h.mappingPlayerIDToClient[message.Recipient()]; ok {
					client.Send <- message.Content()
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
