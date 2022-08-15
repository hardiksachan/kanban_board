package native

import (
	"github.com/hardiksachan/kanban_board/backend/internal/users"
	"github.com/hardiksachan/kanban_board/backend/internal/users/core/domain"
)

type UserMetadataStore struct {
	users []*domain.User
}

func NewUserMetadataStore() *UserMetadataStore {
	return &UserMetadataStore{}
}

func (s *UserMetadataStore) Update(user *domain.User) error {
	for _, metadata := range s.users {
		if metadata.UserId == user.UserId {
			metadata.Name = user.Name
			metadata.ImageURL = user.ImageURL
			return nil
		}
	}

	s.users = append(s.users, user)
	return nil
}

func (s *UserMetadataStore) Get(userID string) (*domain.User, error) {
	for _, metadata := range s.users {
		if metadata.UserId == userID {
			return metadata, nil
		}
	}

	return nil, &users.Error{Op: "native.UserMetadataStore.Get", Code: users.ENOTFOUND}
}
