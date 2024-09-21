package session

import (
	"galere.se/oss-codenames-api/internal/response"
	"github.com/gin-gonic/gin"
)

func (c *Controller) PostSessions(gc *gin.Context) {

	session, err := c.service.CreateSession(gc.Request.Context())
	if err != nil {
		c.APIError(gc, "Unexpected error on session creation", err, 500)
		return
	}

	c.APIResponse(gc, response.NewSessionResponse(session), 200)
}
