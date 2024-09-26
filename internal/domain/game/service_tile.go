package game

import (
	"context"
	"time"

	"galere.se/oss-codenames-api/pkg/domain_util"
)

// PointTile points to a tile in the game room
func (s *Service) PointTile(ctx context.Context, room *GameRoom, actor *Player, tileId int) (*GameRoom, error) {
	s.baseService.Logger.Infof("Pointing tile [%d] for game in room [%s] by player [%s]", tileId, room.Id, actor.Name)

	// Validation
	if room.State != GameRoomStateSelectGuess {
		return nil, domain_util.NewStateValidationError("You can only point to tiles during your team's guessing phase!")
	}

	if room.CurrentRound == nil {
		return nil, domain_util.NewUnexpectedError(nil, "Expected a game round to be created already before pointing a tile!")
	}

	if room.CurrentRound.CurrentTurn == nil {
		return nil, domain_util.NewUnexpectedError(nil, "Expected a game turn to be created already before pointing a tile!")
	}

	// Logic

	if room.CurrentRound.CurrentTurn.PointedTiles[tileId] == nil {
		room.CurrentRound.CurrentTurn.PointedTiles[tileId] = make(map[string]bool)
	}

	// Nothing to do if the tile is already pointed by the player :)
	if room.CurrentRound.CurrentTurn.PointedTiles[tileId][actor.Id] {
		return room, nil
	}

	room.CurrentRound.CurrentTurn.PointedTiles[tileId][actor.Id] = true

	err := s.repository.SaveGameRoom(ctx, room)
	if err != nil {
		return nil, domain_util.NewUnexpectedError(err, "failed to save game room for pointing tile")
	}

	// Wrap up
	s.triggerGameRoomEvents(room, GameRoomEventTilePointed)

	s.baseService.Logger.Infof("Pointed tile [%d] for game in room [%s] by player [%s], state [%s]", tileId, room.Id, actor.Name, room.State)

	return room, nil
}

// UnpointTile removes the point from a tile in the game room
func (s *Service) UnpointTile(ctx context.Context, room *GameRoom, actor *Player, tileId int) (*GameRoom, error) {
	s.baseService.Logger.Infof("Unpointing tile [%d] for game in room [%s] by player [%s]", tileId, room.Id, actor.Name)

	// Validation
	if room.State != GameRoomStateSelectGuess {
		return nil, domain_util.NewStateValidationError("You can only unpoint a tile during your team's guessing phase!")
	}

	if room.CurrentRound == nil {
		return nil, domain_util.NewUnexpectedError(nil, "Expected a game round to be created already before unpointing a tile!")
	}

	if room.CurrentRound.CurrentTurn == nil {
		return nil, domain_util.NewUnexpectedError(nil, "Expected a game turn to be created already before unpointing a tile!")
	}

	// Logic

	// Nothing to do if the tile is not pointed by the player
	if room.CurrentRound.CurrentTurn.PointedTiles[tileId] == nil ||
		room.CurrentRound.CurrentTurn.PointedTiles[tileId][actor.Id] == false {
		return room, nil
	}

	delete(room.CurrentRound.CurrentTurn.PointedTiles[tileId], actor.Id)

	err := s.repository.SaveGameRoom(ctx, room)
	if err != nil {
		return nil, domain_util.NewUnexpectedError(err, "failed to save game room for unpointing tile")
	}

	// Wrap up
	s.triggerGameRoomEvents(room, GameRoomEventTilePointed)

	s.baseService.Logger.Infof("unpointed tile [%d] for game in room [%s] by player [%s], state [%s]", tileId, room.Id, actor.Name, room.State)

	return room, nil
}

// GuessTile guesses a tile in the game room
func (s *Service) GuessTile(ctx context.Context, room *GameRoom, actor *Player, tileId int) (*GameRoom, error) {
	s.baseService.Logger.Infof("Guessing tile [%d] for game in room [%s] by player [%s]", tileId, room.Id, actor.Name)

	// Validation
	if room.State != GameRoomStateSelectGuess {
		return nil, domain_util.NewStateValidationError("You can only make guesses during your team's guessing phase!")
	}

	if room.CurrentRound == nil {
		return nil, domain_util.NewUnexpectedError(nil, "Expected a game round to be created already before guessing a tile!")
	}

	if room.CurrentRound.CurrentTurn == nil {
		return nil, domain_util.NewUnexpectedError(nil, "Expected a game turn to be created already before guessing a tile!")
	}

	// Tile has already been guessed :o
	if room.CurrentRound.GuessedTiles[tileId] {
		return nil, domain_util.NewInvalidParameterError("This tile has already been guessed!")
	}

	// Logic
	currentRound := room.CurrentRound
	currentTurn := currentRound.CurrentTurn

	currentRound.GuessedTiles[tileId] = true
	currentTurn.GuessedTiles[tileId] = true
	currentTurn.GuessesMade++

	// Validate what to do based on the selected tile :)
	shouldStartNewTurn := false
	switch currentRound.Identities[tileId] {
	case IdentityBlue:
		currentRound.BlueScore++

		// Game has ended for team blue!
		if currentRound.BlueScore >= room.CurrentRound.BlueScoreLimit {
			currentRound.WinnerTeam = TeamNameBlue
			currentRound.WinReason = WinReasonScore
			break
		}

		if currentTurn.CurrentTeam != TeamNameBlue {
			shouldStartNewTurn = true
		}

	case IdentityRed:
		currentRound.RedScore++

		// Game has ended for team red!
		if currentRound.RedScore >= room.CurrentRound.RedScoreLimit {
			currentRound.WinnerTeam = TeamNameRed
			currentRound.WinReason = WinReasonScore
			break
		}

		if currentTurn.CurrentTeam != TeamNameRed {
			shouldStartNewTurn = true
		}

	case IdentityNeutral:
		shouldStartNewTurn = true

	case IdentityAssassin:
		currentRound.WinReason = WinReasonAssassin
		currentRound.WinnerTeam = currentTurn.OtherTeam()

	}

	// Team hasn't finished yet but guesses might be over
	if currentRound.WinnerTeam == "" && !shouldStartNewTurn &&
		(currentTurn.GuessesMade > currentTurn.GuessAmount && !currentTurn.UnlimitedGuesses) {
		shouldStartNewTurn = true
	}

	// Game has ended!
	if currentRound.WinnerTeam != "" {
		// Game has ended, act accordingly :)
		room.State = GameRoomStateRoundEnded
		room.CurrentRound.FinishedAt = time.Now()
		shouldStartNewTurn = false

	}

	// Starts a new turn if conditions are met
	if shouldStartNewTurn {
		var err error
		room, err = s.startNewGameTurn(room)
		if err != nil {
			return nil, err
		}
	}

	err := s.repository.SaveGameRoom(ctx, room)
	if err != nil {
		return nil, domain_util.NewUnexpectedError(err, "failed to save game room for tile guessing")
	}

	// Wrap up
	s.triggerGameRoomEvents(room, GameRoomEventTileGuessed)

	if shouldStartNewTurn || room.State == GameRoomStateRoundEnded {
		s.triggerGameRoomEvents(room, GameRoomEventTurnEnded)
	}

	if room.State == GameRoomStateRoundEnded {
		s.triggerGameRoomEvents(room, GameRoomEventRoundEnded)
	}

	s.baseService.Logger.Infof("Guessed tile [%s] for game in room [%s] by player [%s], state [%s]", tileId, room.Id, actor.Name, room.State)

	return room, nil
}
