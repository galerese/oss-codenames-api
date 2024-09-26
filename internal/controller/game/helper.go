package game

import (
	"galere.se/oss-codenames-api/internal/domain/game"
	"galere.se/oss-codenames-api/internal/domain/session"
	"github.com/gin-gonic/gin"
)

func (c *Controller) EnsureSessionHasGameRoom(gc *gin.Context) (*game.GameRoom, *session.Session) {

	session := c.EnsureSessionForRequest(gc)
	if c.HasErrors(gc) {
		return nil, nil
	}

	if session.CurrentRoom == nil {
		c.APIError(gc, "You are not in a game room", nil, 400)
		return nil, nil
	}

	room, err := c.service.GetGameRoomByName(gc.Request.Context(), session.CurrentRoom.Name)
	if err != nil {
		c.APIError(gc, "Unexpected error when getting game room by name", err, 500)
		return nil, nil
	}

	return room, session
}

func (c *Controller) EnsureSessionPlayerIsInPathRoomByName(gc *gin.Context) (*game.GameRoom, *session.Session) {

	room, session := c.EnsureSessionHasGameRoom(gc)
	if c.HasErrors(gc) {
		return nil, nil
	}

	//
	// Validate path parameters
	//

	roomName := gc.Param("name")
	if roomName == "" {
		c.APIError(gc, "Path parameter [name] representing the room name is required", nil, 400)
		return nil, nil
	}

	if room.Name != roomName {
		c.APIError(gc, "You can only act on the room you are currently in", nil, 400)
		return nil, nil
	}

	return room, session

}

func (c *Controller) GetGameRoomByNameFromPath(gc *gin.Context) *game.GameRoom {
	roomName := gc.Param("name")
	if roomName == "" {
		c.APIError(gc, "Path parameter [name] representing the room name is required", nil, 400)
		return nil
	}

	room, err := c.service.GetGameRoomByName(gc.Request.Context(), roomName)
	if err != nil {
		c.APIError(gc, "Unexpected error when getting game room by name", err, 500)
		return nil
	}

	return room
}

// Ensures that a game room exists by name on the path
func (c *Controller) EnsureGameRoomExistsByNameOnPath(gc *gin.Context) *game.GameRoom {
	room := c.GetGameRoomByNameFromPath(gc)
	if room == nil {
		c.APIError(gc, "That room doesn't seem to exist!", nil, 404)
		return nil
	}

	return room
}

func (c *Controller) EnsurePlayerIdMatchesSessionPlayer(gc *gin.Context, session *session.Session) {

	playerId := gc.Param("playerId")

	if playerId == "" {
		c.APIError(gc, "Path parameter [playerId] is required", nil, 400)
		return
	}

	if playerId != session.Player.Id {
		c.APIError(gc, "You are not allowed to perform any actions in another player's behalf", nil, 403)
		return
	}

}
