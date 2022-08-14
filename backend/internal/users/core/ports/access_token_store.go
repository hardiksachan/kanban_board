package ports

import "github.com/hardiksachan/kanban_board/backend/internal/users/core/domain"

type AccessTokenStore interface {
	Create(userId string) (*domain.AccessToken, error)
	VerifyAndDecode(encodedToken string) (*domain.AccessToken, error)
	Encode(token *domain.AccessToken) (string, error)
}
