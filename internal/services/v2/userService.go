package services

import (
	"context"
	"errors"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/repositories/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/pkg"
)

type UserService interface {
	GetAll(ctx context.Context, page, perPage int) (*UserPagination, error)
	Create(ctx context.Context, user models.User) (*models.User, error)
	GetDetail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, email, password, role, status string) (*models.User, error)
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

type UserPagination struct {
	Data       []*models.User `json:"data"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	PerPage    int            `json:"per_page"`
	TotalPages int            `json:"total_pages"`
}

func (svc *userService) GetAll(ctx context.Context, page, perPage int) (*UserPagination, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	users, err := svc.repo.GetAllUsers(ctx, perPage, offset)
	if err != nil {
		return nil, err
	}

	total, err := svc.repo.GetTotalUsers(ctx)
	if err != nil {
		return nil, err
	}

	totalPages := total / perPage
	if total%perPage > 0 {
		totalPages++
	}

	return &UserPagination{
		Data:       users,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (svc *userService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := svc.repo.GetByUsername(ctx, email)
	if err != nil || user == nil {
		return "", errors.New("Invalid credentials")
	}

	if !pkg.CheckPassword(user.Password, password) {
		return "", errors.New("Invalid credentials")
	}

	token, err := pkg.CreateToken(user.Email, user.Status, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (svc *userService) Create(ctx context.Context, user models.User) (*models.User, error) {
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

func (svc *userService) GetDetail(ctx context.Context, email string) (*models.User, error) {
	user, err := svc.repo.GetByUsername(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("Pasien not found")
	}

	return user, nil
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

func (svc *userService) Update(ctx context.Context, email, password, role, status string) (*models.User, error) {
	existing, err := svc.repo.GetByUsername(ctx, email)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return nil, errors.New("User not found")
	}

	updatedUser, err := svc.repo.UpdateData(ctx, models.User{
		Email:    email,
		Password: password,
		Role:     role,
		Status:   status,
	})
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (svc *userService) DeleteUser(ctx context.Context, email string) (*models.User, error) {
	existing, err := svc.repo.GetByUsername(ctx, email)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return nil, errors.New("User not found")
	}

	if existing.Status == "tidak aktif" {
		return nil, errors.New("User is already inactive")
	}

	updatedUser, err := svc.repo.UpdateStatus(ctx, models.User{
		Email:  email,
		Status: "tidak aktif",
	})
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
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
