package game

import (
	"galere.se/oss-codenames-api/internal/response"
	"github.com/gin-gonic/gin"
)

// Note: you can only see a game room you are part of!
func (c *Controller) GetGameRoom(gc *gin.Context) {

	session := c.EnsureSessionForRequest(gc)
	if session == nil {
		return
	}

	if session.CurrentRoom == nil {
		c.APIError(gc, "You can only see details from a room you are currently in!", nil, 403)
		return
	}

	c.APIResponse(gc, response.NewGameRoomResponse(session.CurrentRoom), 200)

}
