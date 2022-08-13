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
