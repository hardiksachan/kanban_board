package native

import (
	"github.com/hardiksachan/kanban_board/backend/internal/users"
	"github.com/hardiksachan/kanban_board/backend/internal/users/core/domain"
)

type UserMetadataStore struct {
	users []*domain.UserMetadata
}

func NewUserMetadataStore() *UserMetadataStore {
	return &UserMetadataStore{}
}

func (s *UserMetadataStore) Update(user *domain.UserMetadata) error {
	for _, metadata := range s.users {
		if metadata.UserId == user.UserId {
			metadata.DisplayName = user.DisplayName
			metadata.ImageURL = user.ImageURL
			return nil
		}
	}

	s.users = append(s.users, user)
	return nil
}

func (s *UserMetadataStore) Get(userID string) (*domain.UserMetadata, error) {
	for _, metadata := range s.users {
		if metadata.UserId == userID {
			return metadata, nil
		}
	}

	return nil, &users.Error{Op: "native.UserMetadataStore.Get", Code: users.ENOTFOUND}
}
