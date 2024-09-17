package game

import "github.com/gin-gonic/gin"

func (c *Controller) SetupRoutes(rg gin.RouterGroup) {
	rg.POST("/rooms", c.PostGameRoom)
	rg.GET("/rooms/:name", c.GetGameRoom)
	rg.PATCH("/rooms/:name", c.PatchGameRoom)

	rg.POST("/rooms/:name/players", c.PostGameRoomPlayer)
	rg.PATCH("/rooms/:id/players/:playerId", c.PatchGameRoomPlayer)
	// rg.DELETE("/rooms/:id/players", c.RemoveGameRoomPlayer)

	rg.PATCH("/rooms/:name/tiles/:tileId", c.PatchGameRoomTile)
}
