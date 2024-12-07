package entities

type Player struct {
	PlayerID       string
	Nickname       string
	CurrentMatchID string
}

func NewPlayer(playerID string) *Player {
	return &Player{
		PlayerID: playerID,
	}
}

func (p *Player) SetNickname(nickname string) {
	p.Nickname = nickname
}

func (p *Player) SetMatch(matchID string) {
	p.CurrentMatchID = matchID
}

func (p *Player) IsInMatch() bool {
	return p.CurrentMatchID != ""
}
