package dto

import "github.com/onunkwor/flypro-backend/internal/utils"

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
