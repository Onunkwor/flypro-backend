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

func TestDeleteExpense_Success(t *testing.T) {
	mockRepo := new(mocks.MockExpenseRepo)
	svc := services.NewExpenseService(mockRepo, nil)
	expenseID := uint(1)

	mockRepo.On("Delete", expenseID).Return(nil)

	err := svc.DeleteExpense(expenseID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteExpense_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockExpenseRepo)
	svc := services.NewExpenseService(mockRepo, nil)
	expenseID := uint(99)

	mockRepo.On("Delete", expenseID).Return(repository.ErrExpenseNotFound)

	err := svc.DeleteExpense(expenseID)

	assert.ErrorIs(t, err, repository.ErrExpenseNotFound)
	mockRepo.AssertExpectations(t)
}

func TestUpdateExpense_Success(t *testing.T) {
	mockRepo := new(mocks.MockExpenseRepo)
	svc := services.NewExpenseService(mockRepo, nil)
	exp := &models.Expense{ID: 1, Amount: 200}

	mockRepo.On("Update", exp).Return(nil)

	err := svc.UpdateExpense(exp)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateExpense_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockExpenseRepo)
	svc := services.NewExpenseService(mockRepo, nil)
	exp := &models.Expense{ID: 999, Amount: 300}

	mockRepo.On("Update", exp).Return(repository.ErrExpenseNotFound)

	err := svc.UpdateExpense(exp)

	assert.ErrorIs(t, err, repository.ErrExpenseNotFound)
	mockRepo.AssertExpectations(t)
}
