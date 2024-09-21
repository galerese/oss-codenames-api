package session

import (
	"strings"

	"galere.se/oss-codenames-api/internal/domain/session"
	"galere.se/oss-codenames-api/pkg/http_controller"
	"github.com/gin-gonic/gin"
)

// Generic function to get a session from the request authorization header
func (c *Controller) GetSessionFromRequest(gc *gin.Context) (*session.Session, error) {

	// Validate and parse token from headers

	authorization := gc.GetHeader("Authorization")
	if authorization == "" {
		return nil, http_controller.NewAPIError("Authorization header is required to get a session", nil, 401)
	}

	parts := strings.Split(authorization, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return nil, http_controller.NewAPIError("Authorization header must be in the format 'Bearer <token>'", nil, 400)
	}

	token := parts[1]

	// Get the session from the service

	session, err := c.service.GetSession(token)
	if err != nil {
		return nil, http_controller.NewAPIError("Unexpected error when getting the session", err, 500)
	}

	return session, nil
}

// Makes sure a session exists for the current request or fails
func (c *Controller) EnsureSessionForRequest(gc *gin.Context) *session.Session {

	session, err := c.GetSessionFromRequest(gc)
	if err != nil {
		c.Error(gc, "You must have an authenticated session to perform this action, and we could not find one", err)
		return nil
	}

	if session == nil {
		c.APIError(gc, "You must have an authenticated session to perform this action, and a session was not found for the specified token", nil, 401)
		return nil
	}

	return session
}
