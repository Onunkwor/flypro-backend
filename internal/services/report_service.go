package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/onunkwor/flypro-backend/internal/models"
	"github.com/onunkwor/flypro-backend/internal/repository"
)

var (
	ErrInvalidReportState = errors.New("report cannot be submitted in current state")
)

type ReportService struct {
	repo        repository.ReportRepository
	expenseRepo repository.ExpenseRepository
	currencySvc CurrencyConverter
	redis       *redis.Client
}

func NewReportService(r repository.ReportRepository, exp repository.ExpenseRepository, cur CurrencyConverter, redis *redis.Client) *ReportService {
	return &ReportService{
		repo:        r,
		expenseRepo: exp,
		currencySvc: cur,
		redis:       redis,
	}
}

// service.go
func (s *ReportService) CreateReport(report *models.ExpenseReport) error {
	exists, err := s.repo.UserExists(report.UserID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("user with id %d does not exist", report.UserID)
	}

	if err := s.repo.Create(report); err != nil {
		return err
	}

	if err := s.repo.LoadReportRelations(report); err != nil {
		return err
	}

	return nil
}

func (s *ReportService) AddExpense(reportID, userID, expenseID uint) error {
	return s.repo.AddExpense(reportID, userID, expenseID)
}

func (s *ReportService) ListReports(userID uint, offset, limit int) ([]models.ExpenseReport, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("report_summaries:%d:%d:%d", userID, offset, limit)

	// if s.redis != nil {
	// 	if cached, err := s.redis.Get(ctx, cacheKey).Result(); err == nil {
	// 		var reports []models.ExpenseReport
	// 		if json.Unmarshal([]byte(cached), &reports) == nil {
	// 			return reports, nil
	// 		}
	// 	}
	// }

	reports, err := s.repo.List(userID, offset, limit)
	if err != nil {
		return nil, err
	}

	for i := range reports {
		totalUSD := 0.0
		for j := range reports[i].Expenses {
			exp := &reports[i].Expenses[j]

			if exp.AmountUSD == 0 && exp.Currency != "USD" {
				converted, err := s.currencySvc.Convert(ctx, exp.Amount, exp.Currency, "USD")
				if err != nil {
					return nil, fmt.Errorf("failed to convert expense %d: %w", exp.ID, err)
				}
				exp.AmountUSD = converted
				_ = s.expenseRepo.UpdateExpenseAmountUSD(exp.ID, converted)
			} else if exp.Currency == "USD" {
				exp.AmountUSD = exp.Amount
			}

			totalUSD += exp.AmountUSD
		}
		reports[i].Total = totalUSD
		_ = s.repo.UpdateReportTotal(reports[i].ID, totalUSD)
	}

	if s.redis != nil {
		if data, err := json.Marshal(reports); err == nil {
			_ = s.redis.Set(ctx, cacheKey, data, 30*time.Minute).Err()
		}
	}

	return reports, nil
}

func (s *ReportService) SubmitReport(id uint, userID uint) error {
	report, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if report.UserID != userID {
		return repository.ErrInvalidOwnership
	}
	if report.Status != "draft" {
		return ErrInvalidReportState
	}
	report.Status = "submitted"
	return s.repo.Update(report)
}
