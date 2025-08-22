package dto

type ConvertCurrencyRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
	From   string  `json:"from" binding:"required,len=3,currency"`
	To     string  `json:"to" binding:"required,len=3,currency"`
}
