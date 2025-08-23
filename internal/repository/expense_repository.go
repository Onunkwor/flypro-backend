package repository

import (
	"errors"

	"github.com/onunkwor/flypro-backend/internal/models"
	"gorm.io/gorm"
)

var (
	ErrExpenseNotFound = errors.New("expense not found")
)

type ExpenseRepository interface {
	Create(expense *models.Expense) error
	GetExpenseByID(id uint) (*models.Expense, error)
	FindAll(filters map[string]interface{}, offset, limit int) ([]models.Expense, error)
	Update(expense *models.Expense) error
	Delete(id uint) error
	UpdateExpenseAmountUSD(expenseID uint, amountUSD float64) error
}

type expenseRepo struct {
	db *gorm.DB
}

func NewExpenseRepository(db *gorm.DB) ExpenseRepository {
	return &expenseRepo{db: db}
}

func (r *expenseRepo) Create(expense *models.Expense) error {
	if err := r.db.Create(expense).Error; err != nil {
		return err
	}
	return r.db.Preload("User").First(expense, expense.ID).Error
}

func (r *expenseRepo) GetExpenseByID(id uint) (*models.Expense, error) {
	var expense models.Expense
	if err := r.db.Preload("User").First(&expense, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrExpenseNotFound
		}
		return nil, err
	}
	return &expense, nil
}

func (r *expenseRepo) Update(expense *models.Expense) error {
	res := r.db.Model(&models.Expense{}).
		Where("id = ?", expense.ID).
		Updates(expense)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrExpenseNotFound
	}
	return nil
}

func (r *expenseRepo) Delete(id uint) error {
	res := r.db.Delete(&models.Expense{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrExpenseNotFound
	}
	return nil
}

func (r *expenseRepo) FindAll(filters map[string]interface{}, offset, limit int) ([]models.Expense, error) {
	var expenses []models.Expense
	q := r.db.Model(&models.Expense{})
	for k, v := range filters {
		q = q.Where(k+" = ?", v)
	}
	if err := q.Offset(offset).Limit(limit).Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

func (r *expenseRepo) UpdateExpenseAmountUSD(expenseID uint, amountUSD float64) error {
	res := r.db.Model(&models.Expense{}).
		Where("id = ?", expenseID).
		Update("amount_usd", amountUSD)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrExpenseNotFound
	}
	return nil

}
