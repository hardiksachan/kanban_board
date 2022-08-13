package ports

import (
	"kanban_board/internal/users/core/domain"
)

type UserStore interface {
	Insert(user *domain.User) (*domain.User, error)
	Update(user *domain.User) error
	Remove(user *domain.User) error
	FindById(userId string) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	CheckByEmail(email string) (bool, error)
}
