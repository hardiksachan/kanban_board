package ports

import "kanban_board/internal/users/core/domain"

type SessionStore interface {
	Create(userId string) (*domain.Session, error)
	Get(sessionId string) (*domain.Session, error)
	Delete(sessionId string) error
}
