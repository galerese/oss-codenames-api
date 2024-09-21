package game

import (
	"context"
	"math/rand"
	"time"

	"github.com/pkg/errors"
)

func (s *Service) createEmptyGameRound() *GameRound {
	return &GameRound{
		RedScore:      0,
		BlueScore:     0,
		CurrentTurn:   nil,
		TurnHistory:   []*GameTurn{},
		StartedAt:     time.Now(),
		WinnerTeam:    "",
		RedSpymaster:  nil,
		BlueSpymaster: nil,
		Identities:    nil,
		BoardTiles:    nil,
		GuessedTiles:  map[int]bool{},
		WinReason:     "",
	}
}

// This sets up the game round with the correct number of tiles and identities (by default 25 tiles)
func (s *Service) initializeGameRound(ctx context.Context, round *GameRound) (*GameRound, error) {

	// Get random board tiles
	tiles, err := s.repository.GetRandomBoardTiles(ctx, DefaultBoardTileCount) // At some point maybe we make this configurable
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get random board tiles")
	}
	round.BoardTiles = tiles

	// Create the first turn :)
	round.CurrentTurn = s.createEmptyGameTurn()
	round.CurrentTurn.CurrentTeam = []TeamName{TeamNameRed, TeamNameBlue}[rand.Intn(2)]

	return round, nil
}
