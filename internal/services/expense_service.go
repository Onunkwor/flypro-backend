package services

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/onunkwor/flypro-backend/internal/models"
	"github.com/onunkwor/flypro-backend/internal/repository"
	"github.com/onunkwor/flypro-backend/internal/utils"
)

type ExpenseService struct {
	repo  repository.ExpenseRepository
	redis *redis.Client
}

func NewExpenseService(r repository.ExpenseRepository, redis *redis.Client) *ExpenseService {
	return &ExpenseService{repo: r, redis: redis}
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

func (s *ExpenseService) ListExpenses(filters map[string]interface{}, offset, limit int, ctx context.Context) ([]models.Expense, error) {
	key := utils.MakeCacheKey(filters, offset, limit)
	var expenses []models.Expense
	if cached, err := s.redis.Get(ctx, key).Result(); err == nil {
		if err := json.Unmarshal([]byte(cached), &expenses); err == nil {
			return expenses, nil
		}
	}
	expenses, err := s.repo.FindAll(filters, offset, limit)
	if err != nil {
		return nil, err
	}
	if data, err := json.Marshal(expenses); err == nil {
		_ = s.redis.Set(ctx, key, data, 2*time.Minute).Err()
	}

	return expenses, nil
}
