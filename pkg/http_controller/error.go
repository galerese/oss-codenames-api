package http_controller

// TODO - remove hard dependency on Gin?
import (
	"fmt"
	"net/http"

	"galere.se/oss-codenames-api/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type ApiError struct {
	message       string
	originalError error
	code          int
}

func (ae *ApiError) Error() string {
	if ae.originalError != nil {
		return fmt.Sprintf("APIError: %s: %s", ae.message, ae.originalError.Error())
	}
	return fmt.Sprintf("APIError: %s", ae.message)
}

func NewAPIError(message string, originalError error, code int) *ApiError {
	return &ApiError{
		message:       message,
		originalError: originalError,
		code:          code,
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func GenericErrorHandler(l logging.Logger) gin.HandlerFunc {
	l.Info("Setting up generic error handler")

	return func(gc *gin.Context) {
		gc.Next()

		if gc.Errors != nil && len(gc.Errors) > 0 {
			err := gc.Errors[0]

			if apiError, ok := err.Err.(*ApiError); ok {
				// API error sent by the application
				gc.JSON(apiError.code, ErrorResponse{Error: apiError.message})
				// l.Error(apiError)
			} else {
				// Any other error, uncaught :o
				gc.JSON(http.StatusInternalServerError, ErrorResponse{Error: "there was an internal server error, please try again later - the developers have been notified"})
				l.Error(errors.Wrap(err, "Unexpected API error").Error())
			}

			gc.Abort()
			return
		}
	}

}
