package services

import (
	"context"
	"errors"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/repositories/v2"
)

type RetensiService interface {
	GetAll(ctx context.Context, page, perPage int) (*RetensiPagination, error)
	GetByID(ctx context.Context, id int) (*models.RetensiJoin, error)
	Create(ctx context.Context, retensi models.Retensi) (*models.Retensi, error)
	Update(ctx context.Context, retensi models.Retensi) (*models.Retensi, error)
	Delete(ctx context.Context, id int) error
}

type retensiService struct {
	repo repositories.RetensiRepository
}

func NewServiceRetensi(repo repositories.RetensiRepository) RetensiService {
	return &retensiService{repo: repo}
}

type RetensiPagination struct {
	Data       []*models.RetensiJoin `json:"data"`
	Total      int                   `json:"total"`
	Page       int                   `json:"page"`
	PerPage    int                   `json:"per_page"`
	TotalPages int                   `json:"total_pages"`
}

func (svc *retensiService) GetAll(ctx context.Context, page, perPage int) (*RetensiPagination, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	retensi, err := svc.repo.GetAllRetensi(ctx, perPage, offset)
	if err != nil {
		return nil, err
	}

	total, err := svc.repo.GetTotalRetensi(ctx)
	if err != nil {
		return nil, err
	}

	totalPages := total / perPage
	if total%perPage > 0 {
		totalPages++
	}

	return &RetensiPagination{
		Data:       retensi,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (svc *retensiService) GetByID(ctx context.Context, id int) (*models.RetensiJoin, error) {
	if id <= 0 {
		return nil, errors.New("ID must can't be negative")
	}

	retensi, err := svc.repo.GetRetensiByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if retensi == nil {
		return nil, errors.New("Retensi not found")
	}

	return retensi, nil
}

func (svc *retensiService) Create(ctx context.Context, retensi models.Retensi) (*models.Retensi, error) {
	newRetensi, err := svc.repo.CreateRetensi(ctx, &retensi)
	if err != nil {
		return nil, err
	}

	return newRetensi, nil
}

func (svc *retensiService) Update(ctx context.Context, retensi models.Retensi) (*models.Retensi, error) {
	existing, err := svc.repo.GetRetensiByID(ctx, retensi.ID)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return nil, errors.New("Alih Media not found")
	}

	newRetensi, err := svc.repo.UpdateRetensi(ctx, retensi)

	if err != nil {
		return nil, err
	}

	return newRetensi, nil
}

func (svc *retensiService) Delete(ctx context.Context, id int) error {
	existing, err := svc.repo.GetRetensiByID(ctx, id)
	if err != nil {
		return err
	}

	if existing == nil {
		return errors.New("Alih Media not found")
	}

	return svc.repo.DeleteRetensi(ctx, id)
}
