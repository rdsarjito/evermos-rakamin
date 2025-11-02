package repositories

import (
	"errors"

	"github.com/rdsarjito/evermos-rakamin/domain"
	"gorm.io/gorm"
)

// UserRepository handles user database operations
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: GetDB(),
	}
}

// Create creates a new user
func (r *UserRepository) Create(user *domain.User) error {
	// Check if email already exists
	var existingUser domain.User
	if err := r.db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return ErrUserExists
	}

	if err := r.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

// FindByEmail finds user by email
func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindByID finds user by ID
func (r *UserRepository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

