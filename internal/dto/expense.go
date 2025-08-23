package dto

import "github.com/onunkwor/flypro-backend/internal/utils"

type CreateExpenseRequest struct {
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Currency    string  `json:"currency" binding:"required,len=3,currency"`
	Category    string  `json:"category" binding:"required,oneof=travel meals office supplies"`
	Description string  `json:"description" binding:"max=500"`
	UserID      uint    `json:"user_id" binding:"required"`
}

func (r *CreateExpenseRequest) Sanitize() {
	r.Currency = utils.NormalizeCurrency(r.Currency)
	r.Category = utils.NormalizeCategory(r.Category)
	r.Description = utils.SanitizeString(r.Description)
}

type UpdateExpenseRequest struct {
	UserID      uint    `json:"user_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Currency    string  `json:"currency" binding:"required,len=3,currency"`
	Category    string  `json:"category" binding:"required,oneof=travel meals office supplies"`
	Description string  `json:"description"`
}

func (r *UpdateExpenseRequest) Sanitize() {
	r.Currency = utils.NormalizeCurrency(r.Currency)
	r.Category = utils.NormalizeCategory(r.Category)
	r.Description = utils.SanitizeString(r.Description)
}
