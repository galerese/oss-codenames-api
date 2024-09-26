package game

import (
	"fmt"
	"slices"
	"strings"

	"galere.se/oss-codenames-api/internal/domain/game"
	"galere.se/oss-codenames-api/internal/response"
	"github.com/gin-gonic/gin"
)

var validTeamNames = []game.TeamName{game.TeamNameBlue, game.TeamNameRed}
var validTeamNamesString = strings.Join([]string{string(game.TeamNameBlue), string(game.TeamNameRed)}, ", ")

type PatchGameRoomPlayerRequest struct {
	Team      *string `json:"team"`
	Spymaster *bool   `json:"spymaster"`
}

func (c *Controller) PatchGameRoomPlayer(gc *gin.Context) {

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

	c.EnsurePlayerIdMatchesSessionPlayer(gc, session)
	if c.HasErrors(gc) {
		return
	}

	//
	// Validate body parameters
	//

	request := PatchGameRoomPlayerRequest{}
	err := c.ParseBody(gc, &request)
	if err != nil || (request.Spymaster == nil && request.Team == nil) {
		c.APIError(gc, "A request body with either 'team' or 'spymaster' is required", err, 400)
		return
	}

	if request.Spymaster != nil && request.Team != nil {
		c.APIError(gc, "You can only set either 'team' or 'spymaster', not both at the same time!", nil, 400)
		return
	}

	//
	// Execute action :)
	//

	// Spymaster selection
	if request.Spymaster != nil {
		room, err := c.service.SetSpymaster(ctx, room, session.Player)
		if err != nil {
			c.APIError(gc, "Unexpected error while setting spymaster", err, 500)
			return
		}

		// Return the updated room
		c.APIResponse(gc, response.NewGameRoomResponse(room), 200)

	}

	// Team selection
	if request.Team != nil {

		if !slices.Contains(validTeamNames, game.TeamName(*request.Team)) {
			c.APIError(gc, fmt.Sprintf("Invalid team name provided '%s', please use one of: %s", *request.Team, validTeamNamesString), nil, 400)
		}

		room, err := c.service.AddPlayerToTeam(ctx, room, session.Player, game.TeamName(*request.Team))
		if err != nil {
			c.APIError(gc, "Unexpected error while adding player to team", err, 500)
			return
		}

		// Return the updated room
		c.APIResponse(gc, response.NewGameRoomResponse(room), 200)

	}

	// We should never hit this due to the validation above
	c.APIError(gc, "Unexpected state when updating a player in the game", nil, 500)
	return

}
