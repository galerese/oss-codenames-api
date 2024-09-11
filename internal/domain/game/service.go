package game

import "github.com/pkg/errors"

type Service struct {
	repository GameRepository
}

func NewService(repository GameRepository) *Service {
	return &Service{
		repository: repository,
	}
}

// Gets a get game room by its code
func (s *Service) GetGameRoomByCode(code string) (*GameRoom, error) {
	return s.repository.GetGameRoomByCode(code)
}

// Gets a get game room by its name
func (s *Service) GetGameRoomByName(name string) (*GameRoom, error) {
	return s.repository.GetGameRoomByName(name)
}

// Creates a game room with a random name
func (s *Service) CreateGameRoom(player Player) (*GameRoom, error) {

	gameRoom := &GameRoom{
		State:   GameRoomStateWaitingForPlayers,
		Players: []Player{player},
	}

	err := s.repository.SaveGameRoom(gameRoom)
	if err != nil {
		return nil, errors.Wrap(err, "failed to save new game room")
	}

	return gameRoom, nil
}
