package game

import (
	"galere.se/oss-codenames-api/internal/response"
	"github.com/gin-gonic/gin"
)

type PostGameRoomRequest struct {
	PlayerName string `json:"playerName"`
}

func (c *Controller) PostGameRoom(gc *gin.Context) {

	ctx := gc.Request.Context()

	//
	// Validate request context
	//

	session := c.EnsureSessionForRequest(gc)
	if c.HasErrors(gc) {
		return
	}

	if session.CurrentRoom != nil {
		c.APIError(gc, "You are already in a game room", nil, 400)
		return
	}

	//
	// Validate body parameters
	//

	request := PostGameRoomRequest{}
	err := c.ParseBody(gc, &request)
	if err != nil || request.PlayerName == "" {
		c.APIError(gc, "A body parameter 'playerName' is required to indicate the player name", err, 400)
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

	// Update session :)
	session.Player.Name = request.PlayerName

	// Create a game room
	room, err := c.service.CreateGameRoom(ctx, session.Player)
	if err != nil {
		c.APIError(gc, "Unexpected error creating a game room", err, 500)
		return
	}

	session.CurrentRoom = room

	// Save the session
	_, err = c.sessionService.SaveSession(ctx, session)
	if err != nil {
		c.APIError(gc, "Unexpected errror saving session with the updated player name", err, 500)
		return
	}

	c.APIResponse(gc, response.NewGameRoomResponse(room), 201)

}
