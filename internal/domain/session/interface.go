package session

import "context"

type SessionRepository interface {
	GetSession(ctx context.Context, token string) (*Session, error)
	SaveSession(ctx context.Context, session *Session) (*Session, error)
}
