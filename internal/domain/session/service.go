package session

type Service struct {
	repository SessionRepository
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) CreateSession() *Session {
	return &Session{}
}

func (s *Service) GetSession(token string) (*Session, error) {
	return s.repository.GetSession(token), nil
}

// Saves a session
func (s *Service) SaveSession(session *Session) (*Session, error) {
	return s.repository.SaveSession(session)
}
