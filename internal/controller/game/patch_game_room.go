package game

import (
	"slices"
	"strings"

	gamed "galere.se/oss-codenames-api/internal/domain/game"
	"github.com/gin-gonic/gin"
)

type PatchGameRoomRequest struct {
	State string `json:"state"`
	// Only used when state is "SELECT_CLUE"
	Clue             string `json:"clue"`
	GuessAmount      int    `json:"guess_amount"`
	UnlimitedGuesses bool   `json:"unlimited_guesses"`
}

var validGameRoomTransitionStates = []string{
	gamed.GameRoomStateStarted,
	gamed.GameRoomStateSpymastersSettled,
	gamed.GameRoomStateClueSelected,
	gamed.GameRoomStateGuessingStopped,
}

// Note: you can only see a game room you are part of!
func (c *Controller) PatchGameRoom(gc *gin.Context) {

	room, session := c.EnsureSessionPlayerIsInPathRoomByName(gc)
	if c.HasErrors(gc) {
		return
	}

	//
	// Validate body parameters
	//

	request := PatchGameRoomRequest{}
	if c.ParseBody(gc, &request) != nil {
		return
	}

	if !slices.Contains(validGameRoomTransitionStates, request.State) {
		c.APIError(gc, "Body parameter 'state' is invalid: you can only change states to one of "+strings.Join(validGameRoomTransitionStates, ", "), nil, 400)
		return
	}

	//
	// Execute action :)
	//

	var err error

	switch request.State {
	// A game is starting
	case gamed.GameRoomStateStarted:
		room, err = c.service.StartGame(room, session.Player)
	// Spymasters have been selected
	case gamed.GameRoomStateSpymastersSettled:
		room, err = c.service.SettleSpymasters(room, session.Player)
	// A clue has been selected
	case gamed.GameRoomStateClueSelected:
		room, err = c.service.SelectClue(room, session.Player, gamed.SelectClueInput{
			Clue:             request.Clue,
			GuessAmount:      request.GuessAmount,
			UnlimitedGuesses: request.UnlimitedGuesses,
		})
	// Guessing has been stopped prematurely :)
	case gamed.GameRoomStateGuessingStopped:
		room, err = c.service.StopGuessing(room, session.Player)
	}

	if err != nil {
		gc.Error(err)
		return
	}

	c.APIResponse(gc, room, 200)

}
