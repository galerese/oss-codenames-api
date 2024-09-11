package session

type SessionRepository interface {
	GetSession() *Session
	SaveSession(session *Session)
}
