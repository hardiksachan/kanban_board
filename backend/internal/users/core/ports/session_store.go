package ports

import (
	"github.com/hardiksachan/kanban_board/backend/internal/users/core/domain"
)

type SessionStore interface {
	Create(userId string) (*domain.Session, error)
	Get(sessionId string) (*domain.Session, error)
	Delete(sessionId string) error
}
