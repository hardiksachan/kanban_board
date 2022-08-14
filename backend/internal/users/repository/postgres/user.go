package postgres

import (
	"context"
	"github.com/hardiksachan/kanban_board/backend/internal/users"
	"github.com/hardiksachan/kanban_board/backend/internal/users/core/domain"
	"github.com/hardiksachan/kanban_board/backend/internal/users/repository/postgres/user/dao"
	"github.com/hardiksachan/kanban_board/backend/shared"
	"github.com/jackc/pgx/v4"
)

type UserStore struct {
	q *dao.Queries
}

func NewUserStore(q *dao.Queries) *UserStore {
	return &UserStore{q}
}

func (s *UserStore) Insert(user *domain.User) (*domain.User, error) {
	op := "postgres.UserStore.Insert"

	ctx := context.Background()
	dbUser, err := s.q.InsertUser(ctx, dao.InsertUserParams{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		return nil, &users.Error{Code: users.EINTERNAL, Op: op, Err: err}
	}
	return &domain.User{
		ID:         dbUser.UserID.String(),
		Name:       dbUser.Name,
		Email:      dbUser.Email,
		Password:   dbUser.Password,
		CreatedAt:  dbUser.CreatedAt,
		ModifiedAt: dbUser.ModifiedAt,
	}, nil
}

func (s *UserStore) Update(user *domain.User) error {
	op := "postgres.UserStore.Update"

	ctx := context.Background()

	userUuid, err := shared.GetUUIDFromString(user.ID)
	if err != nil {
		return &users.Error{Code: users.EINVALID, Message: "Unable to parse UUID", Op: op, Err: err}
	}

	_, err = s.q.UpdateUser(ctx, dao.UpdateUserParams{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		UserID:   *userUuid,
	})
	if err != nil {
		return &users.Error{Code: users.EINTERNAL, Op: op, Err: err}
	}
	return nil
}

func (s *UserStore) Remove(user *domain.User) error {
	op := "postgres.UserStore.Remove"

	ctx := context.Background()

	userUuid, err := shared.GetUUIDFromString(user.ID)
	if err != nil {
		return &users.Error{Code: users.EINVALID, Message: "Unable to parse UUID", Op: op, Err: err}
	}

	_, err = s.q.DeleteUser(ctx, *userUuid)
	if err != nil {
		return &users.Error{Code: users.EINTERNAL, Op: op, Err: err}
	}
	return nil
}

func (s *UserStore) FindById(userId string) (*domain.User, error) {
	op := "postgres.UserStore.FindById"

	userUuid, err := shared.GetUUIDFromString(userId)
	if err != nil {
		return nil, &users.Error{Code: users.EINVALID, Message: "Unable to parse UUID", Op: op, Err: err}
	}

	ctx := context.Background()
	dbUser, err := s.q.GetUserById(ctx, *userUuid)
	if err == pgx.ErrNoRows {
		return nil, &users.Error{Code: users.ENOTFOUND, Op: op, Err: err}
	}
	if err != nil {
		return nil, &users.Error{Code: users.EINTERNAL, Op: op, Err: err}
	}
	return &domain.User{
		ID:         dbUser.UserID.String(),
		Name:       dbUser.Name,
		Email:      dbUser.Email,
		Password:   dbUser.Password,
		CreatedAt:  dbUser.CreatedAt,
		ModifiedAt: dbUser.ModifiedAt,
	}, nil
}

func (s *UserStore) FindByEmail(email string) (*domain.User, error) {
	op := "postgres.UserStore.FindByEmail"

	ctx := context.Background()
	dbUser, err := s.q.GetUserByEmail(ctx, email)
	if err == pgx.ErrNoRows {
		return nil, &users.Error{Code: users.ENOTFOUND, Op: op, Err: err}
	}
	if err != nil {
		return nil, &users.Error{Code: users.EINTERNAL, Op: op, Err: err}
	}
	return &domain.User{
		ID:         dbUser.UserID.String(),
		Name:       dbUser.Name,
		Email:      dbUser.Email,
		Password:   dbUser.Password,
		CreatedAt:  dbUser.CreatedAt,
		ModifiedAt: dbUser.ModifiedAt,
	}, nil
}

func (s *UserStore) CheckByEmail(email string) (bool, error) {
	op := "postgres.UserStore.CheckByEmail"

	ctx := context.Background()
	_, err := s.q.GetUserByEmail(ctx, email)
	if err == pgx.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, &users.Error{Code: users.EINTERNAL, Op: op, Err: err}
	}
	return true, nil
}
