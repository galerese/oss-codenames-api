package database

import (
	"context"
	"fmt"

	"galere.se/oss-codenames-api/internal/domain/game"
	"galere.se/oss-codenames-api/internal/domain/session"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Database struct {
	db *mongo.Database
}

func NewDatabase(db *mongo.Database) *Database {
	return &Database{
		db: db,
	}
}

func (d *Database) GetSession(ctx context.Context, token string) (*session.Session, error) {
	session := session.Session{}
	err := d.db.Collection("sessions").FindOne(ctx, bson.M{"player.token": token}).Decode(&session)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get session")
	}
	return &session, nil
}

func (d *Database) SaveSession(ctx context.Context, session *session.Session) (*session.Session, error) {
	_, err := d.db.Collection("sessions").UpdateOne(ctx,
		bson.M{"player.token": session.Player.Token},
		bson.M{"$set": session},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to save session")
	}

	return session, nil
}

func (d *Database) GetGameRoomByCode(ctx context.Context, code string) (*game.GameRoom, error) {

	room := game.GameRoom{}
	err := d.db.Collection("game_rooms").FindOne(ctx, bson.M{"roomCode": code}).Decode(&room)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get game room by code")
	}

	return &room, nil
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
	_, err := d.db.Collection("game_rooms").UpdateOne(ctx,
		bson.M{"roomCode": room.RoomCode},
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
