package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/onunkwor/flypro-backend/internal/models"
	"github.com/onunkwor/flypro-backend/internal/repository"
)

var (
	ErrInvalidReportState = errors.New("report cannot be submitted in current state")
)

type ReportService struct {
	repo        repository.ReportRepository
	expenseRepo repository.ExpenseRepository
	currencySvc *CurrencyService
}

func NewReportService(r repository.ReportRepository, exp repository.ExpenseRepository, cur *CurrencyService) *ReportService {
	return &ReportService{repo: r, expenseRepo: exp, currencySvc: cur}
}

func (s *ReportService) CreateReport(report *models.ExpenseReport) error {

	return s.repo.Create(report)
}

func (s *ReportService) AddExpense(reportID, userID, expenseID uint) error {
	return s.repo.AddExpense(reportID, userID, expenseID)
}

func (s *ReportService) ListReports(userID uint, offset, limit int) ([]models.ExpenseReport, error) {
	reports, err := s.repo.List(userID, offset, limit)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	for i := range reports {
		totalUSD := 0.0
		for j := range reports[i].Expenses {
			exp := &reports[i].Expenses[j]

			// If AmountUSD is not yet set, convert it
			if exp.AmountUSD == 0 && exp.Currency != "USD" {
				converted, err := s.currencySvc.Convert(ctx, exp.Amount, exp.Currency, "USD")
				if err != nil {
					return nil, fmt.Errorf("failed to convert expense %d: %w", exp.ID, err)
				}
				exp.AmountUSD = converted

				s.expenseRepo.UpdateExpenseAmountUSD(exp.ID, converted)
			} else if exp.Currency == "USD" {
				exp.AmountUSD = exp.Amount
			}

			totalUSD += exp.AmountUSD
		}
		reports[i].Total = totalUSD

		s.repo.UpdateReportTotal(reports[i].ID, totalUSD)
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
