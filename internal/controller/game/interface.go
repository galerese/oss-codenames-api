package game

import "context"

type Locker interface {
	AcquireLock(ctx context.Context, lockName string) error
	ReleaseLock(ctx context.Context, lockName string) error
}
