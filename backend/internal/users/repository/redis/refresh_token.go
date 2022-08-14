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

type RefreshTokenStore struct {
	client *redis.Client
}

func NewRefreshTokenStore(client *redis.Client) *RefreshTokenStore {
	return &RefreshTokenStore{client}
}

func (s *RefreshTokenStore) Encode(token *domain.RefreshToken) (string, error) {
	return token.TokenID, nil
}

func (s *RefreshTokenStore) Create(userId string) (*domain.RefreshToken, error) {
	op := "redis.RefreshTokenStore.Create"

	refreshTokenId, err := uuid.NewV4()
	if err != nil {
		return nil, &users.Error{Op: op, Err: err}
	}

	expiration := time.Minute * 10

	session := &domain.RefreshToken{
		TokenID:   refreshTokenId.String(),
		UserID:    userId,
		ExpiresAt: time.Now().Add(expiration),
	}

	es, err := json.Marshal(*session)
	if err != nil {
		return nil, &users.Error{Op: op, Err: err}
	}

	err = s.client.Set(context.Background(), session.TokenID, string(es), expiration).Err()
	if err != nil {
		return nil, &users.Error{Op: op, Err: err}
	}

	return session, nil
}

func (s *RefreshTokenStore) VerifyAndDecode(encodedToken string) (*domain.RefreshToken, error) {
	op := "redis.RefreshTokenStore.VerifyAndDecode"

	storedToken, err := s.Get(encodedToken)
	if err != nil {
		return nil, &users.Error{Op: op, Err: err}
	}

	if storedToken.ExpiresAt.Unix() < time.Now().Unix() {
		return nil, nil
	}

	return storedToken, nil
}

func (s *RefreshTokenStore) Get(sessionId string) (*domain.RefreshToken, error) {
	op := "redis.RefreshTokenStore.Get"

	var refreshToken domain.RefreshToken
	es, err := s.client.Get(context.Background(), sessionId).Result()
	if err == redis.Nil {
		return nil, &users.Error{Op: op, Code: users.ENOTFOUND, Message: "invalid refreshToken", Err: err}
	}
	if err != nil {
		return nil, &users.Error{Op: op, Err: err}
	}

	err = json.Unmarshal([]byte(es), &refreshToken)
	if err != nil {
		return nil, &users.Error{Op: op, Err: err}
	}

	return &refreshToken, nil
}

func (s *RefreshTokenStore) Delete(tokenId string) error {
	op := "redis.RefreshTokenStore.Delete"

	err := s.client.Del(context.Background(), tokenId).Err()
	if err == redis.Nil {
		return &users.Error{Op: op, Code: users.ENOTFOUND, Message: "invalid session TokenID", Err: err}
	}
	if err != nil {
		return &users.Error{Op: op, Err: err}
	}

	return nil
}
