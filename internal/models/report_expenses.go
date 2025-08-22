package models

// ReportExpense is the join table linking reports and expenses
type ReportExpense struct {
	ReportID  uint `gorm:"primaryKey"`
	ExpenseID uint `gorm:"primaryKey"`
}
