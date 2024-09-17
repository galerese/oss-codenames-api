package game

import "time"

type GameTurn struct {
	CurrentTeam      TeamName
	StartedAt        time.Time
	FinishedAt       time.Time
	Clue             string
	GuessAmount      int
	GuessesMade      int
	GuessedTiles     map[int]bool
	PointedTiles     map[int]map[string]bool
	UnlimitedGuesses bool
}

func (t *GameTurn) OtherTeam() TeamName {
	if t.CurrentTeam == TeamNameBlue {
		return TeamNameRed
	}

	if t.CurrentTeam == TeamNameRed {
		return TeamNameBlue
	}

	return ""
}
