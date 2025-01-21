package spectrum

type User struct {
	UserID        string
	Nickname      string
	Color         string
	currentRoomID string
}

func NewUser(userID string) *User {
	return &User{
		UserID: userID,
	}
}

func (u *User) SetNickname(nickname string) {
	u.Nickname = nickname
}

func (u *User) SetRoom(roomID string) {
	u.currentRoomID = roomID
}

func (u *User) Room() string {
	return u.currentRoomID
}

func (u *User) SetColor(color string) {
	u.Color = color
}

func (u *User) IsInRoom() bool {
	return u.currentRoomID != ""
}
