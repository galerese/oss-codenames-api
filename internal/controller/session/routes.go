package session

import "github.com/gin-gonic/gin"

func (c *Controller) SetupRoutes(rg *gin.RouterGroup) {
	rg.GET("/session", c.GetSession)
	rg.POST("/sessions", c.PostSessions)
}
