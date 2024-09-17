package http_controller

import (
	"galere.se/oss-codenames-api/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// This acts as a base controller for all controllers, with common functions
type BaseController struct {
	Logger logging.Logger
}

func NewBaseController(l logging.Logger) *BaseController {
	return &BaseController{Logger: l}
}

// This function is used to emit API errors with specific status codes
func (c *BaseController) APIError(gc *gin.Context, msgError string, originalError error, status int) {
	gc.Error(NewAPIError(msgError, status))
}

func (c *BaseController) APIResponse(gc *gin.Context, response interface{}, status int) {
	gc.JSON(status, response)
}

// This function is used to emit unexpected errors from the controller, which will likely become a 500 error
func (c *BaseController) Error(gc *gin.Context, msgError string, originalError error) {
	gc.Error(errors.Wrap(originalError, msgError))
}

func (c *BaseController) HasErrors(gc *gin.Context) bool {
	return len(gc.Errors) > 0
}

func (c *BaseController) ParseBody(gc *gin.Context, request interface{}) error {
	if err := gc.BindJSON(request); err != nil {
		return errors.Wrap(err, "could not parse request body")
	}
	return nil
}
