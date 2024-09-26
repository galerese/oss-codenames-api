package http_controller

// TODO - remove hard dependency on Gin?
import (
	"fmt"
	"net/http"

	"galere.se/oss-codenames-api/pkg/domain_util"
	"galere.se/oss-codenames-api/pkg/logging"
	"github.com/gin-gonic/gin"
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

func GenericErrorResponseHandler(l logging.Logger) gin.HandlerFunc {
	l.Info("Setting up generic error response handler")

	return func(gc *gin.Context) {
		gc.Next()

		handled := false

		for _, err := range gc.Errors {
			if apiError, ok := err.Err.(*ApiError); ok {
				gc.JSON(apiError.code, ErrorResponse{Error: apiError.message})
				handled = true
				break
			} else if _, ok := err.Err.(*domain_util.StateValidationError); ok {
				gc.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Err.Error()})
				handled = true
				break
			} else if _, ok := err.Err.(*domain_util.InvalidParameterError); ok {
				gc.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Err.Error()})
				handled = true
				break
			} else if _, ok := err.Err.(*domain_util.InvalidActionError); ok {
				gc.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Err.Error()})
				handled = true
				break
			}
		}

		if len(gc.Errors) > 0 && !handled {
			// Any other error, uncaught :o
			gc.JSON(http.StatusInternalServerError, ErrorResponse{Error: "there was an internal server error, please try again later - the developers have been notified"})
		}

	}

}
