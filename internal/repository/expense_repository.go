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
}

type expenseRepo struct {
	db *gorm.DB
}

func NewExpenseRepository(db *gorm.DB) ExpenseRepository {
	return &expenseRepo{db: db}
}

func (r *expenseRepo) Create(expense *models.Expense) error {
	return r.db.Create(expense).Error
}

func (r *expenseRepo) GetExpenseByID(id uint) (*models.Expense, error) {
	var expense models.Expense
	if err := r.db.First(&expense, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrExpenseNotFound
		}
		return nil, err
	}
	return &expense, nil
}

func (r *expenseRepo) Update(expense *models.Expense) error {
	return r.db.Save(expense).Error
}

func (r *expenseRepo) Delete(id uint) error {
	if err := r.db.Delete(&models.Expense{}, id).Error; err != nil {
		return err
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
