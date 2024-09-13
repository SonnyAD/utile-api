package entities

type Player struct {
	PlayerID       string
	CurrentMatchID string
}

func NewPlayer(playerID string) *Player {
	return &Player{
		PlayerID: playerID,
	}
}

func (p *Player) SetMatch(matchID string) {
	p.CurrentMatchID = matchID
}

func (p *Player) IsInMatch() bool {
	return p.CurrentMatchID != ""
}
