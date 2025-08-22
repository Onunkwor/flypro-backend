package dto

type CreateReportRequest struct {
	Title  string `json:"title" binding:"required"`
	UserID uint   `json:"user_id" binding:"required"`
}

type AddExpenseToReportRequest struct {
	UserID    uint `json:"user_id" binding:"required"`
	ExpenseID uint `json:"expense_id" binding:"required"`
}
