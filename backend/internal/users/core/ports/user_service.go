package ports

import (
	"github.com/hardiksachan/kanban_board/backend/internal/users"
	"github.com/hardiksachan/kanban_board/backend/internal/users/core/domain"
)

type UserService struct {
	store UserMetadataStore
}

func NewUserService(store UserMetadataStore) *UserService {
	return &UserService{store}
}

func (s *UserService) Update(user *domain.UserMetadata) error {
	op := "ports.UserService.Update"

	err := s.store.Update(user)
	if err != nil {
		return &users.Error{Op: op, Err: err}
	}

	return nil
}

func (s *UserService) Find(userId string) (*domain.UserMetadata, error) {
	op := "ports.UserService.Update"

	data, err := s.store.Get(userId)
	if err != nil {
		return nil, &users.Error{Op: op, Err: err}
	}

	return data, nil
}
