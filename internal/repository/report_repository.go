package repository

import (
	"errors"

	"github.com/onunkwor/flypro-backend/internal/models"
	"gorm.io/gorm"
)

var (
	ErrReportNotFound   = errors.New("report not found")
	ErrInvalidOwnership = errors.New("invalid ownership")
)

type ReportRepository interface {
	Create(report *models.ExpenseReport) error
	AddExpense(reportID, userID, expenseID uint) error
	GetByID(id uint) (*models.ExpenseReport, error)
	List(userID uint, offset, limit int) ([]models.ExpenseReport, error)
	Update(report *models.ExpenseReport) error
	UpdateReportTotal(reportID uint, total float64) error
}

type reportRepo struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepo{db: db}
}

func (r *reportRepo) Create(report *models.ExpenseReport) error {
	return r.db.Create(report).Error
}

func (r *reportRepo) AddExpense(reportID, userID, expenseID uint) error {
	// Check ownership
	var expense models.Expense
	if err := r.db.First(&expense, expenseID).Error; err != nil {
		return ErrReportNotFound
	}
	if expense.UserID != userID {
		return ErrInvalidOwnership
	}

	link := models.ReportExpense{ReportID: reportID, ExpenseID: expenseID}
	return r.db.Create(&link).Error
}

func (r *reportRepo) GetByID(id uint) (*models.ExpenseReport, error) {
	var report models.ExpenseReport
	if err := r.db.Preload("Expenses").First(&report, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReportNotFound
		}
		return nil, err
	}
	return &report, nil
}

func (r *reportRepo) List(userID uint, offset, limit int) ([]models.ExpenseReport, error) {
	var reports []models.ExpenseReport
	if err := r.db.Where("user_id = ?", userID).
		Offset(offset).Limit(limit).
		Preload("Expenses").
		Find(&reports).Error; err != nil {
		return nil, err
	}
	return reports, nil
}

func (r *reportRepo) Update(report *models.ExpenseReport) error {
	return r.db.Save(report).Error
}

func (r *reportRepo) UpdateReportTotal(reportID uint, total float64) error {
	return r.db.Model(&models.ExpenseReport{}).
		Where("id = ?", reportID).
		Update("total", total).Error
}
