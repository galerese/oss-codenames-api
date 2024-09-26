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
	gameRoomLocker Locker
	*sc.Controller
}

func NewController(service *game.Service, sessionService *session.Service, gameRoomLocker Locker, logger *logging.Logger) *Controller {
	return &Controller{
		service:        service,
		sessionService: sessionService,
		gameRoomLocker: gameRoomLocker,
		Controller:     sc.NewController(sessionService, logger),
	}
}
