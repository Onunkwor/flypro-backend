package mocks

import (
	"github.com/onunkwor/flypro-backend/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockReportRepo struct {
	mock.Mock
}

func (m *MockReportRepo) Create(report *models.ExpenseReport) error {
	args := m.Called(report)
	return args.Error(0)
}

func (m *MockReportRepo) AddExpense(reportID, userID, expenseID uint) error {
	args := m.Called(reportID, userID, expenseID)
	return args.Error(0)
}
func (m *MockReportRepo) GetByID(id uint) (*models.ExpenseReport, error) {
	args := m.Called(id)
	if r, ok := args.Get(0).(*models.ExpenseReport); ok {
		return r, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockReportRepo) List(userID uint, offset, limit int) ([]models.ExpenseReport, error) {
	args := m.Called(userID, offset, limit)
	if r, ok := args.Get(0).([]models.ExpenseReport); ok {
		return r, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockReportRepo) Update(report *models.ExpenseReport) error {
	args := m.Called(report)
	return args.Error(0)
}
func (m *MockReportRepo) UpdateReportTotal(reportID uint, total float64) error {
	args := m.Called(reportID, total)
	return args.Error(0)
}
