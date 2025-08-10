package services

import (
	"context"
	"errors"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/repositories/v2"
)

type AlihMediaService interface {
	GetAll(ctx context.Context, page, perPage int) (*AlihMediaPagination, error)
	GetByID(ctx context.Context, id int) (*models.AlihMediaJoin, error)
	Create(ctx context.Context, alihMedia models.AlihMedia) (*models.AlihMedia, error)
	Update(ctx context.Context, alihMedia models.AlihMedia) (*models.AlihMedia, error)
	Delete(ctx context.Context, id int) error
}

type alihMediaService struct {
	repo repositories.AlihMediaRepository
}

func NewServiceAlihMedia(repo repositories.AlihMediaRepository) AlihMediaService {
	return &alihMediaService{repo: repo}
}

type AlihMediaPagination struct {
	Data       []*models.AlihMediaJoin `json:"data"`
	Total      int                     `json:"total"`
	Page       int                     `json:"page"`
	PerPage    int                     `json:"per_page"`
	TotalPages int                     `json:"total_pages"`
}

func (svc *alihMediaService) GetAll(ctx context.Context, page, perPage int) (*AlihMediaPagination, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	alihMedia, err := svc.repo.GetAllAlihMedia(ctx, perPage, offset)
	if err != nil {
		return nil, err
	}

	total, err := svc.repo.GetTotalAlihMedia(ctx)
	if err != nil {
		return nil, err
	}

	totalPages := total / perPage
	if total%perPage > 0 {
		totalPages++
	}

	return &AlihMediaPagination{
		Data:       alihMedia,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (svc *alihMediaService) GetByID(ctx context.Context, id int) (*models.AlihMediaJoin, error) {
	if id <= 0 {
		return nil, errors.New("ID must can't be negative")
	}

	alihMedia, err := svc.repo.GetAlihMediaByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if alihMedia == nil {
		return nil, errors.New("Kunjungan not found")
	}

	return alihMedia, nil
}

func (svc *alihMediaService) Create(ctx context.Context, alihMedia models.AlihMedia) (*models.AlihMedia, error) {
	newAlihMedia, err := svc.repo.CreateAlihMedia(ctx, &alihMedia)
	if err != nil {
		return nil, err
	}

	return newAlihMedia, nil
}

func (svc *alihMediaService) Update(ctx context.Context, alihMedia models.AlihMedia) (*models.AlihMedia, error) {
	existing, err := svc.repo.GetAlihMediaByID(ctx, alihMedia.ID)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return nil, errors.New("Alih Media not found")
	}

	newAlihMedia, err := svc.repo.UpdateAlihMedia(ctx, alihMedia)

	if err != nil {
		return nil, err
	}

	return newAlihMedia, nil
}

func (svc *alihMediaService) Delete(ctx context.Context, id int) error {
	existing, err := svc.repo.GetAlihMediaByID(ctx, id)
	if err != nil {
		return err
	}

	if existing == nil {
		return errors.New("Alih Media not found")
	}

	return svc.repo.DeleteAlihMedia(ctx, id)
}
