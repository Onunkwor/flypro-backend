package services_test

import (
	"errors"
	"testing"

	"github.com/onunkwor/flypro-backend/internal/config"
	"github.com/onunkwor/flypro-backend/internal/models"
	"github.com/onunkwor/flypro-backend/internal/repository"
	"github.com/onunkwor/flypro-backend/internal/repository/mocks"
	"github.com/onunkwor/flypro-backend/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestCreateExpense_Success(t *testing.T) {
	mockRepo := new(mocks.MockExpenseRepo)
	svc := services.NewExpenseService(mockRepo, config.Redis)

	expense := &models.Expense{
		ID:       1,
		UserID:   10,
		Amount:   100.50,
		Currency: "USD",
		Category: "travel",
	}

	// Expect Create to be called once and return no error
	mockRepo.On("Create", expense).Return(nil)

	err := svc.CreateExpense(expense)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateExpense_Failure(t *testing.T) {
	mockRepo := new(mocks.MockExpenseRepo)
	svc := services.NewExpenseService(mockRepo, config.Redis)

	expense := &models.Expense{ID: 2, UserID: 20, Amount: 200, Currency: "EUR", Category: "meals"}
	expectedErr := errors.New("db error")

	// Simulate repo failure
	mockRepo.On("Create", expense).Return(expectedErr)

	err := svc.CreateExpense(expense)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

func TestGetExpenseByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockExpenseRepo)
	svc := services.NewExpenseService(mockRepo, config.Redis)

	expenseID := uint(3)
	expectedExpense := &models.Expense{ID: expenseID, UserID: 30, Amount: 300, Currency: "NGN", Category: "office"}

	mockRepo.On("GetExpenseByID", expenseID).Return(expectedExpense, nil)

	expense, err := svc.GetExpenseByID(expenseID)

	assert.NoError(t, err)
	assert.Equal(t, expectedExpense, expense)
	mockRepo.AssertExpectations(t)
}

func TestGetExpenseByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockExpenseRepo)
	svc := services.NewExpenseService(mockRepo, config.Redis)

	expenseID := uint(4)

	mockRepo.On("GetExpenseByID", expenseID).Return(nil, repository.ErrExpenseNotFound)

	expense, err := svc.GetExpenseByID(expenseID)

	assert.Error(t, err)
	assert.Nil(t, expense)
	assert.Equal(t, repository.ErrExpenseNotFound, err)
	mockRepo.AssertExpectations(t)
}
