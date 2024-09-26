package game

import (
	"galere.se/oss-codenames-api/internal/response"
	"github.com/gin-gonic/gin"
)

type PostGameRoomPlayerRequest struct {
	PlayerName string `json:"playerName"`
}

func (c *Controller) PostGameRoomPlayer(gc *gin.Context) {
	ctx := gc.Request.Context()

	//
	// Validate request context
	//

	session := c.EnsureSessionForRequest(gc)
	if c.HasErrors(gc) {
		return
	}

	if session.CurrentRoom != nil {
		c.APIError(gc, "You are playing in another room!", nil, 400)
		return
	}

	//
	// Validate path
	//

	room := c.EnsureGameRoomExistsByNameOnPath(gc)
	if c.HasErrors(gc) {
		return
	}

	//
	// Validate body parameters
	//

	request := PostGameRoomPlayerRequest{}
	if c.ParseBody(gc, &request) != nil || request.PlayerName == "" {
		c.APIError(gc, "A body parameter 'playerName' is required to indicate the player name", nil, 400)
		return
	}

	// No validation required for the player name other than it not being empty
	if request.PlayerName == "" {
		c.APIError(gc, "Body parameter 'playerName' cannot be empty", nil, 400)
		return
	}

	//
	// Execute action :)
	//

	session.Player.Name = request.PlayerName

	// Add the player to the game room
	room, err := c.service.AddPlayerToGameRoom(ctx, room, session.Player)
	if err != nil {
		c.APIError(gc, "Unexpected error when adding player to game room", err, 500)
		return
	}

	session.CurrentRoom = room

	// Save the session
	_, err = c.sessionService.SaveSession(ctx, session)
	if err != nil {
		c.APIError(gc, "Unexpected errror saving session with the updated player name", err, 500)
		return
	}

	// Return the updated room
	c.APIResponse(gc, response.NewGameRoomResponse(room), 200)

}
