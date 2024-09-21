package game

import "context"

type GameRepository interface {
	GetGameRoomByCode(ctx context.Context, code string) (*GameRoom, error)
	GetGameRoomByName(ctx context.Context, name string) (*GameRoom, error)
	SaveGameRoom(ctx context.Context, gameRoom *GameRoom) error
	GetRandomBoardTiles(ctx context.Context, count int) (map[int]BoardTile, error)
}
