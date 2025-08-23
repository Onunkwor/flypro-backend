package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/onunkwor/flypro-backend/internal/dto"
	"github.com/onunkwor/flypro-backend/internal/models"
	"github.com/onunkwor/flypro-backend/internal/repository"
	"github.com/onunkwor/flypro-backend/internal/services"
	"github.com/onunkwor/flypro-backend/internal/utils"
)

type ReportHandler struct {
	svc *services.ReportService
}

func NewReportHandler(svc *services.ReportService) *ReportHandler {
	return &ReportHandler{svc: svc}
}

func (h *ReportHandler) CreateReport(c *gin.Context) {
	var req dto.CreateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, utils.FormatValidationError(err))
		return
	}
	report := &models.ExpenseReport{
		Title:  req.Title,
		UserID: req.UserID,
		Status: "draft",
	}
	if err := h.svc.CreateReport(report); err != nil {
		utils.InternalServerErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"report": dto.ToReportDTO(*report)})
}

func (h *ReportHandler) AddExpense(c *gin.Context) {
	reportIDParam := c.Param("id")
	if reportIDParam == "" {
		utils.BadRequestResponse(c, "report ID is required")
		return
	}

	reportID, err := strconv.Atoi(reportIDParam)
	if err != nil {
		utils.BadRequestResponse(c, "report ID must be a valid number")
		return
	}
	var req dto.AddExpenseToReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, utils.FormatValidationError(err))
		return
	}
	err = h.svc.AddExpense(uint(reportID), req.UserID, req.ExpenseID)
	if err != nil {
		if errors.Is(err, repository.ErrReportNotFound) {
			utils.NotFoundResponse(c, "report not found")
		}
		if errors.Is(err, repository.ErrExpenseNotFound) {
			utils.NotFoundResponse(c, "expense not found")
		}
		if errors.Is(err, repository.ErrInvalidOwnership) {
			utils.ForbiddenResponse(c, "you do not own this report or expense")
		}
		utils.InternalServerErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "expense added to report"})
}

func (h *ReportHandler) ListReports(c *gin.Context) {
	userIDParam := c.Query("user_id")
	if userIDParam == "" {
		utils.BadRequestResponse(c, "user_id is required")
		return
	}

	userID, err := strconv.Atoi(userIDParam)
	if err != nil || userID <= 0 {
		utils.BadRequestResponse(c, "invalid user ID")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	reports, err := h.svc.ListReports(uint(userID), offset, limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"reports": dto.ToReportDTOs(reports)})
}

func (h *ReportHandler) SubmitReport(c *gin.Context) {
	reportIDParam := c.Param("id")
	reportID, err := strconv.Atoi(reportIDParam)
	if err != nil || reportID <= 0 {
		utils.BadRequestResponse(c, "invalid report ID")
		return
	}

	userIDParam := c.Query("user_id")
	if userIDParam == "" {
		utils.BadRequestResponse(c, "user_id is required")
		return
	}
	userID, err := strconv.Atoi(userIDParam)
	if err != nil || userID <= 0 {
		utils.BadRequestResponse(c, "invalid user ID")
		return
	}

	err = h.svc.SubmitReport(uint(reportID), uint(userID))
	if err != nil {

		switch {
		case errors.Is(err, repository.ErrInvalidOwnership):
			utils.ForbiddenResponse(c, "you do not own this report")
		case errors.Is(err, services.ErrInvalidReportState):
			utils.BadRequestResponse(c, "report cannot be submitted in its current state")
		default:
			utils.InternalServerErrorResponse(c, err)
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "report submitted"})
}
