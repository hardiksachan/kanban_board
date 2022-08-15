package postgres

import (
	"context"
	"github.com/hardiksachan/kanban_board/backend/internal/users"
	"github.com/hardiksachan/kanban_board/backend/internal/users/core/domain"
	"github.com/hardiksachan/kanban_board/backend/internal/users/repository/postgres/user/dao"
	"github.com/hardiksachan/kanban_board/backend/shared"
	"github.com/jackc/pgx/v4"
)

type CredentialStore struct {
	q *dao.Queries
}

func NewCredentialStore(q *dao.Queries) *CredentialStore {
	return &CredentialStore{q}
}

func (s *CredentialStore) Insert(credential *domain.Credential) (*domain.Credential, error) {
	op := "postgres.CredentialStore.Insert"

	ctx := context.Background()
	credentialRow, err := s.q.InsertCredential(ctx, dao.InsertCredentialParams{
		Email:    credential.Email,
		Password: credential.Password,
	})
	if err != nil {
		return nil, &users.Error{Code: users.EINTERNAL, Op: op, Err: err}
	}

	return &domain.Credential{
		UserID:   credentialRow.UserID.String(),
		Email:    credentialRow.Email,
		Password: credentialRow.Password,
	}, nil
}

func (s *CredentialStore) Update(credential *domain.Credential) error {
	op := "postgres.CredentialStore.Update"

	ctx := context.Background()

	userUuid, err := shared.GetUUIDFromString(credential.UserID)
	if err != nil {
		return &users.Error{Code: users.EINVALID, Message: "Unable to parse UUID", Op: op, Err: err}
	}

	_, err = s.q.UpdatePassword(ctx, dao.UpdatePasswordParams{
		Password: credential.Password,
		UserID:   *userUuid,
	})
	if err != nil {
		return &users.Error{Code: users.EINTERNAL, Op: op, Err: err}
	}
	return nil
}

func (s *CredentialStore) Remove(credential *domain.Credential) error {
	op := "postgres.CredentialStore.Remove"

	ctx := context.Background()

	userUuid, err := shared.GetUUIDFromString(credential.UserID)
	if err != nil {
		return &users.Error{Code: users.EINVALID, Message: "Unable to parse UUID", Op: op, Err: err}
	}

	_, err = s.q.DeleteUser(ctx, *userUuid)
	if err != nil {
		return &users.Error{Code: users.EINTERNAL, Op: op, Err: err}
	}
	return nil
}

func (s *CredentialStore) FindById(userId string) (*domain.Credential, error) {
	op := "postgres.CredentialStore.FindById"

	userUuid, err := shared.GetUUIDFromString(userId)
	if err != nil {
		return nil, &users.Error{Code: users.EINVALID, Message: "Unable to parse UUID", Op: op, Err: err}
	}

	ctx := context.Background()
	dbUser, err := s.q.FindById(ctx, *userUuid)
	if err == pgx.ErrNoRows {
		return nil, &users.Error{Code: users.ENOTFOUND, Op: op, Err: err}
	}
	if err != nil {
		return nil, &users.Error{Code: users.EINTERNAL, Op: op, Err: err}
	}
	return &domain.Credential{
		UserID:   dbUser.UserID.String(),
		Email:    dbUser.Email,
		Password: dbUser.Password,
	}, nil
}

func (s *CredentialStore) FindByEmail(email string) (*domain.Credential, error) {
	op := "postgres.CredentialStore.FindByEmail"

	ctx := context.Background()
	dbUser, err := s.q.FindByEmail(ctx, email)
	if err == pgx.ErrNoRows {
		return nil, &users.Error{Code: users.ENOTFOUND, Op: op, Err: err}
	}
	if err != nil {
		return nil, &users.Error{Code: users.EINTERNAL, Op: op, Err: err}
	}

	return &domain.Credential{
		UserID:   dbUser.UserID.String(),
		Email:    dbUser.Email,
		Password: dbUser.Password,
	}, nil
}

func (s *CredentialStore) CountByEmail(email string) (int, error) {
	op := "postgres.CredentialStore.CountByEmail"

	ctx := context.Background()
	count, err := s.q.CountByEmail(ctx, email)
	if err == pgx.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return -1, &users.Error{Code: users.EINTERNAL, Op: op, Err: err}
	}

	return int(count), nil
}
