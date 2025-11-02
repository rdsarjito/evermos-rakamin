package domain

import (
	"time"

	"gorm.io/gorm"
)

// User represents user entity
type User struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Name         string         `json:"name" gorm:"not null"`
	Email        string         `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string         `json:"-" gorm:"column:password_hash;not null"`
	Role         string         `json:"role" gorm:"default:'user'"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name for User
func (User) TableName() string {
	return "users"
}
