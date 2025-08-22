package services

import (
	"errors"

	"github.com/onunkwor/flypro-backend/internal/models"
	"github.com/onunkwor/flypro-backend/internal/repository"
)

type ExpenseService struct {
	repo repository.ExpenseRepository
}

func NewExpenseService(r repository.ExpenseRepository) *ExpenseService {
	return &ExpenseService{repo: r}
}

func (s *ExpenseService) CreateExpense(expense *models.Expense) error {
	return s.repo.Create(expense)
}

func (s *ExpenseService) GetExpenseByID(id uint) (*models.Expense, error) {
	expense, err := s.repo.GetExpenseByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrExpenseNotFound) {
			return nil, repository.ErrExpenseNotFound
		}
		return nil, err
	}
	return expense, nil
}

func (s *ExpenseService) UpdateExpense(expense *models.Expense) error {
	return s.repo.Update(expense)
}

func (s *ExpenseService) DeleteExpense(id uint) error {
	return s.repo.Delete(id)
}

func (s *ExpenseService) ListExpenses(filters map[string]interface{}, offset, limit int) ([]models.Expense, error) {
	return s.repo.FindAll(filters, offset, limit)
}
