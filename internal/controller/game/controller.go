package game

import (
	sc "galere.se/oss-codenames-api/internal/controller/session"
	"galere.se/oss-codenames-api/internal/domain/game"
	"galere.se/oss-codenames-api/internal/domain/session"
	"galere.se/oss-codenames-api/pkg/logging"
)

type Controller struct {
	service        *game.Service
	sessionService *session.Service
	*sc.Controller
}

func NewController(service *game.Service, sessionService *session.Service, logger logging.Logger) *Controller {
	return &Controller{
		service:        service,
		sessionService: sessionService,
		Controller:     sc.NewController(sessionService, logger),
	}
}
