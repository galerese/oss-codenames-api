package game

import "github.com/gin-gonic/gin"

type PostGameRoomRequest struct {
	PlayerName string `json:"playerName"`
}

func (c *Controller) PostGameRoom(gc *gin.Context) {

	//
	// Validate request context
	//

	session := c.EnsureSessionForRequest(gc)
	if c.HasErrors(gc) {
		return
	}

	if session.CurrentGameRoom != nil {
		c.APIError(gc, "You are already in a game room", nil, 400)
		return
	}

	//
	// Validate body parameters
	//

	request := PostGameRoomRequest{}
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

	// Save the session
	session, err := c.sessionService.SaveSession(session)
	if err != nil {
		c.APIError(gc, "Unexpected errror saving session with the updated player name", err, 500)
		return
	}

	// Create a game room
	room, err := c.service.CreateGameRoom(session.Player)
	if err != nil {
		c.APIError(gc, "Unexpected error creating a game room", err, 500)
		return
	}

	c.APIResponse(gc, room, 201)

}
