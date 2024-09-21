package session

import "galere.se/oss-codenames-api/internal/domain/game"

type Session struct {
	Player      *game.Player
	CurrentRoom *game.GameRoom
}
