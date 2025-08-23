package dto

import "github.com/onunkwor/flypro-backend/internal/utils"

type ConvertCurrencyRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
	From   string  `json:"from" binding:"required,len=3,currency"`
	To     string  `json:"to" binding:"required,len=3,currency"`
}

func (r *ConvertCurrencyRequest) Sanitize() {
	r.From = utils.NormalizeCurrency(r.From)
	r.To = utils.NormalizeCurrency(r.To)
}
