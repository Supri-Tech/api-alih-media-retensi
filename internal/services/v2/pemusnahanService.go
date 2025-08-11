package services

import (
	"context"
	"errors"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/repositories/v2"
)

type PemusnahanService interface {
	GetAll(ctx context.Context, page, perPage int) (*PemusnahanPagination, error)
	GetByID(ctx context.Context, id int) (*models.PemusnahanJoin, error)
	Create(ctx context.Context, pemusnahan models.Pemusnahan) (*models.Pemusnahan, error)
	Update(ctx context.Context, pemusnahan models.Pemusnahan) (*models.Pemusnahan, error)
	Delete(ctx context.Context, id int) error
}

type pemusnahanService struct {
	repo repositories.PemusnahanRepository
}

func NewServicePemusnahan(repo repositories.PemusnahanRepository) PemusnahanService {
	return &pemusnahanService{repo: repo}
}

type PemusnahanPagination struct {
	Data       []*models.PemusnahanJoin `json:"data"`
	Total      int                      `json:"total"`
	Page       int                      `json:"page"`
	PerPage    int                      `json:"per_page"`
	TotalPages int                      `json:"total_pages"`
}

func (svc *pemusnahanService) GetAll(ctx context.Context, page, perPage int) (*PemusnahanPagination, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	pemusnahan, err := svc.repo.GetAllPemusnahan(ctx, perPage, offset)
	if err != nil {
		return nil, err
	}

	total, err := svc.repo.GetTotalPemusnahan(ctx)
	if err != nil {
		return nil, err
	}

	totalPages := total / perPage
	if total%perPage > 0 {
		totalPages++
	}

	return &PemusnahanPagination{
		Data:       pemusnahan,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (svc *pemusnahanService) GetByID(ctx context.Context, id int) (*models.PemusnahanJoin, error) {
	if id <= 0 {
		return nil, errors.New("ID must can't be negative")
	}

	pemusnahan, err := svc.repo.GetPemusnahanByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if pemusnahan == nil {
		return nil, errors.New("Pemusnahan not found")
	}

	return pemusnahan, nil
}

func (svc *pemusnahanService) Create(ctx context.Context, pemusnahan models.Pemusnahan) (*models.Pemusnahan, error) {
	newPemusnahan, err := svc.repo.CreatePemusnahan(ctx, &pemusnahan)
	if err != nil {
		return nil, err
	}

	return newPemusnahan, nil
}

func (svc *pemusnahanService) Update(ctx context.Context, pemusnahan models.Pemusnahan) (*models.Pemusnahan, error) {
	existing, err := svc.repo.GetPemusnahanByID(ctx, pemusnahan.ID)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return nil, errors.New("Alih Media not found")
	}

	newPemusnahan, err := svc.repo.UpdatePemusnahan(ctx, pemusnahan)

	if err != nil {
		return nil, err
	}

	return newPemusnahan, nil
}

func (svc *pemusnahanService) Delete(ctx context.Context, id int) error {
	existing, err := svc.repo.GetPemusnahanByID(ctx, id)
	if err != nil {
		return err
	}

	if existing == nil {
		return errors.New("Alih Media not found")
	}

	return svc.repo.DeletePemusnahan(ctx, id)
}
