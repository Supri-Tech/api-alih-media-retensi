package services

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/repositories/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/pkg"
	"github.com/xuri/excelize/v2"
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
	Import(ctx context.Context, filePath string) error
	Export(ctx context.Context, filter PasienFilter) ([]byte, error)
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
	Statistik  PasienStatistik  `json:"statistik"`
	Total      int              `json:"total"`
	Page       int              `json:"page"`
	PerPage    int              `json:"per_page"`
	TotalPages int              `json:"total_pages"`
}

type PasienStatistik struct {
	Total           int `json:"total"`
	TotalAktif      int `json:"total_aktif"`
	TotalTidakAktif int `json:"total_tidak_aktif"`
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

	total, aktif, tidak_aktif, err := svc.repo.GetStatistikPasien(ctx)
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
		Statistik: PasienStatistik{
			Total:           total,
			TotalAktif:      aktif,
			TotalTidakAktif: tidak_aktif,
		},
	}, nil
}

// func (svc *pasienService) GetAll(ctx context.Context, page, perPage int) (*PasienPagination, error) {
// 	if page < 1 {
// 		page = 1
// 	}
// 	if perPage < 1 || perPage > 100 {
// 		perPage = 10
// 	}

// 	offset := (page - 1) * perPage

// 	pasien, err := svc.repo.GetAllPasien(ctx, perPage, offset)
// 	if err != nil {
// 		return nil, err
// 	}

// 	total, err := svc.repo.GetTotalPasien(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	totalPages := total / perPage
// 	if total%perPage > 0 {
// 		totalPages++
// 	}

// 	return &PasienPagination{
// 		Data:       pasien,
// 		Total:      total,
// 		Page:       page,
// 		PerPage:    perPage,
// 		TotalPages: totalPages,
// 	}, nil
// }

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

func (svc *pasienService) Import(ctx context.Context, filepath string) error {
	f, err := excelize.OpenFile(filepath)
	if err != nil {
		return fmt.Errorf("Failed to open Excel file: %v", err)
	}
	defer f.Close()

	rows, err := f.GetRows("Worksheet")
	if err != nil {
		return fmt.Errorf("Failed to get rows: %v", err)
	}

	for i := 4; i < len(rows); i++ {
		if len(rows[i]) < 6 {
			continue
		}

		pasien := models.Pasien{
			NoRM:         rows[i][0],
			NamaPasien:   rows[i][1],
			JenisKelamin: rows[i][2],
			TanggalLahir: pkg.ParseDate(rows[i][3]),
			NIK:          rows[i][4],
			Alamat:       rows[i][5],
			Status:       rows[i][6],
			CreatedAt:    time.Now(),
		}

		existing, err := svc.repo.GetPasienByNoRM(ctx, pasien.NoRM)
		if err != nil || existing == nil {
			_, err := svc.repo.CreatePasien(ctx, pasien)
			if err != nil {
				log.Printf("Failed to create pasien %s: %v", pasien.NoRM, err)
			}
		} else {
			pasien.ID = existing.ID
			_, err := svc.repo.UpdatePasien(ctx, pasien)
			if err != nil {
				log.Printf("Failed to update pasien %s: %v", pasien.NoRM, err)
			}
		}
	}

	return nil
}

func (svc *pasienService) Export(ctx context.Context, filter PasienFilter) ([]byte, error) {
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

	pasiens, err := svc.repo.FindPasien(ctx, filterMap)
	if err != nil {
		return nil, err
	}

	// f := excelize.NewFile()
	f, err := excelize.OpenFile("./templates/pasien-template.xlsx")
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

	// index, err := f.NewSheet(sheetName)
	// if err != nil {
	// 	return nil, err
	// }
	// f.SetActiveSheet(index)

	// headers := []string{"No RM", "Nama Pasien", "Jenis Kelamin", "Tanggal Lahir", "NIK", "Alamat", "Status", "Tanggal Dibuat"}
	// for col, header := range headers {
	// 	cell, _ := excelize.CoordinatesToCellName(col+1, 1)
	// 	f.SetCellValue(sheetName, cell, header)
	// 	f.SetCellStyle(sheetName, cell, cell, pkg.GetHeaderStyle(f))
	// }

	for row, pasien := range pasiens {
		rowNum := row + startRow

		f.SetCellValue(sheetName, pkg.GetCell(1, rowNum), pasien.NoRM)
		f.SetCellValue(sheetName, pkg.GetCell(2, rowNum), pasien.NamaPasien)
		f.SetCellValue(sheetName, pkg.GetCell(3, rowNum), pasien.JenisKelamin)
		f.SetCellValue(sheetName, pkg.GetCell(4, rowNum), pasien.TanggalLahir.Format("2006-01-02"))
		f.SetCellValue(sheetName, pkg.GetCell(5, rowNum), pasien.NIK)
		f.SetCellValue(sheetName, pkg.GetCell(6, rowNum), pasien.Alamat)
		f.SetCellValue(sheetName, pkg.GetCell(7, rowNum), pasien.Status)
		f.SetCellValue(sheetName, pkg.GetCell(8, rowNum), pasien.CreatedAt.Format("2006-01-02 15:04:05"))

		// if row%2 == 0 {
		// 	pkg.SetRowStyle(f, sheetName, rowNum, 8, "E2EFDA")
		// } else {
		// 	pkg.SetRowStyle(f, sheetName, rowNum, 8, "FFFFFF")
		// }
	}

	// for col := 1; col <= 8; col++ {
	// 	f.SetColWidth(sheetName, pkg.GetColumnName(col), pkg.GetColumnName(col), 20)
	// }

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
