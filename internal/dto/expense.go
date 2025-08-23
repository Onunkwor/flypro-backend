package dto

import (
	"time"

	"github.com/onunkwor/flypro-backend/internal/models"
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

func ToExpenseResponse(expense *models.Expense) ExpenseResponse {
	return ExpenseResponse{
		ID:          expense.ID,
		UserID:      expense.UserID,
		Amount:      expense.Amount,
		AmountUSD:   expense.AmountUSD,
		Currency:    expense.Currency,
		Category:    expense.Category,
		Description: expense.Description,
		Receipt:     expense.Receipt,
		Status:      expense.Status,
		CreatedAt:   expense.CreatedAt,
		UpdatedAt:   expense.UpdatedAt,
	}
}

func ToExpenseResponses(expenses []models.Expense) []ExpenseResponse {
	res := make([]ExpenseResponse, len(expenses))
	for i, exp := range expenses {
		res[i] = ToExpenseResponse(&exp)
	}
	return res
}
