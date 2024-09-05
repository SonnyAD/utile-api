package entities

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"utile.space/api/domain/valueobjects"
)

var (
	r = regexp.MustCompile(`^(signin|startgame|join|shoot|giveup|commit|emoji|miss|hit|prove|proof|lose)(\s+([0-9a-f-]*))?(\s+([0-9]+,[0-9]+))?(\s+([\x{1F600}-\x{1F6FF}|[\x{2600}-\x{26FF}]|[\x{1FAE3}]|[\x{1F92F}]|[\x{1FAE1}]|[\x{1F6DF}]))?$`)
)

var (
	ErrCommandNotRecognized = errors.New("command not recognized")
)

func (c *Client) EvaluateRPC(command string) error {
	subMatch := r.FindStringSubmatch(command)
	if subMatch == nil {
		return ErrCommandNotRecognized
	}

	hub := c.Hub

	switch {
	case subMatch[1] == "signin":
		c.SetPlayerID(subMatch[3])
		hub.RecordPlayerIDClientMapping(c.PlayerID, c)
		c.Send <- valueobjects.RPC_ACK.Export()
	case subMatch[1] == "startgame":
		matchID := hub.NewMatch(c.PlayerID)
		c.CurrentMatchID = matchID
		c.Send <- valueobjects.RPC_ACK.Export()
	case subMatch[1] == "join":
		matchID, err := hub.QuickMatch(c.PlayerID)
		if err == nil {
			c.CurrentMatchID = matchID
			c.Send <- valueobjects.RPC_ACK.Export()
		} else {
			c.Send <- valueobjects.RPC_NACK.Export()
		}
	case subMatch[1] == "commit":
		index, _ := strconv.Atoi(strings.Split(subMatch[5], ",")[0])
		c.Send <- valueobjects.RPC_ACK.Export()
		hub.MessageOpponent(c.PlayerID, c.CurrentMatchID, command)
		hub.PlayerCommit(c.CurrentMatchID, c.PlayerID, subMatch[3], index)
	case subMatch[1] == "shoot":
		hub.MessageOpponent(c.PlayerID, c.CurrentMatchID, "shot "+subMatch[5])
	case subMatch[1] == "emoji":
		hub.MessageOpponent(c.PlayerID, c.CurrentMatchID, "receive "+subMatch[7])
	case subMatch[1] == "giveup":
		hub.EndMatch(c.PlayerID, "gaveup")
	case subMatch[1] == "lose":
		hub.EndMatch(c.PlayerID, "lose")
	case subMatch[1] == "hit":
		hub.MessageOpponent(c.PlayerID, c.CurrentMatchID, command)
		c.Send <- []byte("turn")
	case subMatch[1] == "miss":
		hub.MessageOpponent(c.PlayerID, c.CurrentMatchID, command)
		c.Send <- []byte("turn")
	case subMatch[1] == "prove":
		hub.MessageOpponent(c.PlayerID, c.CurrentMatchID, command)
	case subMatch[1] == "proof":
		hub.MessageOpponent(c.PlayerID, c.CurrentMatchID, command)
	default:
		return ErrCommandNotRecognized
	}

	return nil
}
