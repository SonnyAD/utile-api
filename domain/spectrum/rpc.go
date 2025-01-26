package spectrum

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"

	log "github.com/sirupsen/logrus"
	"utile.space/api/domain/valueobjects"
)

const spectrum = "spectrum "

var (
	r = regexp.MustCompile(`^(emoji|signin|nickname|startspectrum|joinspectrum|leavespectrum|resetpositions|update|claim)(\s+([0-9a-f-]*))?(\s+([0-9]+,[0-9]+))?(\s+([\x{1F600}-\x{1F6FF}|[\x{2600}-\x{26FF}]|[\x{1FAE3}]|[\x{1F92F}]|[\x{1FAE1}]|[\x{1F6DF}]))?(\s+(.+))?$`)
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

	switch {
	case subMatch[1] == "emoji":
		c.hub.MessageRoom(c.hub.users[c.UserID()].Room(), "receive "+subMatch[7])
	case subMatch[1] == "signin":
		c.SetUserID(subMatch[3])
		c.hub.LinkUserWithClient(c.UserID(), c)
		c.send <- valueobjects.RPC_ACK.Export()
		if c.hub.users[c.userID].IsInRoom() {
			roomID := c.hub.users[c.userID].currentRoomID
			admin := slices.Contains(c.hub.rooms[roomID].admins, c.userID)
			c.send <- []byte(spectrum + c.hub.users[c.userID].currentRoomID + " " + fmt.Sprintf("%t", admin))
		}
	case subMatch[1] == "nickname":
		c.send <- valueobjects.RPC_ACK.Export()
		c.hub.users[c.UserID()].SetNickname(subMatch[9])
	case subMatch[1] == "startspectrum":
		roomID, err := c.hub.NewRoom(c.UserID(), subMatch[3])
		if err != nil {
			c.send <- valueobjects.RPC_NACK.Export()
			break
		}
		c.hub.users[c.UserID()].SetRoom(roomID)
		c.send <- []byte(spectrum + roomID)
	case subMatch[1] == "joinspectrum":
		spt := strings.Split(subMatch[9], " ")
		roomID := spt[0]
		c.hub.users[c.UserID()].SetNickname(spt[1])
		c.hub.users[c.UserID()].SetColor(spt[2])
		err := c.hub.JoinRoom(roomID, c.UserID(), spt[2])
		if err != nil {
			// Nothing
			log.Error(err.Error())
		} else {
			c.hub.users[c.UserID()].SetRoom(roomID)
			c.send <- []byte(spectrum + roomID + " " + subMatch[9])
		}
	case subMatch[1] == "leavespectrum":
		roomID := c.hub.users[c.userID].currentRoomID
		c.hub.users[c.userID].SetRoom("")
		err := c.hub.rooms[roomID].Leave(c.hub.users[c.userID])
		if err != nil {
			c.send <- valueobjects.RPC_NACK.Export()
			break
		}
		c.send <- valueobjects.RPC_ACK.Export()
	case subMatch[1] == "update":
		if c.hub.users[c.UserID()].IsInRoom() {
			c.hub.MessageRoom(c.hub.users[c.UserID()].Room(), command)
		}
	case subMatch[1] == "resetpositions":
		if c.hub.users[c.UserID()].IsInRoom() {
			newPositions := []string{"405,383", "376,413", "322,421", "323,381", "279,389", "360,381"}
			room := c.hub.rooms[c.hub.users[c.UserID()].Room()]
			var i = 0
			for _, user := range room.participants {
				if slices.Contains(room.admins, user.UserID) {
					continue
				}
				c.hub.MessageUser(c.UserID(), user.UserID, "newposition "+newPositions[i%len(newPositions)])
				i = i + 1
			}
		}
	case subMatch[1] == "claim":
		if c.hub.users[c.UserID()].IsInRoom() {
			c.hub.MessageRoom(c.hub.users[c.UserID()].Room(), command)
		}
	default:
		return ErrCommandNotRecognized
	}

	return nil
}
