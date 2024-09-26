package game

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Wraps a handler with a game room lock
func (c *Controller) WrapGameRoomRouteWithLock(f gin.HandlerFunc) gin.HandlerFunc {

	return func(gc *gin.Context) {
		roomName := gc.Param("name")
		lockName := fmt.Sprintf("game-room-%s-lock", roomName)

		err := c.gameRoomLocker.AcquireLock(gc.Request.Context(), lockName)
		if err != nil {
			c.APIError(gc, "Failed to acquire game room lock for room ["+roomName+"]", err, 500)
			return
		}

		f(gc)

		defer func() {
			err := c.gameRoomLocker.ReleaseLock(gc.Request.Context(), lockName)
			if err != nil {
				c.APIError(gc, "Failed to release game room lock for room ["+roomName+"]", err, 500)
			}
		}()
	}

}
