package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/hardiksachan/kanban_board/backend/internal/users"
	"github.com/hardiksachan/kanban_board/backend/internal/users/core/domain"
	"github.com/hardiksachan/kanban_board/backend/internal/users/repository/postgres/user/dao"
	"github.com/hardiksachan/kanban_board/backend/shared"
	"github.com/jackc/pgx/v4"
)

type UserMetadataStore struct {
	q *dao.Queries
}

func NewUserMetadataStore(q *dao.Queries) *UserMetadataStore {
	return &UserMetadataStore{q}
}

func (s *UserMetadataStore) Update(user *domain.UserMetadata) error {
	op := "postgres.UserMetadataStore.Update"

	ctx := context.Background()

	userUuid, err := shared.GetUUIDFromString(user.UserId)
	if err != nil {
		return &users.Error{Code: users.EINVALID, Message: "Unable to parse UUID", Op: op, Err: err}
	}

	_, err = s.q.UpdateUserMetadata(ctx, dao.UpdateUserMetadataParams{
		DisplayName: user.DisplayName,
		ProfileImageUrl: sql.NullString{
			String: user.ImageURL,
			Valid:  true,
		},
		UserID: *userUuid,
	})
	if err != nil {
		return &users.Error{Code: users.EINTERNAL, Op: op, Err: err}
	}
	return nil
}

func (s *UserMetadataStore) Get(userID string) (*domain.UserMetadata, error) {
	op := "postgres.UserMetadataStore.Get"

	userUuid, err := shared.GetUUIDFromString(userID)
	if err != nil {
		return nil, &users.Error{Code: users.EINVALID, Message: fmt.Sprintf("Unable to parse UUID. id: %v", userID), Op: op, Err: err}
	}

	ctx := context.Background()
	dbUser, err := s.q.GetUserMetadata(ctx, *userUuid)
	if err == pgx.ErrNoRows {
		return nil, &users.Error{Code: users.ENOTFOUND, Op: op, Err: err}
	}
	if err != nil {
		return nil, &users.Error{Code: users.EINTERNAL, Op: op, Err: err}
	}

	return &domain.UserMetadata{
		UserId:      dbUser.UserID.String(),
		DisplayName: dbUser.DisplayName,
		ImageURL:    dbUser.ProfileImageUrl.String,
	}, nil
}
