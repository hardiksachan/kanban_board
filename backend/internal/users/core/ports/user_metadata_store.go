package ports

import "github.com/hardiksachan/kanban_board/backend/internal/users/core/domain"

type UserMetadataStore interface {
	Update(user *domain.UserMetadata) error
	Get(userID string) (*domain.UserMetadata, error)
}
