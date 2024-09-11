package game

type GameRepository interface {
	GetGameRoomByCode(code string) (*GameRoom, error)
	GetGameRoomByName(name string) (*GameRoom, error)
	SaveGameRoom(gameRoom *GameRoom) error
}
