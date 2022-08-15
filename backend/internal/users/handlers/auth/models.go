package auth

import (
	"github.com/hardiksachan/kanban_board/backend/internal/users/core/domain"
)

type SignUpRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LogInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LogInResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	UserId       string `json:"user_id,omitempty"`
}

type LogOutRequest struct {
	RefreshToken string `json:"refresh_token,omitempty"`
}

type RefreshAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func (r *SignUpRequest) toDomain() *domain.Credential {
	return &domain.Credential{
		UserID:   "",
		Email:    r.Email,
		Password: r.Password,
	}
}
