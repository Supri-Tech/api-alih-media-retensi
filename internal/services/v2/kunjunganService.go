package services

import (
	"context"
	"errors"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/repositories/v2"
)

type KunjunganService interface {
	GetAll(ctx context.Context, page, perPage int) (*KunjunganPagination, error)
	GetByID(ctx context.Context, id int) (*models.KunjunganJoin, error)
	Create(ctx context.Context, kunjungan models.Kunjungan) (*models.Kunjungan, error)
	Update(ctx context.Context, kunjungan models.Kunjungan) (*models.Kunjungan, error)
	Delete(ctx context.Context, id int) error
}

type kunjunganService struct {
	repo repositories.KunjunganRepository
}

func NewServiceKunjungan(repo repositories.KunjunganRepository) KunjunganService {
	return &kunjunganService{repo: repo}
}

type KunjunganPagination struct {
	Data       []*models.KunjunganJoin `json:"data"`
	Total      int                     `json:"total"`
	Page       int                     `json:"page"`
	PerPage    int                     `json:"per_page"`
	TotalPages int                     `json:"total_pages"`
}

func (svc *kunjunganService) GetAll(ctx context.Context, page, perPage int) (*KunjunganPagination, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	kunjungan, err := svc.repo.GetAllKunjungan(ctx, perPage, offset)
	if err != nil {
		return nil, err
	}

	total, err := svc.repo.GetTotalKunjungan(ctx)
	if err != nil {
		return nil, err
	}

	totalPages := total / perPage
	if total%perPage > 0 {
		totalPages++
	}

	return &KunjunganPagination{
		Data:       kunjungan,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (svc *kunjunganService) GetByID(ctx context.Context, id int) (*models.KunjunganJoin, error) {
	if id <= 0 {
		return nil, errors.New("ID must can't be negative")
	}

	kunjungan, err := svc.repo.GetKunjunganByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if kunjungan == nil {
		return nil, errors.New("Kunjungan not found")
	}

	return kunjungan, nil
}

func (svc *kunjunganService) Create(ctx context.Context, kunjungan models.Kunjungan) (*models.Kunjungan, error) {
	newKunjungan, err := svc.repo.CreateKunjungan(ctx, &kunjungan)
	if err != nil {
		return nil, err
	}

	return newKunjungan, nil
}

func (svc *kunjunganService) Update(ctx context.Context, kunjungan models.Kunjungan) (*models.Kunjungan, error) {
	existing, err := svc.repo.GetKunjunganByID(ctx, kunjungan.ID)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return nil, errors.New("Kunjungan not found")
	}

	newKunjungan, err := svc.repo.UpdateKunjungan(ctx, kunjungan)

	if err != nil {
		return nil, err
	}

	return newKunjungan, nil
}

func (svc *kunjunganService) Delete(ctx context.Context, id int) error {
	existing, err := svc.repo.GetKunjunganByID(ctx, id)
	if err != nil {
		return err
	}

	if existing == nil {
		return errors.New("Pasien not found")
	}

	return svc.repo.DeleteKunjungan(ctx, id)
}
