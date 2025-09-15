package models

import "time"

type Task struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CategoryID  uint      `json:"category_id"` // Thay Category string bằng CategoryID
	Category    Category  `json:"category"`    // Thêm quan hệ với Category
	DueDate     time.Time `json:"due_date"`
	Completed   bool      `json:"completed"`
	UserID      uint      `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
