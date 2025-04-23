package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID          int            `gorm:"primaryKey" json:"id"`
	UserID      int            `gorm:"not null" json:"user_id"`
	User        User           `gorm:"foreignKey:UserID" json:"-"`
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Status      string         `gorm:"type:varchar(50);not null" json:"status"`
	DueDate     *time.Time     `gorm:"type:timestamp" json:"due_date"`
	CreatedAt   time.Time      `gorm:"not null;default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"not null;default:current_timestamp" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
