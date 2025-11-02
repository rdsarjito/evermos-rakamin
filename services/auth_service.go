package services

import (
	"errors"

	"github.com/rdsarjito/evermos-rakamin/constants"
	"github.com/rdsarjito/evermos-rakamin/domain"
	"github.com/rdsarjito/evermos-rakamin/domain/dto"
	"github.com/rdsarjito/evermos-rakamin/helpers"
	"github.com/rdsarjito/evermos-rakamin/repositories"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo *repositories.UserRepository
}

// NewAuthService creates a new auth service
func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: repositories.NewUserRepository(),
	}
}

// Register registers a new user
func (s *AuthService) Register(req dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Hash password
	hashedPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Set default role if not provided
	role := req.Role
	if role == "" {
		role = constants.RoleUser
	}

	// Create user
	user := &domain.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         role,
	}

	if err := s.userRepo.Create(user); err != nil {
		if err == repositories.ErrUserExists {
			return nil, err
		}
		return nil, err
	}

	// Generate token
	token, err := helpers.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token:     token,
		TokenType: "Bearer",
		ExpiresIn: 24 * 60 * 60, // 24 hours in seconds
		User: dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}

// Login authenticates a user
func (s *AuthService) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if err == repositories.ErrUserNotFound {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Verify password
	if !helpers.CheckPassword(req.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	// Generate token
	token, err := helpers.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token:     token,
		TokenType: "Bearer",
		ExpiresIn: 24 * 60 * 60, // 24 hours in seconds
		User: dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}

// GetUserProfile gets user profile by ID
func (s *AuthService) GetUserProfile(userID uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

