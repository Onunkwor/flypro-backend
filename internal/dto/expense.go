package dto

import (
	"time"

	"github.com/onunkwor/flypro-backend/internal/utils"
)

type CreateExpenseRequest struct {
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Currency    string  `json:"currency" binding:"required,len=3,currency"`
	Category    string  `json:"category" binding:"required,category"`
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
	Category    string  `json:"category" binding:"required,category"`
	Description string  `json:"description"`
}

func (r *UpdateExpenseRequest) Sanitize() {
	r.Currency = utils.NormalizeCurrency(r.Currency)
	r.Category = utils.NormalizeCategory(r.Category)
	r.Description = utils.SanitizeString(r.Description)
}

type ExpenseResponse struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	Amount      float64   `json:"amount"`
	AmountUSD   float64   `json:"amount_usd"`
	Currency    string    `json:"currency"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Receipt     string    `json:"receipt"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
