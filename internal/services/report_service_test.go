package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/onunkwor/flypro-backend/internal/models"
	"github.com/onunkwor/flypro-backend/internal/repository"
	"github.com/onunkwor/flypro-backend/internal/repository/mocks"
	"github.com/onunkwor/flypro-backend/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCurrencySvc struct {
	mock.Mock
}

func (m *MockCurrencySvc) Convert(ctx context.Context, amount float64, from, to string) (float64, error) {
	args := m.Called(ctx, amount, from, to)
	return args.Get(0).(float64), args.Error(1)
}
func TestSubmitReport_Success(t *testing.T) {
	mockRepo := new(mocks.MockReportRepo)
	mockExpRepo := new(mocks.MockExpenseRepo)

	svc := services.NewReportService(mockRepo, mockExpRepo, nil)

	report := &models.ExpenseReport{
		ID:     1,
		UserID: 42,
		Status: "draft",
	}
	mockRepo.On("GetByID", uint(1)).Return(report, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.ExpenseReport")).Return(nil)

	err := svc.SubmitReport(1, 42)

	assert.NoError(t, err)
	assert.Equal(t, "submitted", report.Status)
	mockRepo.AssertExpectations(t)
}

func TestSubmitReport_InvalidOwnership(t *testing.T) {
	mockRepo := new(mocks.MockReportRepo)
	mockExpRepo := new(mocks.MockExpenseRepo)

	svc := services.NewReportService(mockRepo, mockExpRepo, nil)

	report := &models.ExpenseReport{
		ID:     1,
		UserID: 99,
		Status: "draft",
	}
	mockRepo.On("GetByID", uint(1)).Return(report, nil)

	err := svc.SubmitReport(1, 42)
	assert.Error(t, err)
	assert.Equal(t, repository.ErrInvalidOwnership, err)
	mockRepo.AssertExpectations(t)
}

func TestSubmitReport_InvalidState(t *testing.T) {
	mockRepo := new(mocks.MockReportRepo)
	mockExpRepo := new(mocks.MockExpenseRepo)
	svc := services.NewReportService(mockRepo, mockExpRepo, nil)

	report := &models.ExpenseReport{
		ID:     1,
		UserID: 42,
		Status: "submitted",
	}

	mockRepo.On("GetByID", uint(1)).Return(report, nil)

	err := svc.SubmitReport(1, 42)

	assert.Error(t, err)
	assert.Equal(t, services.ErrInvalidReportState, err)
	mockRepo.AssertExpectations(t)
}

func TestListReports_ConvertsToUSD(t *testing.T) {
	mockRepo := new(mocks.MockReportRepo)
	mockExpRepo := new(mocks.MockExpenseRepo)
	mockCur := new(MockCurrencySvc)

	svc := services.NewReportService(mockRepo, mockExpRepo, mockCur)

	reports := []models.ExpenseReport{
		{
			ID:     1,
			UserID: 42,
			Expenses: []models.Expense{
				{ID: 10, Amount: 100, Currency: "EUR"},
				{ID: 11, Amount: 50, Currency: "USD"},
			},
		},
	}
	mockRepo.On("List", uint(42), 0, 10).Return(reports, nil)
	mockCur.On("Convert", mock.Anything, 100.0, "EUR", "USD").Return(120.0, nil)
	mockExpRepo.On("UpdateExpenseAmountUSD", uint(10), 120.0).Return(nil)
	mockRepo.On("UpdateReportTotal", uint(1), 170.0).Return(nil)
	result, err := svc.ListReports(42, 0, 10)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, 170.0, result[0].Total)
	assert.Equal(t, 120.0, result[0].Expenses[0].AmountUSD)
	assert.Equal(t, 50.0, result[0].Expenses[1].AmountUSD)

	mockRepo.AssertExpectations(t)
	mockExpRepo.AssertExpectations(t)
	mockCur.AssertExpectations(t)
}

func TestListReports_FailsOnCurrencyError(t *testing.T) {
	mockRepo := new(mocks.MockReportRepo)
	mockExpRepo := new(mocks.MockExpenseRepo)
	mockCur := new(MockCurrencySvc)

	svc := services.NewReportService(mockRepo, mockExpRepo, mockCur)

	reports := []models.ExpenseReport{
		{
			ID:     1,
			UserID: 42,
			Expenses: []models.Expense{
				{ID: 10, Amount: 100, Currency: "EUR"},
			},
		},
	}

	mockRepo.On("List", uint(42), 0, 10).Return(reports, nil)
	mockCur.On("Convert", mock.Anything, 100.0, "EUR", "USD").Return(0.0, errors.New("API down"))

	result, err := svc.ListReports(42, 0, 10)

	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
	mockCur.AssertExpectations(t)
}
