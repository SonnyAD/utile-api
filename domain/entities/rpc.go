package entities

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"utile.space/api/domain/valueobjects"
)

var (
	r = regexp.MustCompile(`^(commit|emoji|giveup|hit|join|lose|miss|proof|prove|shoot|signin|startgame)(\s+([0-9a-f-]*))?(\s+([0-9]+,[0-9]+))?(\s+([\x{1F600}-\x{1F6FF}|[\x{2600}-\x{26FF}]|[\x{1FAE3}]|[\x{1F92F}]|[\x{1FAE1}]|[\x{1F6DF}]))?$`)
)

var (
	ErrCommandNotRecognized = errors.New("command not recognized")
	ErrCannotReachOpponent  = errors.New("cannot reach opponent")
	ErrUnexpected           = errors.New("unexpected error")
)

//nolint:gocyclo
func (c *Client) EvaluateRPC(command string) error {
	subMatch := r.FindStringSubmatch(command)
	if subMatch == nil {
		return ErrCommandNotRecognized
	}

	hub := c.Hub

	switch {
	case subMatch[1] == "commit":
		index, err := strconv.Atoi(strings.Split(subMatch[5], ",")[0])
		if err != nil {
			return errors.Join(ErrUnexpected, err)
		}
		/*c.Send <- valueobjects.RPC_ACK.Export()
		err := hub.MessageOpponent(c.PlayerID, c.CurrentMatchID, command)
		hub.PlayerCommit(c.CurrentMatchID, c.PlayerID, subMatch[3], index)*/
		err = hub.MessageOpponent(c.PlayerID, c.CurrentMatchID, command)
		if err != nil {
			c.Send <- valueobjects.RPC_NACK.Export()
			return ErrCannotReachOpponent
		}
		hub.PlayerCommit(c.CurrentMatchID, c.PlayerID, subMatch[3], index)
		c.Send <- valueobjects.RPC_ACK.Export()
	case subMatch[1] == "emoji":
		err := hub.MessageOpponent(c.PlayerID, c.CurrentMatchID, "receive "+subMatch[7])
		if err != nil {
			c.Send <- valueobjects.RPC_NACK.Export()
			return ErrCannotReachOpponent
		}
	case subMatch[1] == "giveup":
		hub.EndMatch(c.PlayerID, "gaveup")
	case subMatch[1] == "hit":
		err := hub.MessageOpponent(c.PlayerID, c.CurrentMatchID, command)
		if err != nil {
			c.Send <- valueobjects.RPC_NACK.Export()
			return ErrCannotReachOpponent
		}
		c.Send <- []byte("turn")
	case subMatch[1] == "join":
		matchID, err := hub.QuickMatch(c.PlayerID)
		if err != nil {
			c.Send <- valueobjects.RPC_NACK.Export()
		}
		c.CurrentMatchID = matchID
		c.Send <- valueobjects.RPC_ACK.Export()
	case subMatch[1] == "lose":
		hub.EndMatch(c.PlayerID, "lose")
	case subMatch[1] == "miss":
		err := hub.MessageOpponent(c.PlayerID, c.CurrentMatchID, command)
		if err != nil {
			c.Send <- valueobjects.RPC_NACK.Export()
			return ErrCannotReachOpponent
		}
		c.Send <- []byte("turn")
	case subMatch[1] == "proof":
		err := hub.MessageOpponent(c.PlayerID, c.CurrentMatchID, command)
		if err != nil {
			c.Send <- valueobjects.RPC_NACK.Export()
			return ErrCannotReachOpponent
		}
	case subMatch[1] == "prove":
		err := hub.MessageOpponent(c.PlayerID, c.CurrentMatchID, command)
		if err != nil {
			c.Send <- valueobjects.RPC_NACK.Export()
			return ErrCannotReachOpponent
		}
	case subMatch[1] == "shoot":
		err := hub.MessageOpponent(c.PlayerID, c.CurrentMatchID, "shot "+subMatch[5])
		if err != nil {
			c.Send <- valueobjects.RPC_NACK.Export()
			return ErrCannotReachOpponent
		}
	case subMatch[1] == "signin":
		c.SetPlayerID(subMatch[3])
		hub.RecordPlayerIDClientMapping(c.PlayerID, c)
		c.Send <- valueobjects.RPC_ACK.Export()
	case subMatch[1] == "startgame":
		matchID := hub.NewMatch(c.PlayerID)
		c.CurrentMatchID = matchID
		c.Send <- valueobjects.RPC_ACK.Export()
	default:
		return ErrCommandNotRecognized
	}

	return nil
}
