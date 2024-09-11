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

func (s *Service) GetSession() *Session {
	return s.repository.GetSession()
}
