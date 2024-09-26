package game

import (
	"slices"
	"strings"

	gamed "galere.se/oss-codenames-api/internal/domain/game"
	"galere.se/oss-codenames-api/internal/response"
	"github.com/gin-gonic/gin"
)

type PatchGameRoomRequest struct {
	State string `json:"state"`
	// Only used when state is "SELECT_CLUE"
	Clue             string `json:"clue"`
	GuessAmount      int    `json:"guessAmount"`
	UnlimitedGuesses bool   `json:"unlimitedGuesses"`
}

var validGameRoomTransitionStates = []string{
	gamed.GameRoomStateStarted,
	gamed.GameRoomStateSpymastersSettled,
	gamed.GameRoomStateClueSelected,
	gamed.GameRoomStateGuessingStopped,
}

// Note: you can only see a game room you are part of!
func (c *Controller) PatchGameRoom(gc *gin.Context) {

	ctx := gc.Request.Context()

	room, session := c.EnsureSessionPlayerIsInPathRoomByName(gc)
	if c.HasErrors(gc) {
		return
	}

	//
	// Validate body parameters
	//

	request := PatchGameRoomRequest{}
	err := c.ParseBody(gc, &request)
	if err != nil {
		c.APIError(gc, "A request body with a valid 'state' parameter is required", err, 400)
		return
	}

	if !slices.Contains(validGameRoomTransitionStates, request.State) {
		c.APIError(gc, "Body parameter 'state' is invalid: you can only change states to one of "+strings.Join(validGameRoomTransitionStates, ", "), nil, 400)
		return
	}

	//
	// Execute action :)
	//

	switch request.State {
	// A game is starting
	case gamed.GameRoomStateStarted:
		room, err = c.service.StartGame(ctx, room, session.Player)
	// Spymasters have been selected
	case gamed.GameRoomStateSpymastersSettled:
		room, err = c.service.SettleSpymasters(ctx, room, session.Player)
	// A clue has been selected
	case gamed.GameRoomStateClueSelected:
		room, err = c.service.SelectClue(ctx, room, session.Player, gamed.SelectClueInput{
			Clue:             request.Clue,
			GuessAmount:      request.GuessAmount,
			UnlimitedGuesses: request.UnlimitedGuesses,
		})
	// Guessing has been stopped prematurely :)
	case gamed.GameRoomStateGuessingStopped:
		room, err = c.service.StopGuessing(ctx, room, session.Player)
	}

	if err != nil {
		gc.Error(err)
		return
	}

	c.APIResponse(gc, response.NewGameRoomResponse(room), 200)

}
