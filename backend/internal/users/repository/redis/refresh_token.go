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

type Claims struct {
	ExpiresAt time.Time
	Token     string
	*domain.Credential
}

type RefreshStore struct {
	client *redis.Client
}

func NewRefreshTokenStore(client *redis.Client) *RefreshStore {
	return &RefreshStore{client}
}

func (s *RefreshStore) Create(credential *domain.Credential) (*domain.RefreshToken, error) {
	op := "redis.RefreshStore.Create"

	token, err := uuid.NewV4()
	if err != nil {
		return nil, &users.Error{Op: op, Err: err}
	}

	expiration := time.Minute * 10

	claims := Claims{
		ExpiresAt:  time.Now().Add(expiration),
		Credential: credential,
		Token:      token.String(),
	}

	encodedClaims, err := json.Marshal(claims)
	if err != nil {
		return nil, &users.Error{Op: op, Err: err}
	}

	err = s.client.Set(context.Background(), claims.Token, encodedClaims, expiration).Err()
	if err != nil {
		return nil, &users.Error{Op: op, Err: err}
	}

	return (*domain.RefreshToken)(&claims.Token), nil
}

func (s *RefreshStore) Revoke(token *domain.RefreshToken) error {
	op := "redis.RefreshStore.Revoke"

	err := s.client.Del(context.Background(), string(*token)).Err()
	if err == redis.Nil {
		return &users.Error{Op: op, Code: users.ENOTFOUND, Message: "invalid session TokenID", Err: err}
	}
	if err != nil {
		return &users.Error{Op: op, Err: err}
	}

	return nil
}

func (s *RefreshStore) Verify(token *domain.RefreshToken) (*domain.Credential, error) {
	op := "redis.RefreshStore.Verify"

	claimsJSON, err := s.client.Get(context.Background(), string(*token)).Result()
	if err == redis.Nil {
		return nil, &users.Error{Op: op, Code: users.ENOTFOUND, Message: "invalid refreshToken", Err: err}
	}
	if err != nil {
		return nil, &users.Error{Op: op, Err: err}
	}

	var claims Claims
	err = json.Unmarshal([]byte(claimsJSON), &claims)
	if err != nil {
		return nil, &users.Error{Op: op, Err: err}
	}

	if claims.ExpiresAt.Unix() < time.Now().Unix() {
		return nil, nil
	}

	return claims.Credential, nil
}