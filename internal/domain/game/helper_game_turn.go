package game

import (
	"fmt"
	"time"

	"galere.se/oss-codenames-api/pkg/domain_util"
)

func (s *Service) startNewGameTurn(room *GameRoom) (*GameRoom, error) {

	room.State = GameRoomStateSelectClue

	// Adjust last turn :)
	lastTurn := room.CurrentRound.CurrentTurn
	lastTurn.FinishedAt = time.Now()

	// Create a new turn
	newTurn := s.createEmptyGameTurn()
	newTurn.CurrentTeam = lastTurn.OtherTeam()

	if newTurn.CurrentTeam == "" {
		return nil, domain_util.NewUnexpectedError(nil, fmt.Sprintf("Unexpected team name for new turn: %s", lastTurn.CurrentTeam))
	}

	// Set new turn and adjust history :)
	room.CurrentRound.CurrentTurn = newTurn
	room.CurrentRound.TurnHistory = append(room.CurrentRound.TurnHistory, lastTurn)

	return room, nil
}

func (s *Service) createEmptyGameTurn() *GameTurn {
	return &GameTurn{
		StartedAt:    time.Now(),
		GuessedTiles: make(map[int]bool),
		PointedTiles: make(map[int]map[string]bool),
	}
}
