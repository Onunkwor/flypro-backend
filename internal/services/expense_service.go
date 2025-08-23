package services

import (
	"context"
	"encoding/json"
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
	return s.repo.GetExpenseByID(id)

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
	cached, err := s.redis.Get(ctx, key).Result()
	if err == nil {
		if unmarshalErr := json.Unmarshal([]byte(cached), &expenses); unmarshalErr == nil {
			return expenses, nil
		}
	} else if err != redis.Nil {
		return nil, err
	}

	expenses, e := s.repo.FindAll(filters, offset, limit)
	if e != nil {
		return nil, e
	}
	if data, err := json.Marshal(expenses); err == nil {
		_ = s.redis.Set(ctx, key, data, 30*time.Minute).Err()
	}

	return expenses, nil
}
