package mocks

import (
	"github.com/onunkwor/flypro-backend/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockExpenseRepo struct {
	mock.Mock
}

func (m *MockExpenseRepo) Create(expense *models.Expense) error {
	args := m.Called(expense)
	return args.Error(0)
}

func (m *MockExpenseRepo) GetExpenseByID(id uint) (*models.Expense, error) {
	args := m.Called(id)
	if e, ok := args.Get(0).(*models.Expense); ok {
		return e, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockExpenseRepo) FindAll(filters map[string]interface{}, offset, limit int) ([]models.Expense, error) {
	args := m.Called(filters, offset, limit)
	if e, ok := args.Get(0).([]models.Expense); ok {
		return e, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockExpenseRepo) Update(expense *models.Expense) error {
	args := m.Called(expense)
	return args.Error(0)
}

func (m *MockExpenseRepo) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockExpenseRepo) UpdateExpenseAmountUSD(expenseID uint, amountUSD float64) error {
	args := m.Called(expenseID, amountUSD)
	return args.Error(0)
}
