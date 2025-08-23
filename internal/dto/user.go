package dto

import (
	"github.com/onunkwor/flypro-backend/internal/models"
	"github.com/onunkwor/flypro-backend/internal/utils"
)

type CreateUserRequest struct {
	Email string `json:"email" binding:"required,email"`
	Name  string `json:"name" binding:"required"`
}

func (r *CreateUserRequest) Sanitize() {
	r.Email = utils.SanitizeString(r.Email)
	r.Name = utils.SanitizeString(r.Name)
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func NewUserResponse(user *models.User) UserResponse {
	return UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}
}
