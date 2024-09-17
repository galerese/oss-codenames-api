package session

import "github.com/gin-gonic/gin"

func (c *Controller) GetSession(gc *gin.Context) {

	session, err := c.GetSessionFromRequest(gc)
	if err != nil {
		gc.Error(err)
		return
	}

	c.APIResponse(gc, session, 200)

}
