package services

import (
	"context"
	"errors"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/repositories/v2"
)

type KasusService interface {
	GetAll(ctx context.Context, page, perPage int) (*KasusPagination, error)
	GetByID(ctx context.Context, id int) (*models.Kasus, error)
	Search(ctx context.Context, filter KasusFilter) ([]*models.Kasus, error)
	Create(ctx context.Context, kasus models.Kasus) (*models.Kasus, error)
	Update(ctx context.Context, kasus models.Kasus) (*models.Kasus, error)
	Delete(ctx context.Context, id int) error
}

type KasusFilter struct {
	JenisKasus string
	Limit      int
}

type kasusService struct {
	repo repositories.KasusRepository
}

func NewServiceKasus(repo repositories.KasusRepository) KasusService {
	return &kasusService{repo: repo}
}

type KasusPagination struct {
	Data       []*models.Kasus `json:"data"`
	Total      int             `json:"total"`
	Page       int             `json:"page"`
	PerPage    int             `json:"per_page"`
	TotalPages int             `json:"total_pages"`
}

func (svc *kasusService) GetAll(ctx context.Context, page, perPage int) (*KasusPagination, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	kasus, err := svc.repo.GetAllKasus(ctx, perPage, offset)
	if err != nil {
		return nil, err
	}

	total, err := svc.repo.GetTotalKasus(ctx)
	if err != nil {
		return nil, err
	}

	totalPages := total / perPage
	if total%perPage > 0 {
		totalPages++
	}

	return &KasusPagination{
		Data:       kasus,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (svc *kasusService) GetByID(ctx context.Context, id int) (*models.Kasus, error) {
	if id <= 0 {
		return nil, errors.New("ID must can't be negative")
	}

	kasus, err := svc.repo.GetKasusByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if kasus == nil {
		return nil, errors.New("Kasus not found")
	}

	return kasus, nil
}

func (svc *kasusService) Search(ctx context.Context, filter KasusFilter) ([]*models.Kasus, error) {
	filterMap := make(map[string]string)
	if filter.JenisKasus != "" {
		filterMap["JenisKasus"] = filter.JenisKasus
	}

	kasus, err := svc.repo.FindKasus(ctx, filterMap)
	if err != nil {
		return nil, err
	}

	if len(kasus) == 0 {
		return nil, errors.New("No kasus found")
	}

	return kasus, nil
}

func (svc *kasusService) Create(ctx context.Context, kasus models.Kasus) (*models.Kasus, error) {
	newKasus, err := svc.repo.CreateKasus(ctx, kasus)
	if err != nil {
		return nil, err
	}

	return newKasus, nil
}

func (svc *kasusService) Update(ctx context.Context, kasus models.Kasus) (*models.Kasus, error) {
	existing, err := svc.repo.GetKasusByID(ctx, kasus.ID)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return nil, errors.New("Kasus not found")
	}

	newKasus, err := svc.repo.UpdateKasus(ctx, kasus)
	if err != nil {
		return nil, err
	}

	return newKasus, nil
}

func (svc *kasusService) Delete(ctx context.Context, id int) error {
	existing, err := svc.repo.GetKasusByID(ctx, id)
	if err != nil {
		return err
	}

	if existing == nil {
		return errors.New("Kasus not found")
	}

	return svc.repo.DeleteKasus(ctx, id)
}
