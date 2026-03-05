package service

import (
	"context"
	"fmt"
	"time"

	"go-web-demo/internal/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, username, email string) (*repository.User, error)
	GetUser(ctx context.Context, id int64) (*repository.User, error)
	GetUserByEmail(ctx context.Context, email string) (*repository.User, error)
	ListUsers(ctx context.Context, page, pageSize int) ([]*repository.User, error)
	UpdateUser(ctx context.Context, id int64, username, email string) (*repository.User, error)
	DeleteUser(ctx context.Context, id int64) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(ctx context.Context, username, email string) (*repository.User, error) {
	if username == "" {
		return nil, fmt.Errorf("username is required")
	}
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}

	user := &repository.User{
		Username:  username,
		Email:     email,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *userService) GetUser(ctx context.Context, id int64) (*repository.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid user id")
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*repository.User, error) {
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

func (s *userService) ListUsers(ctx context.Context, page, pageSize int) ([]*repository.User, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	users, err := s.userRepo.List(ctx, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

func (s *userService) UpdateUser(ctx context.Context, id int64, username, email string) (*repository.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid user id")
	}
	if username == "" && email == "" {
		return nil, fmt.Errorf("username or email is required")
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if username != "" {
		user.Username = username
	}
	if email != "" {
		user.Email = email
	}
	user.UpdatedAt = time.Now().UTC()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid user id")
	}

	if err := s.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
