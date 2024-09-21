package session

import "github.com/gin-gonic/gin"

func (c *Controller) PostSessions(gc *gin.Context) {

	session, err := c.service.CreateSession()
	if err != nil {
		c.APIError(gc, "Unexpected error on session creation", err, 500)
		return
	}

	c.APIResponse(gc, session, 200)
}
