package entities

import (
	"errors"
	"math/rand"
)

const (
	commitmentSize = 100
)

type MatchedPlayer struct {
	playerID    string
	commitments []string
}

func NewMatchedPlayer(playerID string) *MatchedPlayer {
	return &MatchedPlayer{
		playerID:    playerID,
		commitments: make([]string, commitmentSize),
	}
}

func (mp *MatchedPlayer) IsPlayerReady() bool {
	for _, commitment := range mp.commitments {
		if commitment == "" {
			return false
		}
	}
	return len(mp.commitments) == 100
}

type Match struct {
	players     []*MatchedPlayer
	matchOver   bool
	player1Turn bool
	turn        int
}

func (m *Match) PlayerCommit(playerID string, commit string, index int) error {
	for _, player := range m.players {
		if player.playerID == playerID {
			player.commitments[index] = commit
			return nil
		}
	}

	return errors.New("unknown player")
}

func (m *Match) HasPlayer1Won() bool {
	return m.matchOver && m.player1Turn
}

func (m *Match) HasPlayer2Won() bool {
	return m.matchOver && !m.player1Turn
}

func (m *Match) IsPendingPlayer() bool {
	return m.players[1] == nil
}

func (m *Match) Player2Join(player2 string) error {
	if !m.IsPendingPlayer() {
		return errors.New("cannot reassign a player")
	}
	m.players[1] = NewMatchedPlayer(player2)
	return nil
}

func NewPendingMatch(player1ID string) *Match {
	player1 := NewMatchedPlayer(player1ID)
	return &Match{
		players:     []*MatchedPlayer{player1, nil},
		matchOver:   false,
		player1Turn: rand.Intn(2) == 1, // randomly decide on whose turn
		turn:        1,
	}
}

func NewMatch(player1ID string, player2ID string) *Match {
	player1 := NewMatchedPlayer(player1ID)
	player2 := NewMatchedPlayer(player2ID)
	return &Match{
		players:     []*MatchedPlayer{player1, player2},
		matchOver:   false,
		player1Turn: rand.Intn(2) == 1, // randomly decide on whose turn
		turn:        1,
	}
}
