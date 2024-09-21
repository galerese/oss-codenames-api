package session

import (
	"galere.se/oss-codenames-api/internal/domain/session"
	"galere.se/oss-codenames-api/pkg/http_controller"
	"galere.se/oss-codenames-api/pkg/logging"
)

type Controller struct {
	service *session.Service
	*http_controller.BaseController
}

func NewController(service *session.Service, logger *logging.Logger) *Controller {
	return &Controller{
		service:        service,
		BaseController: http_controller.NewBaseController(logger),
	}
}
