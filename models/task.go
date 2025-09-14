package models

import "time"

type Task struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Category    string     `json:"category"`
	DueDate     *time.Time `json:"due_date"`
	Completed   bool       `json:"completed" gorm:"default:false"`
	UserID      uint       `json:"user_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
