package services

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/repositories/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/pkg"
	"github.com/xuri/excelize/v2"
)

type KasusService interface {
	GetAll(ctx context.Context, page, perPage int) (*KasusPagination, error)
	GetByID(ctx context.Context, id int) (*models.Kasus, error)
	Search(ctx context.Context, filter KasusFilter) ([]*models.Kasus, error)
	Create(ctx context.Context, kasus models.Kasus) (*models.Kasus, error)
	Update(ctx context.Context, kasus models.Kasus) (*models.Kasus, error)
	Delete(ctx context.Context, id int) error
	Import(ctx context.Context, filepath string) error
	Export(ctx context.Context, filter KasusFilter) ([]byte, error)
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
	Statistik  KasusStatistik  `json:"statistik"`
}

type KasusStatistik struct {
	RataMasaAktifRI   float64 `json:"rata_masa_aktif_ri"`
	RataMasaInaktifRI float64 `json:"rata_masa_inaktif_ri"`
	RataMasaAktifRJ   float64 `json:"rata_masa_aktif_rj"`
	RataMasaInaktifRJ float64 `json:"rata_masa_inaktif_rj"`
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

	rataAktifRI, rataInaktifRI, rataAktifRJ, rataInaktifRJ, err := svc.repo.GetStatistikKasus(ctx)
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
		Statistik: KasusStatistik{
			RataMasaAktifRI:   rataAktifRI,
			RataMasaInaktifRI: rataInaktifRI,
			RataMasaAktifRJ:   rataAktifRJ,
			RataMasaInaktifRJ: rataInaktifRJ,
		},
	}, nil
}

// func (svc *kasusService) GetAll(ctx context.Context, page, perPage int) (*KasusPagination, error) {
// 	if page < 1 {
// 		page = 1
// 	}
// 	if perPage < 1 || perPage > 100 {
// 		perPage = 10
// 	}

// 	offset := (page - 1) * perPage

// 	kasus, err := svc.repo.GetAllKasus(ctx, perPage, offset)
// 	if err != nil {
// 		return nil, err
// 	}

// 	total, err := svc.repo.GetTotalKasus(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	totalPages := total / perPage
// 	if total%perPage > 0 {
// 		totalPages++
// 	}

// 	return &KasusPagination{
// 		Data:       kasus,
// 		Total:      total,
// 		Page:       page,
// 		PerPage:    perPage,
// 		TotalPages: totalPages,
// 	}, nil
// }

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

func (svc *kasusService) Import(ctx context.Context, filepath string) error {
	f, err := excelize.OpenFile(filepath)
	if err != nil {
		return fmt.Errorf("Failed to open Excel file: %v", err)
	}
	defer f.Close()

	rows, err := f.GetRows("Worksheet")
	if err != nil {
		return fmt.Errorf("Failed to get rows: %v", err)
	}

	for i := 5; i < len(rows); i++ {
		if len(rows[i]) < 6 {
			continue
		}

		masaAktifRi, _ := strconv.Atoi(rows[i][1])
		masaInaktifRi, _ := strconv.Atoi(rows[i][2])
		masaAktifRj, _ := strconv.Atoi(rows[i][3])
		masaInaktifRj, _ := strconv.Atoi(rows[i][4])

		kasus := models.Kasus{
			JenisKasus:    rows[i][0],
			MasaAktifRI:   masaAktifRi,
			MasaInaktifRI: masaInaktifRi,
			MasaAktifRJ:   masaAktifRj,
			MasaInaktifRJ: masaInaktifRj,
		}

		existing, err := svc.repo.FindKasus(ctx, map[string]string{"JenisKasus": kasus.JenisKasus})
		if err != nil || existing == nil {
			_, err := svc.repo.CreateKasus(ctx, kasus)
			if err != nil {
				log.Printf("Failed to create pasien %s: %v", kasus.JenisKasus, err)
			}
		} else {
			kasus.ID = existing[0].ID
			_, err := svc.repo.UpdateKasus(ctx, kasus)
			if err != nil {
				log.Printf("Failed to update pasien %s: %v", kasus.JenisKasus, err)
			}
		}
	}

	return nil
}

func (svc *kasusService) Export(ctx context.Context, filter KasusFilter) ([]byte, error) {
	filterMap := make(map[string]string)
	if filter.JenisKasus != "" {
		filterMap["JenisKasus"] = filter.JenisKasus
	}

	kasusList, err := svc.repo.FindKasus(ctx, filterMap)
	if err != nil {
		return nil, err
	}

	f, err := excelize.OpenFile("./templates/kasus-template.xlsx")
	if err != nil {
		return nil, fmt.Errorf("Failed to open template: %v", err)
	}
	defer f.Close()

	sheetName := "Worksheet"
	startRow := 6
	endRow := 1000

	for row := startRow; row <= endRow; row++ {
		for col := 1; col <= 8; col++ {
			cell, _ := excelize.CoordinatesToCellName(col, row)
			f.SetCellValue(sheetName, cell, "")
		}
	}

	// headers := []string{"Jenis Kasus", "Masa Aktif RI", "Masa Inaktif RI", "Masa Aktif RJ", "Masa Inaktif RJ", "Info Lain"}
	// for col, header := range headers {
	// 	cell, _ := excelize.CoordinatesToCellName(col+1, 1)
	// 	f.SetCellValue(sheetName, cell, header)
	// 	f.SetCellStyle(sheetName, cell, cell, pkg.GetHeaderStyle(f))
	// }

	for row, kasus := range kasusList {
		rowNum := row + startRow

		f.SetCellValue(sheetName, pkg.GetCell(1, rowNum), kasus.JenisKasus)
		f.SetCellValue(sheetName, pkg.GetCell(2, rowNum), kasus.MasaAktifRI)
		f.SetCellValue(sheetName, pkg.GetCell(3, rowNum), kasus.MasaInaktifRI)
		f.SetCellValue(sheetName, pkg.GetCell(4, rowNum), kasus.MasaAktifRJ)
		f.SetCellValue(sheetName, pkg.GetCell(5, rowNum), kasus.MasaInaktifRJ)
		f.SetCellValue(sheetName, pkg.GetCell(6, rowNum), kasus.InfoLain)

		// if row%2 == 0 {
		// 	pkg.SetRowStyle(f, sheetName, rowNum, 8, "E2EFDA")
		// } else {
		// 	pkg.SetRowStyle(f, sheetName, rowNum, 8, "FFFFFF")
		// }
	}

	for col := 1; col <= 8; col++ {
		f.SetColWidth(sheetName, pkg.GetColumnName(col), pkg.GetColumnName(col), 20)
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
