package redis

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v9"
	"github.com/gofrs/uuid"
	"github.com/hardiksachan/kanban_board/backend/internal/users"
	"github.com/hardiksachan/kanban_board/backend/internal/users/core/domain"
	"time"
)

type SessionStore struct {
	client *redis.Client
}

func NewSessionStore(client *redis.Client) *SessionStore {
	return &SessionStore{client}
}

func (s *SessionStore) Create(userId string) (*domain.Session, error) {
	op := "redis.SessionStore.Create"

	sessionId, err := uuid.NewV4()
	if err != nil {
		return nil, &users.Error{Op: op, Err: err}
	}

	expiration := time.Minute * 10

	session := &domain.Session{
		ID:        sessionId.String(),
		UserID:    userId,
		ExpiresAt: time.Now().Add(expiration),
	}

	es, err := json.Marshal(*session)
	if err != nil {
		return nil, &users.Error{Op: op, Err: err}
	}

	err = s.client.Set(context.Background(), session.ID, string(es), expiration).Err()
	if err != nil {
		return nil, &users.Error{Op: op, Err: err}
	}

	return session, nil
}

func (s *SessionStore) Get(sessionId string) (*domain.Session, error) {
	op := "redis.SessionStore.Get"

	var session domain.Session
	es, err := s.client.Get(context.Background(), sessionId).Result()
	if err == redis.Nil {
		return nil, &users.Error{Op: op, Code: users.ENOTFOUND, Message: "invalid session ID", Err: err}
	}
	if err != nil {
		return nil, &users.Error{Op: op, Err: err}
	}

	err = json.Unmarshal([]byte(es), &session)
	if err != nil {
		return nil, &users.Error{Op: op, Err: err}
	}

	return &session, nil
}

func (s *SessionStore) Delete(sessionId string) error {
	op := "redis.SessionStore.Delete"

	err := s.client.Del(context.Background(), sessionId).Err()
	if err == redis.Nil {
		return &users.Error{Op: op, Code: users.ENOTFOUND, Message: "invalid session ID", Err: err}
	}
	if err != nil {
		return &users.Error{Op: op, Err: err}
	}

	return nil
}
