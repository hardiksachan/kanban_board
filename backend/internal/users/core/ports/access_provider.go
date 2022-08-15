package ports

import "github.com/hardiksachan/kanban_board/backend/internal/users/core/domain"

type AccessProvider interface {
	Create(*domain.AccessClaims) (*domain.AccessToken, error)
	Verify(*domain.AccessToken) (*domain.AccessClaims, error)
}
