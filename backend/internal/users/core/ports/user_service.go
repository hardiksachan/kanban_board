package ports

import (
	"github.com/hardiksachan/kanban_board/backend/internal/users"
	"github.com/hardiksachan/kanban_board/backend/internal/users/core/domain"
)

type UserService struct {
	store UserStore
}

func NewUserService(store UserStore) *UserService {
	return &UserService{store}
}

func (s *UserService) Update(user *domain.User) error {
	op := "ports.UserService.Update"

	_, err := s.store.Update(user)
	if err != nil {
		return &users.Error{Op: op, Err: err}
	}

	return nil
}

func (s *UserService) Find(userId string) (*domain.User, error) {
	op := "ports.UserService.Update"

	data, err := s.store.Get(userId)
	if err != nil {
		return nil, &users.Error{Op: op, Err: err}
	}

	return data, nil
}
