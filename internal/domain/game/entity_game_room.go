package game

import "time"

type GameRoom struct {
	Id           string
	CreatedAt    time.Time
	State        GameRoomState
	Name         string
	RedTeam      []*Player
	BlueTeam     []*Player
	Players      []*Player
	CurrentRound *GameRound
	RoundHistory []*GameRound
}

func (r *GameRoom) IsPlayerInRoom(player *Player) bool {

	for _, p := range r.Players {
		if p.Id == player.Id {
			return true
		}
	}

	return false
}

func (r *GameRoom) getPlayerTeam(player *Player) TeamName {

	for _, member := range r.RedTeam {
		if member.Id == player.Id {
			return TeamNameRed
		}
	}

	for _, member := range r.BlueTeam {
		if member.Id == player.Id {
			return TeamNameBlue
		}
	}

	return ""
}
