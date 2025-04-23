package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"type:varchar(50);unique;not null" json:"username"`
	Email     string         `gorm:"type:varchar(255);unique;not null" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	CreatedAt time.Time      `gorm:"not null;default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null;default:current_timestamp" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
