package dto

type CreateExpenseRequest struct {
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Currency    string  `json:"currency" binding:"required,len=3,currency"`
	Category    string  `json:"category" binding:"required,oneof=travel meals office supplies"`
	Description string  `json:"description" binding:"max=500"`
	UserID      uint    `json:"user_id" binding:"required"`
}

type UpdateExpenseRequest struct {
	UserID      uint    `json:"user_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Currency    string  `json:"currency" binding:"required,len=3,currency"`
	Category    string  `json:"category" binding:"required,oneof=travel meals office supplies"`
	Description string  `json:"description"`
}
