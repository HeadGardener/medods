package storage

import (
	"context"
	"fmt"

	"github.com/HeadGardener/medods/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SessionStorage struct {
	coll *mongo.Collection
}

func NewSessionStorage(coll *mongo.Collection) *SessionStorage {
	return &SessionStorage{coll: coll}
}

func (s *SessionStorage) GetSessionByUserID(ctx context.Context, userID string) (models.Session, error) {
	var session models.Session

	if err := s.coll.FindOne(ctx, bson.D{{"user_id", userID}}).Decode(&session); err != nil {
		return models.Session{}, fmt.Errorf("failed while getting user by id session: %w", err)
	}

	return session, nil
}

func (s *SessionStorage) CreateUserSession(ctx context.Context, session models.Session) error {
	_, err := s.coll.InsertOne(ctx, session)
	if err != nil {
		return fmt.Errorf("failed while creating user sesion: %w", err)
	}

	return nil
}

func (s *SessionStorage) UpdateSession(ctx context.Context, session models.Session) error {
	filter := bson.D{{"user_id", session.UserID}}
	update := bson.D{{"$set", bson.D{
		{"id", session.ID},
		{"refresh_token", session.RefreshToken},
		{"expires_at", session.ExpiresAt},
	}}}

	if _, err := s.coll.UpdateOne(ctx, filter, update); err != nil {
		return fmt.Errorf("failed while updating user session: %w", err)
	}

	return nil
}

func (s *SessionStorage) GetSessionByID(ctx context.Context, id string) (models.Session, error) {
	var session models.Session

	if err := s.coll.FindOne(ctx, bson.D{{"id", id}}).Decode(&session); err != nil {
		return models.Session{}, fmt.Errorf("failed while getting session by id: %w", err)
	}

	return session, nil
}
