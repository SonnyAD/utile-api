package entities

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"utile.space/api/domain/valueobjects"
)

var (
	r = regexp.MustCompile(`^(commit|emoji|giveup|hit|join|lose|miss|proof|prove|shoot|signin|nickname|startgame|startspectrum|joinspectrum|update|claim)(\s+([0-9a-f-]*))?(\s+([0-9]+,[0-9]+))?(\s+([\x{1F600}-\x{1F6FF}|[\x{2600}-\x{26FF}]|[\x{1FAE3}]|[\x{1F92F}]|[\x{1FAE1}]|[\x{1F6DF}]))?(\s+(.+))?$`)
)

var (
	ErrCommandNotRecognized = errors.New("command not recognized")
	ErrCannotReachOpponent  = errors.New("cannot reach opponent")
	ErrCannotParseCoords    = errors.New("cannot parse coords")
	ErrUnexpected           = errors.New("unexpected error")
)

//nolint:gocyclo
func (c *Client) EvaluateRPC(command string) error {
	subMatch := r.FindStringSubmatch(command)
	if subMatch == nil {
		return errors.Join(ErrCommandNotRecognized, errors.New(command))
	}

	log.Debug("RPC " + subMatch[0])

	hub := c.Hub

	switch {
	case subMatch[1] == "commit":
		spt := strings.Split(subMatch[5], ",")
		if len(spt) != 2 {
			return ErrCannotParseCoords
		}
		x, err := strconv.Atoi(spt[0])
		if err != nil {
			return errors.Join(ErrCannotParseCoords, err)
		}
		y, err := strconv.Atoi(spt[1])
		if err != nil {
			return errors.Join(ErrCannotParseCoords, err)
		}
		index := (x - 1) + (y-1)*10

		err = hub.MessageOpponent(c.PlayerID, hub.players[c.PlayerID].CurrentMatchID, command)
		if err != nil {
			c.Send <- valueobjects.RPC_NACK.Export()
			return errors.Join(ErrCannotReachOpponent, err)
		}
		hub.PlayerCommit(hub.players[c.PlayerID].CurrentMatchID, c.PlayerID, subMatch[3], index)
		c.Send <- valueobjects.RPC_ACK.Export()
	case subMatch[1] == "emoji":
		err := hub.MessageOpponent(c.PlayerID, hub.players[c.PlayerID].CurrentMatchID, "receive "+subMatch[7])
		if err != nil {
			c.Send <- valueobjects.RPC_NACK.Export()
			return errors.Join(ErrCannotReachOpponent, err)
		}
	case subMatch[1] == "giveup":
		hub.EndMatch(c.PlayerID, "gaveup")
	case subMatch[1] == "hit":
		err := hub.MessageOpponent(c.PlayerID, hub.players[c.PlayerID].CurrentMatchID, command)
		if err != nil {
			c.Send <- valueobjects.RPC_NACK.Export()
			return errors.Join(ErrCannotReachOpponent, err)
		}
		c.Send <- []byte("turn")
	case subMatch[1] == "join":
		matchID, err := hub.QuickMatch(c.PlayerID)
		if err != nil {
			c.Send <- valueobjects.RPC_NACK.Export()
		}
		hub.players[c.PlayerID].CurrentMatchID = matchID
		c.Send <- valueobjects.RPC_ACK.Export()
	case subMatch[1] == "lose":
		hub.EndMatch(c.PlayerID, "lose")
	case subMatch[1] == "miss":
		err := hub.MessageOpponent(c.PlayerID, hub.players[c.PlayerID].CurrentMatchID, command)
		if err != nil {
			c.Send <- valueobjects.RPC_NACK.Export()
			return errors.Join(ErrCannotReachOpponent, err)
		}
		c.Send <- []byte("turn")
	case subMatch[1] == "proof":
		err := hub.MessageOpponent(c.PlayerID, hub.players[c.PlayerID].CurrentMatchID, command)
		if err != nil {
			c.Send <- valueobjects.RPC_NACK.Export()
			return errors.Join(ErrCannotReachOpponent, err)
		}
	case subMatch[1] == "prove":
		err := hub.MessageOpponent(c.PlayerID, hub.players[c.PlayerID].CurrentMatchID, command)
		if err != nil {
			c.Send <- valueobjects.RPC_NACK.Export()
			return errors.Join(ErrCannotReachOpponent, err)
		}
	case subMatch[1] == "shoot":
		err := hub.MessageOpponent(c.PlayerID, hub.players[c.PlayerID].CurrentMatchID, "shot "+subMatch[5])
		if err != nil {
			c.Send <- valueobjects.RPC_NACK.Export()
			return errors.Join(ErrCannotReachOpponent, err)
		}
	case subMatch[1] == "signin":
		c.SetPlayerID(subMatch[3])
		hub.RecordPlayer(c.PlayerID, c)
		c.Send <- valueobjects.RPC_ACK.Export()
	case subMatch[1] == "nickname":
		c.Send <- valueobjects.RPC_ACK.Export()
		hub.players[c.PlayerID].SetNickname(subMatch[9])
	case subMatch[1] == "startgame":
		matchID := hub.NewMatch(c.PlayerID)
		hub.players[c.PlayerID].CurrentMatchID = matchID
		c.Send <- valueobjects.RPC_ACK.Export()
	case subMatch[1] == "startspectrum":
		matchID, password := hub.NewPrivateMatch(c.PlayerID)
		hub.players[c.PlayerID].CurrentMatchID = matchID
		c.Send <- []byte("spectrum " + matchID + " " + password)
	case subMatch[1] == "joinspectrum":
		spt := strings.Split(subMatch[9], " ")
		hub.players[c.PlayerID].SetNickname(spt[1])
		matchID, err := hub.JoinPrivateMatch(c.PlayerID, spt[0])
		if err != nil {
			// Nothing
		} else {
			hub.players[c.PlayerID].CurrentMatchID = matchID
			c.Send <- []byte("spectrum " + matchID + " " + subMatch[9])
		}
	case subMatch[1] == "update":
		if hub.players[c.PlayerID].CurrentMatchID != "" {
			hub.MessagePlayersInMatch(hub.players[c.PlayerID].CurrentMatchID, command)
		}
	case subMatch[1] == "claim":
		if hub.players[c.PlayerID].CurrentMatchID != "" {
			hub.MessagePlayersInMatch(hub.players[c.PlayerID].CurrentMatchID, command)
		}
	default:
		return ErrCommandNotRecognized
	}

	return nil
}
