package session

import (
	"galere.se/oss-codenames-api/internal/domain/game"
	"galere.se/oss-codenames-api/pkg/domain_util"
	"galere.se/oss-codenames-api/pkg/logging"
	"github.com/google/uuid"
)

type Service struct {
	repository  SessionRepository
	baseService *domain_util.BaseService
}

func NewService(repository SessionRepository, logger *logging.Logger) *Service {
	return &Service{
		repository:  repository,
		baseService: domain_util.NewBaseService(logger),
	}
}

// Creates a new player session
func (s *Service) CreateSession() (*Session, error) {
	session := &Session{
		Player: &game.Player{
			Id: uuid.New().String(),
		},
	}

	session, err := s.SaveSession(session)
	if err != nil {
		return nil, domain_util.NewUnexpectedError(err, "failed to save session on session creation")
	}

	return session, nil
}

// Gets a session by token
func (s *Service) GetSession(token string) (*Session, error) {
	return s.repository.GetSession(token), nil
}

// Saves a session
func (s *Service) SaveSession(session *Session) (*Session, error) {
	return s.repository.SaveSession(session)
}
