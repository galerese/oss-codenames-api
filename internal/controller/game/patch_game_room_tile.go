package game

import (
	"fmt"
	"strconv"

	gamed "galere.se/oss-codenames-api/internal/domain/game"
	"galere.se/oss-codenames-api/internal/response"
	"github.com/gin-gonic/gin"
)

type PatchGameRoomTileRequest struct {
	Pointed *bool `json:"pointed"`
	Guessed *bool `json:"guessed"`
}

func (c *Controller) PatchGameRoomTile(gc *gin.Context) {

	ctx := gc.Request.Context()

	//
	// Validate request context
	//

	room, session := c.EnsureSessionPlayerIsInPathRoomByName(gc)
	if c.HasErrors(gc) {
		return
	}

	//
	// Validate path parameters
	//

	tileId := gc.Param("tileId")
	if tileId == "" {
		c.APIError(gc, "Path parameter [tileId] is required", nil, 400)
		return
	}

	tileIdInt, err := strconv.Atoi(tileId)
	if err != nil || tileIdInt < 1 || tileIdInt > gamed.DefaultBoardTileCount {
		c.APIError(gc, fmt.Sprintf("Path parameter [tileId] must be a valid number between 1 and %d", gamed.DefaultBoardTileCount), err, 400)
		return
	}

	//
	// Validate body parameters
	//

	request := PatchGameRoomTileRequest{}
	if c.ParseBody(gc, &request) != nil {
		return
	}

	if request.Pointed == nil && request.Guessed == nil {
		c.APIError(gc, "You must set either 'pointed' or 'guessed' on this tile", nil, 400)
		return
	}

	if request.Guessed != nil && request.Pointed != nil {
		c.APIError(gc, "You can only set 'pointed' or 'guessed' at one time on a tile", nil, 400)
		return
	}

	if request.Guessed != nil && !*request.Guessed {
		c.APIError(gc, "You cannot unguess a tile!", nil, 400)
		return
	}

	//
	// Execute action :)
	//

	if request.Guessed != nil && *request.Guessed {
		room, err = c.service.GuessTile(ctx, room, session.Player, tileIdInt)
		if err != nil {
			c.APIError(gc, "Unexpected error while guessing tile", err, 500)
			return
		}
	}

	if request.Pointed != nil && *request.Pointed {
		room, err = c.service.PointTile(ctx, room, session.Player, tileIdInt)
		if err != nil {
			c.APIError(gc, "Unexpected error while pointing tile", err, 500)
			return
		}
	}

	if request.Pointed != nil && !*request.Pointed {
		room, err = c.service.UnpointTile(ctx, room, session.Player, tileIdInt)
		if err != nil {
			c.APIError(gc, "Unexpected error while unpointing tile", err, 500)
			return
		}
	}

	c.APIResponse(gc, response.NewGameRoomResponse(room), 200)

}
