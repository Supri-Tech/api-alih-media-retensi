package services

import (
	"context"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/repositories/v2"
)

type KunjunganService interface {
	GetAll(ctx context.Context, page, perPage int) (*KunjunganPagination, error)
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
