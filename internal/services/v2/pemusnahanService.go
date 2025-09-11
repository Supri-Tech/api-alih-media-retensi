package services

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/repositories/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/pkg"
	"github.com/xuri/excelize/v2"
)

type PemusnahanService interface {
	GetAll(ctx context.Context, page, perPage int) (*PemusnahanPagination, error)
	Search(ctx context.Context, filter PemusnahanFilter) ([]*models.PemusnahanJoin, error)
	GetByID(ctx context.Context, id int) (*models.PemusnahanJoin, error)
	Create(ctx context.Context, pemusnahan models.Pemusnahan) (*models.Pemusnahan, error)
	Update(ctx context.Context, pemusnahan models.Pemusnahan) (*models.Pemusnahan, error)
	Delete(ctx context.Context, id int) error
	Export(ctx context.Context) ([]byte, error)
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
	Statistik  PemusnahanStatistik      `json:"statistik"`
}

type PemusnahanStatistik struct {
	TotalDokumen int `json:"total_dokumen"`
	TotalSudah   int `json:"total_sudah"`
	TotalBelum   int `json:"total_belum"`
}

type PemusnahanFilter struct {
	NoRM       string
	NamaPasien string
	Limit      int
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

	total, sudah, belum, err := svc.repo.GetStatistikPemusnahan(ctx)
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
		Statistik: PemusnahanStatistik{
			TotalDokumen: total,
			TotalSudah:   sudah,
			TotalBelum:   belum,
		},
	}, nil
}

func (svc *pemusnahanService) Search(ctx context.Context, filter PemusnahanFilter) ([]*models.PemusnahanJoin, error) {
	filterMap := make(map[string]interface{})
	if filter.NoRM != "" {
		filterMap["NoRM"] = filter.NoRM
	}
	if filter.NamaPasien != "" {
		filterMap["NamaPasien"] = filter.NamaPasien
	}
	if filter.Limit > 0 {
		filterMap["Limit"] = filter.Limit
	}

	pemusnahan, err := svc.repo.FindPemusnahan(ctx, filterMap)
	if err != nil {
		return nil, err
	}

	if len(pemusnahan) == 0 {
		return nil, errors.New("No pemusnahan found")
	}

	return pemusnahan, nil
}

// func (svc *pemusnahanService) GetAll(ctx context.Context, page, perPage int) (*PemusnahanPagination, error) {
// 	if page < 1 {
// 		page = 1
// 	}
// 	if perPage < 1 || perPage > 100 {
// 		perPage = 10
// 	}

// 	offset := (page - 1) * perPage

// 	pemusnahan, err := svc.repo.GetAllPemusnahan(ctx, perPage, offset)
// 	if err != nil {
// 		return nil, err
// 	}

// 	total, err := svc.repo.GetTotalPemusnahan(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	totalPages := total / perPage
// 	if total%perPage > 0 {
// 		totalPages++
// 	}

// 	return &PemusnahanPagination{
// 		Data:       pemusnahan,
// 		Total:      total,
// 		Page:       page,
// 		PerPage:    perPage,
// 		TotalPages: totalPages,
// 	}, nil
// }

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

func (svc *pemusnahanService) Export(ctx context.Context) ([]byte, error) {
	data, err := svc.repo.GetAllPemusnahanForExport(ctx)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		log.Println("[Export] Tidak ada data pemusnahan")
	}

	f, err := excelize.OpenFile("./templates/pemusnahan-template.xlsx")
	if err != nil {
		return nil, fmt.Errorf("Failed to open template: %v", err)
	}
	defer f.Close()

	sheetName := "Worksheet"
	startRow := 6
	endRow := 1000

	for row := startRow; row <= endRow; row++ {
		for col := 1; col <= 16; col++ { // 16 columns for alih media
			cell, _ := excelize.CoordinatesToCellName(col, row)
			f.SetCellValue(sheetName, cell, "")
		}
	}

	for i, row := range data {
		rowNum := i + startRow

		// Use pkg.GetCell for better maintainability
		f.SetCellValue(sheetName, pkg.GetCell(1, rowNum), row.ID)

		if row.TglLaporan != nil {
			f.SetCellValue(sheetName, pkg.GetCell(2, rowNum), row.TglLaporan.Format("2006-01-02"))
		} else {
			f.SetCellValue(sheetName, pkg.GetCell(2, rowNum), "-")
		}

		f.SetCellValue(sheetName, pkg.GetCell(3, rowNum), row.Status)
		f.SetCellValue(sheetName, pkg.GetCell(4, rowNum), row.JenisKunjungan)
		f.SetCellValue(sheetName, pkg.GetCell(5, rowNum), row.NoRM)
		f.SetCellValue(sheetName, pkg.GetCell(6, rowNum), row.NamaPasien)
		f.SetCellValue(sheetName, pkg.GetCell(7, rowNum), row.JenisKelamin)
		f.SetCellValue(sheetName, pkg.GetCell(8, rowNum), row.TglLahir.Format("2006-01-02"))
		f.SetCellValue(sheetName, pkg.GetCell(9, rowNum), row.Alamat)
		f.SetCellValue(sheetName, pkg.GetCell(10, rowNum), row.StatusPasien)
		f.SetCellValue(sheetName, pkg.GetCell(11, rowNum), row.JenisKasus)
		f.SetCellValue(sheetName, pkg.GetCell(12, rowNum), row.MasaAktifRi)
		f.SetCellValue(sheetName, pkg.GetCell(13, rowNum), row.MasaInaktifRi)
		f.SetCellValue(sheetName, pkg.GetCell(14, rowNum), row.MasaAktifRj)
		f.SetCellValue(sheetName, pkg.GetCell(15, rowNum), row.MasaInaktifRj)
		f.SetCellValue(sheetName, pkg.GetCell(16, rowNum), row.InfoLain)

		// Optional: Apply alternating row colors
		// if i%2 == 0 {
		// 	pkg.SetRowStyle(f, sheetName, rowNum, 16, "E2EFDA")
		// } else {
		// 	pkg.SetRowStyle(f, sheetName, rowNum, 16, "FFFFFF")
		// }
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
