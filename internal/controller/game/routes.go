package game

import "github.com/gin-gonic/gin"

func (c *Controller) SetupRoutes(rg *gin.RouterGroup) {
	rg.POST("/rooms", c.PostGameRoom)
	rg.GET("/rooms/:name", c.GetGameRoom)
	rg.PATCH("/rooms/:name", c.WrapGameRoomRouteWithLock(c.PatchGameRoom))

	rg.POST("/rooms/:name/players", c.WrapGameRoomRouteWithLock(c.PostGameRoomPlayer))
	rg.PATCH("/rooms/:name/players/:playerId", c.WrapGameRoomRouteWithLock(c.PatchGameRoomPlayer))
	// rg.DELETE("/rooms/:id/players", c.RemoveGameRoomPlayer)

	rg.PATCH("/rooms/:name/tiles/:tileId", c.WrapGameRoomRouteWithLock(c.PatchGameRoomTile))
}
