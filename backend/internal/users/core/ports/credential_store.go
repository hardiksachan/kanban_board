package ports

import (
	"github.com/hardiksachan/kanban_board/backend/internal/users/core/domain"
)

type CredentialStore interface {
	Insert(*domain.Credential) (*domain.Credential, error)
	Update(*domain.Credential) error
	Remove(*domain.Credential) error
	FindByEmail(email string) (*domain.Credential, error)
	FindById(UserId string) (*domain.Credential, error)
	CountByEmail(email string) (int, error)
}
