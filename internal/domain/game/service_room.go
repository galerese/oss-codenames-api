package game

import (
	"context"
	"fmt"
	"time"

	"galere.se/oss-codenames-api/pkg/domain_util"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Gets a get game room by its name
func (s *Service) GetGameRoomByName(ctx context.Context, name string) (*GameRoom, error) {
	return s.repository.GetGameRoomByName(ctx, name)
}

// Creates a game room with a random name
func (s *Service) CreateGameRoom(ctx context.Context, actor *Player) (*GameRoom, error) {
	s.baseService.Logger.Infof("Creating new game room for player [%s]", actor.Name)

	room := &GameRoom{
		Id:           uuid.NewString(),
		Name:         s.generateRandomGameRoomName(),
		CreatedAt:    time.Now(),
		State:        GameRoomStateWaitingForPlayers,
		Players:      []*Player{actor},
		RedTeam:      []*Player{},
		BlueTeam:     []*Player{},
		CurrentRound: nil,
		RoundHistory: []*GameRound{},
	}

	err := s.repository.SaveGameRoom(ctx, room)
	if err != nil {
		return nil, errors.Wrap(err, "failed to save new game room")
	}

	s.baseService.Logger.Infof("Created new game room [%s] for player [%s]", room.Name, actor.Name)

	return room, nil
}

// Settle teams and start the game!
func (s *Service) StartGame(ctx context.Context, room *GameRoom, actor *Player) (*GameRoom, error) {
	s.baseService.Logger.Infof("Starting game in room [%s] for player [%s]", room.Id, actor.Name)

	// Validation

	if room.State != GameRoomStateWaitingForPlayers {
		return nil, domain_util.NewStateValidationError("You may not start the game because it has already started!")
	}

	if len(room.RedTeam) < 2 || len(room.BlueTeam) < 2 {
		return nil, domain_util.NewStateValidationError("There must be at least 2 players on each team to start the game!")
	}

	// Logic

	room.State = GameRoomStateSelectSpymasters
	room.CurrentRound = s.createEmptyGameRound()

	err := s.repository.SaveGameRoom(ctx, room)
	if err != nil {
		return nil, domain_util.NewUnexpectedError(err, "failed to save game room for starting the game")
	}

	// Wrap up

	s.triggerGameRoomEvents(room, GameRoomEventGameStarted)
	s.triggerGameRoomEvents(room, GameRoomEventRoundStarted)

	s.baseService.Logger.Infof("Started game in room [%s] for player [%s], state [%s]", room.Id, actor.Name, room.State)

	return room, nil
}

// Spymasters have been selected, start the game!
func (s *Service) SettleSpymasters(ctx context.Context, room *GameRoom, actor *Player) (*GameRoom, error) {
	s.baseService.Logger.Infof("Spymasters selection indicated for game in room [%s] by player [%s]", room.Id, actor.Name)

	// Validation

	if room.State != GameRoomStateSelectSpymasters {
		return nil, domain_util.NewStateValidationError("A game round has already begun, you cannot change spymasters before it ends!")
	}

	if room.CurrentRound == nil {
		return nil, domain_util.NewUnexpectedError(nil, "Expected a game round to be created already before selecting spymasters!")
	}

	if room.CurrentRound.RedSpymaster != nil && room.CurrentRound.BlueSpymaster != nil {
		return nil, domain_util.NewUnexpectedError(nil, "Please wait until both spymasters have been selected to start the game!")
	}

	// Logic

	room.State = GameRoomStateSelectClue

	var err error
	room.CurrentRound, err = s.initializeGameRound(ctx, room.CurrentRound)
	if err != nil {
		return nil, domain_util.NewUnexpectedError(err, "Failed to initialize game round after spymasters have been selected :(")
	}

	err = s.repository.SaveGameRoom(ctx, room)
	if err != nil {
		return nil, domain_util.NewUnexpectedError(err, "failed to save game room for spymaster selection")
	}

	// Wrap up

	s.triggerGameRoomEvents(room, GameRoomEventSpymasterSelected)

	s.baseService.Logger.Infof("Spymasters selection indicated for game in room [%s] by player [%s], state [%s]", room.Id, actor.Name, room.State)

	return room, nil
}

type SelectClueInput struct {
	Clue             string
	GuessAmount      int
	UnlimitedGuesses bool
}

// A clue has been selected, start the guessing!
func (s *Service) SelectClue(ctx context.Context, room *GameRoom, actor *Player, input SelectClueInput) (*GameRoom, error) {
	s.baseService.Logger.Infof("Selecting clue for game in room [%s] by player [%s]", room.Id, actor.Name)

	// Validation

	if room.State != GameRoomStateSelectClue {
		return nil, domain_util.NewStateValidationError("It's not currently time to select a clue!")
	}

	if room.CurrentRound == nil {
		return nil, domain_util.NewUnexpectedError(nil, "Expected a game round to be created already before selecting a clue!")
	}

	if room.CurrentRound.CurrentTurn == nil {
		return nil, domain_util.NewUnexpectedError(nil, "Expected a game turn to be created already before selecting a clue!")
	}

	// Validate the current team's spymaster is the actor

	currentSpymaster := room.CurrentRound.RedSpymaster
	if room.CurrentRound.CurrentTurn.CurrentTeam == TeamNameBlue {
		currentSpymaster = room.CurrentRound.BlueSpymaster
	}

	if currentSpymaster.Id != actor.Id {
		return nil, domain_util.NewInvalidActionError(fmt.Sprintf("Only the current team's spymaster can select the clue! That is '%s'.", currentSpymaster.Name))
	}

	// Validate the clue is not empty

	if len(input.Clue) == 0 || len(input.Clue) > MaxClueLength {
		return nil, domain_util.NewInvalidParameterError(fmt.Sprintf("You must select a clue with up to %d characters!", MaxClueLength))
	}

	if !input.UnlimitedGuesses && (input.GuessAmount < 1 || input.GuessAmount > MaxGuessAmount) {
		return nil, domain_util.NewInvalidParameterError(fmt.Sprintf("You must select a guess amount between 1 and %d!", MaxGuessAmount))
	}

	// Logic

	currentTurn := room.CurrentRound.CurrentTurn
	currentTurn.Clue = input.Clue
	currentTurn.GuessAmount = input.GuessAmount
	currentTurn.UnlimitedGuesses = input.UnlimitedGuesses

	room.State = GameRoomStateSelectGuess

	err := s.repository.SaveGameRoom(ctx, room)
	if err != nil {
		return nil, domain_util.NewUnexpectedError(err, "failed to save game room for selecting clue")
	}

	// Wrap up

	s.triggerGameRoomEvents(room, GameRoomEventClueSelected)

	s.baseService.Logger.Infof("Selected clue [%s] with [%d] guesses and unlimited guesses [%t] for game in room [%s] by player [%s], state [%s]",
		currentTurn.Clue, currentTurn.GuessAmount, currentTurn.UnlimitedGuesses, room.Id, actor.Name, room.State)

	return room, nil

}

// Guessing has been stopped prematurely, proceed :)
func (s *Service) StopGuessing(ctx context.Context, room *GameRoom, actor *Player) (*GameRoom, error) {

	s.baseService.Logger.Infof("Stopping guessing for game in room [%s] by player [%s]", room.Id, actor.Name)

	// Validation

	if room.State != GameRoomStateSelectGuess {
		return nil, domain_util.NewStateValidationError("You cannot currently stop guessing because the game is not in the guessing phase!")
	}

	if room.CurrentRound == nil {
		return nil, domain_util.NewUnexpectedError(nil, "Expected a game round to be created already before stopping guessing!")
	}

	if room.CurrentRound.CurrentTurn == nil {
		return nil, domain_util.NewUnexpectedError(nil, "Expected a game turn to be created already before stopping guessing!")
	}

	// Logic

	room, err := s.startNewGameTurn(room)
	if err != nil {
		return nil, err
	}

	err = s.repository.SaveGameRoom(ctx, room)
	if err != nil {
		return nil, domain_util.NewUnexpectedError(err, "failed to save game room for stopping guessing")
	}

	// Wrap up

	s.triggerGameRoomEvents(room, GameRoomEventGuessingStopped)

	s.baseService.Logger.Infof("Stopped guessing for game in room [%s] by player [%s], state [%s]", room.Id, actor.Name, room.State)

	return room, nil

}
