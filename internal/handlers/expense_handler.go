package handlers

import (
	"context"
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

	c.JSON(http.StatusCreated, gin.H{"expense": expense})
}

func (h *ExpenseHandler) GetExpenseByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expense ID"})
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

	c.JSON(http.StatusOK, gin.H{"expense": expense})
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
	expense := &models.Expense{
		ID:          uint(id),
		UserID:      req.UserID,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Category:    req.Category,
		Description: req.Description,
	}

	if err := h.svc.UpdateExpense(expense); err != nil {
		utils.InternalServerErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"expense": expense})
}

func (h *ExpenseHandler) DeleteExpense(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expense ID"})
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
	// optional filters
	filters := make(map[string]interface{})
	if userID := c.Query("user_id"); userID != "" {
		filters["user_id"] = userID
	}
	if category := c.Query("category"); category != "" {
		filters["category"] = category
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	expenses, err := h.svc.ListExpenses(filters, offset, limit, context.Background())
	if err != nil {
		utils.InternalServerErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"expenses": expenses})
}
