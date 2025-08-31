package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/repositories/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/pkg"
	"github.com/xuri/excelize/v2"
)

type KunjunganService interface {
	GetAll(ctx context.Context, page, perPage int) (*KunjunganPagination, error)
	GetByID(ctx context.Context, id int) (*models.KunjunganJoin, error)
	Create(ctx context.Context, kunjungan models.Kunjungan) (*models.Kunjungan, error)
	Update(ctx context.Context, kunjungan models.Kunjungan) (*models.Kunjungan, error)
	Delete(ctx context.Context, id int) error
	Import(ctx context.Context, filePath string) error
}

type kunjunganService struct {
	repo       repositories.KunjunganRepository
	pasienRepo repositories.PasienRepository
	kasusRepo  repositories.KasusRepository
}

func NewServiceKunjungan(
	repo repositories.KunjunganRepository,
	pasienRepo repositories.PasienRepository,
	kasusRepo repositories.KasusRepository,
) KunjunganService {
	return &kunjunganService{
		repo:       repo,
		pasienRepo: pasienRepo,
		kasusRepo:  kasusRepo,
	}
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
		return nil, nil
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

func (svc *kunjunganService) Import(ctx context.Context, filePath string) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("Failed to open Excel file: %v", err)
	}
	defer f.Close()

	rows, err := f.GetRows("Worksheet")
	if err != nil {
		return fmt.Errorf("Failed to get rows: %v", err)
	}

	for i := 4; i < len(rows); i++ {
		if len(rows[i]) < 3 {
			continue
		}

		tglMasuk := pkg.ParseDate(rows[i][7])

		noRM := rows[i][0]
		pasien, err := svc.pasienRepo.GetPasienByNoRM(ctx, noRM)
		if err != nil || pasien == nil {
			continue
		}

		jenisKasus := rows[i][8]
		kasusList, err := svc.kasusRepo.FindKasus(ctx, map[string]string{"JenisKasus": jenisKasus})

		if err != nil || len(kasusList) == 0 {
			continue
		}

		kunjungan := models.Kunjungan{
			IDPasien:       pasien.ID,
			IDKasus:        kasusList[0].ID,
			TanggalMasuk:   tglMasuk,
			JenisKunjungan: rows[i][9],
		}

		_, err = svc.repo.CreateKunjungan(ctx, &kunjungan)
		if err != nil {
			continue
		}
	}

	return nil
}
