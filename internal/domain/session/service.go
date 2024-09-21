package session

import (
	"context"

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
func (s *Service) CreateSession(ctx context.Context) (*Session, error) {
	session := &Session{
		Player: &game.Player{
			Id:    uuid.New().String(),
			Token: uuid.New().String(),
		},
	}

	session, err := s.SaveSession(ctx, session)
	if err != nil {
		return nil, domain_util.NewUnexpectedError(err, "failed to save session on session creation")
	}

	return session, nil
}

// Gets a session by token
func (s *Service) GetSession(ctx context.Context, token string) (*Session, error) {
	return s.repository.GetSession(ctx, token)
}

// Saves a session
func (s *Service) SaveSession(ctx context.Context, session *Session) (*Session, error) {
	return s.repository.SaveSession(ctx, session)
}
