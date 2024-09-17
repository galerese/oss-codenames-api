package session

type SessionRepository interface {
	GetSession(token string) *Session
	SaveSession(session *Session) (*Session, error)
}
