package domain

import (
	"time"

	"gorm.io/gorm"
)

// Category represents category entity
type Category struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"uniqueIndex;not null"`
	Description string         `json:"description" gorm:"type:text"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name for Category
func (Category) TableName() string {
	return "categories"
}
