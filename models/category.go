package models

import "time"

type Category struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Name      string    `gorm:"unique;not null" json:"name"`
    Tasks     []Task    `json:"tasks"` // 1-n: Một Category có nhiều Task
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}