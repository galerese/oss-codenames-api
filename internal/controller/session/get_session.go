package session

import (
	"galere.se/oss-codenames-api/internal/response"
	"github.com/gin-gonic/gin"
)

func (c *Controller) GetSession(gc *gin.Context) {

	session, err := c.GetSessionFromRequest(gc)
	if err != nil {
		gc.Error(err)
		return
	}

	c.APIResponse(gc, response.NewSessionResponse(session), 200)

}
