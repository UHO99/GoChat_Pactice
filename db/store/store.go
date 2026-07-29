package store

import (
	"context"
	db "gochat/db/sqlc"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	query *db.Queries
}

type Message struct {
	ID       int64
	RoomID   int64
	Nickname string
	Content  string
	CreateAt time.Time
}

func New(pool *pgxpool.Pool) *Store {
	return &Store{query: db.New(pool)}
}

func (s *Store) CreateUser(ctx context.Context, nickname string) (db.User, error) {
	return s.query.CreateUser(ctx, nickname)
}

func (s *Store) CreateRoom(ctx context.Context, name string) (db.Room, error) {
	return s.query.CreateRoom(ctx, name)
}

func (s *Store) InsertMessage(ctx context.Context, roomID, userID int64, content string) (db.Message, error) {
	return s.query.InsertMessage(ctx, db.InsertMessageParams{
		RoomID:  roomID,
		UserID:  userID,
		Content: content,
	})
}

func (s *Store) ListRooms(ctx context.Context) ([]db.Room, error) {
	return s.query.ListRooms(ctx)
}
