package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/onunkwor/flypro-backend/internal/dto"
	"github.com/onunkwor/flypro-backend/internal/services"
	"github.com/onunkwor/flypro-backend/internal/utils"
)

type CurrencyHandler struct {
	svc *services.CurrencyService
}

func NewCurrencyHandler(svc *services.CurrencyService) *CurrencyHandler {
	return &CurrencyHandler{svc: svc}
}

func (h *CurrencyHandler) ConvertCurrency(c *gin.Context) {
	var req dto.ConvertCurrencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		formatted := utils.FormatValidationError(err)
		utils.ValidationErrorResponse(c, formatted)
		return
	}
	converted, err := h.svc.Convert(context.Background(), req.Amount, req.From, req.To)
	if err != nil {
		utils.InternalServerErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"from":      req.From,
		"to":        req.To,
		"amount":    req.Amount,
		"converted": converted,
	})
}
