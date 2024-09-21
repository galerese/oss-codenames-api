package database

import (
	"fmt"

	"galere.se/oss-codenames-api/internal/domain/game"
	"galere.se/oss-codenames-api/internal/domain/session"
)

type Database struct {
}

func NewDatabase() *Database {
	return &Database{}
}

func (d *Database) GetSession(token string) *session.Session {
	return nil
}

func (d *Database) SaveSession(session *session.Session) (*session.Session, error) {
	return nil, fmt.Errorf("not implemented")
}

func (d *Database) GetGameRoomByCode(code string) (*game.GameRoom, error) {
	return nil, fmt.Errorf("not implemented")
}

func (d *Database) GetGameRoomByName(name string) (*game.GameRoom, error) {
	return nil, fmt.Errorf("not implemented")
}

func (d *Database) SaveGameRoom(gameRoom *game.GameRoom) error {
	return fmt.Errorf("not implemented")
}

func (d *Database) GetRandomBoardTiles(count int) (map[int]game.BoardTile, error) {
	return nil, fmt.Errorf("not implemented")
}
