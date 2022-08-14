package ports

import (
	"github.com/hardiksachan/kanban_board/backend/internal/users/core/domain"
)

type RefreshTokenStore interface {
	Create(userId string) (*domain.RefreshToken, error)
	Get(tokenId string) (*domain.RefreshToken, error)
	Delete(encodedRefreshToken string) error
	Encode(token *domain.RefreshToken) (string, error)
	VerifyAndDecode(encodedRefreshToken string) (*domain.RefreshToken, error)
}
