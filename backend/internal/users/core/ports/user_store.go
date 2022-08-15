package ports

import "github.com/hardiksachan/kanban_board/backend/internal/users/core/domain"

type UserStore interface {
	Update(*domain.User) (*domain.User, error)
	Get(userID string) (*domain.User, error)
}
