package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/onunkwor/flypro-backend/internal/dto"
	"github.com/onunkwor/flypro-backend/internal/models"
	"github.com/onunkwor/flypro-backend/internal/repository"
	"github.com/onunkwor/flypro-backend/internal/services"
	"github.com/onunkwor/flypro-backend/internal/utils"
)

type ExpenseHandler struct {
	svc *services.ExpenseService
}

func NewExpenseHandler(svc *services.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{svc: svc}
}

// POST /api/expenses
func (h *ExpenseHandler) CreateExpense(c *gin.Context) {
	var req dto.CreateExpenseRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		formatted := utils.FormatValidationError(err)
		utils.ValidationErrorResponse(c, formatted)
		return
	}
	req.Sanitize()
	expense := &models.Expense{
		UserID:      req.UserID,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Category:    req.Category,
		Description: req.Description,
	}

	if err := h.svc.CreateExpense(expense); err != nil {
		utils.InternalServerErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"expense": dto.ToExpenseResponse(expense)})
}

func (h *ExpenseHandler) GetExpenseByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expense ID"})
		return
	}

	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	expense, err := h.svc.GetExpenseByID(uint(id))
	if err != nil {
		if err == repository.ErrExpenseNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "expense not found"})
			return
		}
		utils.InternalServerErrorResponse(c, err)
		return
	}

	if fmt.Sprintf("%d", expense.UserID) != userID {
		utils.ForbiddenResponse(c, "not authorized to access this expense")
		return
	}

	c.JSON(http.StatusOK, gin.H{"expense": dto.ToExpenseResponse(expense)})
}

func (h *ExpenseHandler) UpdateExpense(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expense ID"})
		return
	}

	var req dto.UpdateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		formatted := utils.FormatValidationError(err)
		utils.ValidationErrorResponse(c, formatted)
		return
	}
	req.Sanitize()

	expense, err := h.svc.GetExpenseByID(uint(id))
	if err != nil {
		if err == repository.ErrExpenseNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "expense not found"})
			return
		}
		utils.InternalServerErrorResponse(c, err)
		return
	}

	if expense.UserID != req.UserID {
		utils.ForbiddenResponse(c, "not authorized to update this expense")
		return
	}

	expense.Amount = req.Amount
	expense.Currency = req.Currency
	expense.Category = req.Category
	expense.Description = req.Description

	err = h.svc.UpdateExpense(expense)
	if err != nil {
		utils.InternalServerErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"expense": dto.ToExpenseResponse(expense)})
}

func (h *ExpenseHandler) DeleteExpense(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expense ID"})
		return
	}

	var req struct {
		UserID uint `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, utils.FormatValidationError(err))
		return
	}

	expense, err := h.svc.GetExpenseByID(uint(id))
	if err != nil {
		if err == repository.ErrExpenseNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "expense not found"})
			return
		}
		utils.InternalServerErrorResponse(c, err)
		return
	}

	if expense.UserID != req.UserID {
		utils.ForbiddenResponse(c, "not authorized to delete this expense")
		return
	}

	if err := h.svc.DeleteExpense(uint(id)); err != nil {
		if err == repository.ErrExpenseNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "expense not found"})
			return
		}
		utils.InternalServerErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "expense deleted successfully"})
}

func (h *ExpenseHandler) ListExpenses(c *gin.Context) {
	filters := make(map[string]interface{})
	if userID := c.Query("user_id"); userID != "" {
		if uid, err := strconv.Atoi(userID); err == nil && uid > 0 {
			filters["user_id"] = uint(uid)
		}
	}
	if category := c.Query("category"); category != "" {
		filters["category"] = utils.NormalizeCategory(category)
	}
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offsetStr := c.DefaultQuery("offset", "0")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	expenses, err := h.svc.ListExpenses(filters, offset, limit, context.Background())
	if err != nil {
		utils.InternalServerErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"expenses": dto.ToExpenseResponses(expenses)})

}
