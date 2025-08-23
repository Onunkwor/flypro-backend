package dto

import (
	"time"

	"github.com/onunkwor/flypro-backend/internal/models"
	"github.com/onunkwor/flypro-backend/internal/utils"
)

type CreateReportRequest struct {
	Title  string `json:"title" binding:"required"`
	UserID uint   `json:"user_id" binding:"required"`
}

type AddExpenseToReportRequest struct {
	UserID    uint `json:"user_id" binding:"required"`
	ExpenseID uint `json:"expense_id" binding:"required"`
}

func (r *CreateReportRequest) Sanitize() {
	r.Title = utils.SanitizeString(r.Title)
}

type ExpenseDTO struct {
	ID          uint      `json:"id"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	AmountUSD   float64   `json:"amount_usd"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
}

type UserDTO struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ReportDTO struct {
	ID        uint         `json:"id"`
	Title     string       `json:"title"`
	User      UserDTO      `json:"user"`
	Expenses  []ExpenseDTO `json:"expenses"`
	Total     float64      `json:"total"`
	CreatedAt time.Time    `json:"created_at"`
}

func ToExpenseDTO(e models.Expense) ExpenseDTO {
	return ExpenseDTO{
		ID:          e.ID,
		Description: e.Description,
		Amount:      e.Amount,
		AmountUSD:   e.AmountUSD,
		Category:    e.Category,
		CreatedAt:   e.CreatedAt,
	}
}

func ToUserDTO(u models.User) UserDTO {
	return UserDTO{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}

func ToReportDTO(r models.ExpenseReport) ReportDTO {
	expenses := make([]ExpenseDTO, len(r.Expenses))
	return ReportDTO{
		ID:        r.ID,
		Title:     r.Title,
		User:      ToUserDTO(r.User),
		Expenses:  expenses,
		Total:     r.Total,
		CreatedAt: r.CreatedAt,
	}
}

func ToReportDTOs(reports []models.ExpenseReport) []ReportDTO {
	result := make([]ReportDTO, len(reports))
	for i, r := range reports {
		result[i] = ToReportDTO(r)
	}
	return result
}
