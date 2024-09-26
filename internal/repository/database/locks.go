package database

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"
)

const LOCK_TIME_TO_LIVE = 30 * time.Second
const LOCK_TIMEOUT = 4 * time.Second
const LOCK_SLEEP_INTERVAL = 50 * time.Millisecond

func (d *Database) SetupLockExpiration(ctx context.Context) error {
	d.l.Debugf("Setting up lock expiration index")

	_, err := d.db.Collection("locks").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"createdAt": 1},
		Options: options.Index().SetExpireAfterSeconds(int32(LOCK_TIME_TO_LIVE.Seconds())),
	})
	if err != nil {
		return errors.Wrap(err, "failed to create lock expiration index")
	}
	return nil
}

func (d *Database) AcquireLock(ctx context.Context, lockName string) error {
	d.l.Debugf("Acquiring lock", zap.String("lockName", lockName))

	var err error
	timeoutAt := time.Now().Add(LOCK_TIMEOUT)
	for timeoutAt.After(time.Now()) {
		_, err = d.db.Collection("locks").InsertOne(ctx, bson.M{
			"_id":       lockName,
			"createdAt": time.Now(),
		})
		if err == nil {
			return nil
		}
		time.Sleep(LOCK_SLEEP_INTERVAL)
	}

	return errors.Wrap(err, fmt.Sprintf("failed to acquire lock [%s]", lockName))
}

func (d *Database) ReleaseLock(ctx context.Context, lockName string) (err error) {
	d.l.Debugf("Releasing lock", zap.String("lockName", lockName))

	_, err = d.db.Collection("locks").DeleteOne(ctx, bson.M{"_id": lockName})
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to release lock [%s]", lockName))
	}

	return nil
}
