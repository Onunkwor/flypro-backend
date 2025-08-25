package models

import "time"

type ExpenseReport struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	Title     string    `json:"title" gorm:"not null"`
	Status    string    `json:"status" gorm:"default:'draft'"`
	Total     float64   `json:"total"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User     *User     `json:"user" gorm:"foreignKey:UserID"`
	Expenses []Expense `json:"expenses" gorm:"many2many:report_expenses;joinForeignKey:ReportID;joinReferences:ExpenseID"`
}
