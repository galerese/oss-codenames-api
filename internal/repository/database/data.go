package database

import "galere.se/oss-codenames-api/internal/domain/game"

type SessionData struct {
	Player        *game.Player
	CurrentRoomId *string
}
