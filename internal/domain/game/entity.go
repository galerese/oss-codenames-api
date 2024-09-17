package game

import "time"

type GameRound struct {
	RedScore       Score
	BlueScore      Score
	CurrentTurn    *GameTurn
	StartedAt      time.Time
	FinishedAt     time.Time
	WinnerTeam     TeamName
	RedSpymaster   *Player
	BlueSpymaster  *Player
	RedScoreLimit  Score
	BlueScoreLimit Score
	Identities     map[int]Identity
	BoardTiles     map[int]BoardTile
	GuessedTiles   map[int]bool
	TurnHistory    []*GameTurn
	WinReason      WinReason
}

type BoardTile struct {
	ImageUrl string
}

type Player struct {
	Name  string
	IP    string
	Id    string
	Token string
}
