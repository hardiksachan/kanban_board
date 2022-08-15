package ports

import "github.com/hardiksachan/kanban_board/backend/internal/users/core/domain"

type RefreshStore interface {
	Create(*domain.Credential) (*domain.RefreshToken, error)
	Revoke(*domain.RefreshToken) error
	Verify(*domain.RefreshToken) (*domain.Credential, error)
}
