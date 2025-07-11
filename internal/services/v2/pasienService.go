package services

import (
	"context"
	"errors"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/repositories/v2"
)

type PasienService interface {
	GetAll(ctx context.Context, page, perPage int) (*PasienPagination, error)
	GetByID(ctx context.Context, id int) (*models.Pasien, error)
	GetByNIK(ctx context.Context, NIK string) (*models.Pasien, error)
	GetByNoRM(ctx context.Context, noRM string) (*models.Pasien, error)
	GetByName(ctx context.Context, name string) (*models.Pasien, error)
	Search(ctx context.Context, filter PasienFilter) ([]*models.Pasien, error)
	Create(ctx context.Context, pasien models.Pasien) (*models.Pasien, error)
	Update(ctx context.Context, pasien models.Pasien) (*models.Pasien, error)
	Delete(ctx context.Context, id int) error
}

type PasienFilter struct {
	NoRM       string
	NamaPasien string
	NIK        string
	Limit      int
}

type pasienService struct {
	repo repositories.PasienRepository
}

func NewServicePasien(repo repositories.PasienRepository) PasienService {
	return &pasienService{repo: repo}
}

type PasienPagination struct {
	Data       []*models.Pasien `json:"data"`
	Total      int              `json:"total"`
	Page       int              `json:"page"`
	PerPage    int              `json:"per_page"`
	TotalPages int              `json:"total_pages"`
}

func (svc *pasienService) GetAll(ctx context.Context, page, perPage int) (*PasienPagination, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	pasien, err := svc.repo.GetAllPasien(ctx, perPage, offset)
	if err != nil {
		return nil, err
	}

	total, err := svc.repo.GetTotalPasien(ctx)
	if err != nil {
		return nil, err
	}

	totalPages := total / perPage
	if total%perPage > 0 {
		totalPages++
	}

	return &PasienPagination{
		Data:       pasien,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (svc *pasienService) GetByID(ctx context.Context, id int) (*models.Pasien, error) {
	if id <= 0 {
		return nil, errors.New("ID must can't be negative")
	}

	pasien, err := svc.repo.GetPasienByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if pasien == nil {
		return nil, errors.New("Pasien not found")
	}

	return pasien, nil
}

func (svc *pasienService) Search(ctx context.Context, filter PasienFilter) ([]*models.Pasien, error) {
	filterMap := make(map[string]string)
	if filter.NoRM != "" {
		filterMap["NoRM"] = filter.NoRM
	}
	if filter.NamaPasien != "" {
		filterMap["NamaPasien"] = filter.NamaPasien
	}
	if filter.NIK != "" {
		filterMap["NIK"] = filter.NIK
	}

	pasien, err := svc.repo.FindPasien(ctx, filterMap)
	if err != nil {
		return nil, err
	}

	if len(pasien) == 0 {
		return nil, errors.New("No pasien found")
	}

	return pasien, nil
}

func (svc *pasienService) GetByNIK(ctx context.Context, NIK string) (*models.Pasien, error) {
	if NIK == "" {
		return nil, errors.New("ID can't be empty")
	}

	pasien, err := svc.repo.GetPasienByNIK(ctx, NIK)
	if err != nil {
		return nil, err
	}

	if pasien == nil {
		return nil, errors.New("Pasien not found")
	}

	return pasien, nil
}

func (svc *pasienService) GetByNoRM(ctx context.Context, noRM string) (*models.Pasien, error) {
	if noRM == "" {
		return nil, errors.New("ID can't be empty")
	}

	pasien, err := svc.repo.GetPasienByNIK(ctx, noRM)
	if err != nil {
		return nil, err
	}

	if pasien == nil {
		return nil, errors.New("Pasien not found")
	}

	return pasien, nil
}

func (svc *pasienService) GetByName(ctx context.Context, name string) (*models.Pasien, error) {
	if name == "" {
		return nil, errors.New("ID can't be empty")
	}

	pasien, err := svc.repo.GetPasienByNIK(ctx, name)
	if err != nil {
		return nil, err
	}

	if pasien == nil {
		return nil, errors.New("Pasien not found")
	}

	return pasien, nil
}

func (svc *pasienService) Create(ctx context.Context, pasien models.Pasien) (*models.Pasien, error) {
	newPasien, err := svc.repo.CreatePasien(ctx, pasien)
	if err != nil {
		return nil, err
	}

	return newPasien, nil
}

func (svc *pasienService) Update(ctx context.Context, pasien models.Pasien) (*models.Pasien, error) {
	existing, err := svc.repo.GetPasienByID(ctx, pasien.ID)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return nil, errors.New("Pasien not found")
	}

	newPasien, err := svc.repo.UpdatePasien(ctx, pasien)

	if err != nil {
		return nil, err
	}
	return newPasien, nil
}

func (svc *pasienService) Delete(ctx context.Context, id int) error {
	existing, err := svc.repo.GetPasienByID(ctx, id)
	if err != nil {
		return err
	}

	if existing == nil {
		return errors.New("Pasien not found")
	}

	return svc.repo.DeletePasien(ctx, id)
}
