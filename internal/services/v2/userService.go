package services

import (
	"context"
	"errors"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/repositories/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/pkg"
)

type UserService interface {
	Login(ctx context.Context, email, password string) (string, error)
	Register(ctx context.Context, user models.User) (*models.User, error)
	Activation(ctx context.Context, email string) (*models.User, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewServiceUser(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (svc *userService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := svc.repo.GetByUsername(ctx, email)
	if err != nil || user == nil {
		return "", errors.New("Invalid credentials")
	}

	if !pkg.CheckPassword(user.Password, password) {
		return "", errors.New("Invalid credentials")
	}

	token, err := pkg.CreateToken(user.Email, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (svc *userService) Register(ctx context.Context, user models.User) (*models.User, error) {
	existing, err := svc.repo.GetByUsername(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return nil, errors.New("Email already used")
	}

	hashed, err := pkg.HashPassword(user.Password)
	if err != nil {
		return nil, errors.New("Failed to hash password")
	}
	user.Password = hashed

	newUser, err := svc.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (svc *userService) Activation(ctx context.Context, email string) (*models.User, error) {
	existing, err := svc.repo.GetByUsername(ctx, email)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return nil, errors.New("User not found")
	}

	if existing.Status == "aktif" {
		return nil, errors.New("User is already active")
	}

	updatedUser, err := svc.repo.UpdateStatus(ctx, models.User{
		Email:  email,
		Status: "aktif",
	})
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}
