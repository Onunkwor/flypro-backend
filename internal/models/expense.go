package models

import "time"

type Expense struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"not null"`
	Amount      float64   `json:"amount" gorm:"not null"`
	AmountUSD   float64   `json:"amount_usd"`
	Currency    string    `json:"currency" gorm:"not null"`
	Category    string    `json:"category" gorm:"not null"`
	Description string    `json:"description"`
	Receipt     string    `json:"receipt"`
	Status      string    `json:"status" gorm:"default:'pending'"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	User User `json:"user" gorm:"foreignKey:UserID"`
}
