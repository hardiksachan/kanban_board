package handlers

import (
	"github.com/hardiksachan/kanban_board/backend/internal/users/core/domain"
	"time"
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

func (r *SignUpRequest) toDomain() *domain.User {
	return &domain.User{
		ID:         "",
		Name:       r.Name,
		Email:      r.Email,
		Password:   r.Password,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}
}
