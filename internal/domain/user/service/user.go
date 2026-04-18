package service

import (
	"context"

	"finance/internal/common/errors"
	"finance/internal/domain/user/entity"
	"finance/internal/domain/user/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserService handles user business logic
type UserService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, username, email, phone string) (*entity.User, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, errors.ErrUserAlreadyExists
	}

	user := entity.NewUser(username, email, phone)
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID gets a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}
	return user, nil
}

// RechargeUser recharges a user's balance
func (s *UserService) RechargeUser(ctx context.Context, userID int64, amount float64) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return errors.ErrUserNotFound
	}

	if amount <= 0 {
		return errors.ErrInvalidRechargeAmount
	}

	user.Recharge(amount)
	return s.userRepo.Update(ctx, user)
}

// ListUsers lists users with pagination
func (s *UserService) ListUsers(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	return s.userRepo.List(ctx, limit, offset)
}

// HashPassword hashes a plain text password using bcrypt
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPasswordHash compares a plain text password with a hashed password
func CheckPasswordHash(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
