package database

import (
	"context"
	"fmt"

	"galere.se/oss-codenames-api/internal/domain/game"
	"galere.se/oss-codenames-api/internal/domain/session"
	"galere.se/oss-codenames-api/pkg/logging"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"
)

type Database struct {
	db *mongo.Database
	l  logging.Logger
}

func NewDatabase(db *mongo.Database, l logging.Logger) *Database {
	return &Database{
		db: db,
		l:  l,
	}
}

func (d *Database) GetSession(ctx context.Context, token string) (*session.Session, error) {

	session := session.Session{}
	sessionData := SessionData{}

	// Load session by token
	err := d.db.Collection("sessions").FindOne(ctx, bson.M{"player.token": token}).Decode(&sessionData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to get session")
	}
	session.Player = sessionData.Player

	// Load room from separate collection :)
	if sessionData.CurrentRoomId != nil {
		room := game.GameRoom{}
		err = d.db.Collection("game_rooms").FindOne(ctx, bson.M{"id": sessionData.CurrentRoomId}).Decode(&room)
		if err != nil && err != mongo.ErrNoDocuments {
			return nil, errors.Wrap(err, "failed to get game room")
		}
		session.CurrentRoom = &room
	}

	return &session, nil
}

func (d *Database) SaveSession(ctx context.Context, session *session.Session) (*session.Session, error) {

	sessionData := SessionData{
		Player: session.Player,
	}

	if session.CurrentRoom != nil {
		sessionData.CurrentRoomId = &session.CurrentRoom.Id
	}

	_, err := d.db.Collection("sessions").UpdateOne(ctx,
		bson.M{"player.token": session.Player.Token},
		bson.M{"$set": sessionData},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to save session")
	}

	return session, nil
}

func (d *Database) GetGameRoomByName(ctx context.Context, name string) (*game.GameRoom, error) {
	room := game.GameRoom{}
	err := d.db.Collection("game_rooms").FindOne(ctx, bson.M{"name": name}).Decode(&room)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get game room by name")
	}

	return &room, nil
}

func (d *Database) SaveGameRoom(ctx context.Context, room *game.GameRoom) error {
	d.l.Info("Saving game room", zap.Any("room", fmt.Sprintf("%+v", room)))

	_, err := d.db.Collection("game_rooms").UpdateOne(ctx,
		bson.M{"id": room.Id},
		bson.M{"$set": room},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return errors.Wrap(err, "failed to save game room")
	}

	return nil
}

func (d *Database) GetRandomBoardTiles(ctx context.Context, count int) (map[int]game.BoardTile, error) {
	return nil, fmt.Errorf("not implemented")
}
