package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/onunkwor/flypro-backend/internal/dto"
	"github.com/onunkwor/flypro-backend/internal/models"
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
	c.JSON(http.StatusCreated, gin.H{"report": report})
}

func (h *ReportHandler) AddExpense(c *gin.Context) {
	reportIDParam := c.Param("id")
	if reportIDParam == "" {
		utils.BadRequestResponse(c, "report ID is required")
		return
	}

	// convert string to int
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
	if err := h.svc.AddExpense(uint(reportID), req.UserID, req.ExpenseID); err != nil {
		utils.InternalServerErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "expense added to report"})
}

func (h *ReportHandler) ListReports(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	reports, err := h.svc.ListReports(uint(userID), offset, limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"reports": reports})
}

func (h *ReportHandler) SubmitReport(c *gin.Context) {
	reportID, _ := strconv.Atoi(c.Param("id"))
	userID, _ := strconv.Atoi(c.Query("user_id"))

	if err := h.svc.SubmitReport(uint(reportID), uint(userID)); err != nil {
		utils.InternalServerErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "report submitted"})
}
