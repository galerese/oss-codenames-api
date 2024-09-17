package http_controller

// TODO - remove hard dependency on Gin?
import (
	"net/http"

	"galere.se/oss-codenames-api/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type ApiError struct {
	message string
	code    int
}

func (ae *ApiError) Error() string {
	return ae.message
}

func NewAPIError(message string, code int) *ApiError {
	return &ApiError{message: message, code: code}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func GenericErrorHandler(l logging.Logger) gin.HandlerFunc {
	return func(gc *gin.Context) {
		gc.Next()

		if gc.Errors != nil && len(gc.Errors) > 0 {
			err := gc.Errors[0]

			// Error logging
			l.Error(gc.Request.Context(), errors.Wrap(err, "api error").Error())

			if apiError, ok := err.Err.(*ApiError); ok {
				// API error sent by the application
				gc.JSON(apiError.code, ErrorResponse{Error: apiError.message})
			} else {
				// Any other error, uncaught :o
				gc.JSON(http.StatusInternalServerError, ErrorResponse{Error: "there was an internal server error, please try again later - the developers have been notified"})
			}

			gc.Abort()
			return
		}
	}

}
