package services

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/repositories/v2"
	"github.com/xuri/excelize/v2"
)

type RetensiService interface {
	GetAll(ctx context.Context, page, perPage int) (*RetensiPagination, error)
	GetByID(ctx context.Context, id int) (*models.RetensiJoin, error)
	Create(ctx context.Context, retensi models.Retensi) (*models.Retensi, error)
	Update(ctx context.Context, retensi models.Retensi) (*models.Retensi, error)
	Delete(ctx context.Context, id int) error
	Export(ctx context.Context) ([]byte, error)
}

type retensiService struct {
	repo repositories.RetensiRepository
}

func NewServiceRetensi(repo repositories.RetensiRepository) RetensiService {
	return &retensiService{repo: repo}
}

type RetensiPagination struct {
	Data       []*models.RetensiJoin `json:"data"`
	Total      int                   `json:"total"`
	Page       int                   `json:"page"`
	PerPage    int                   `json:"per_page"`
	TotalPages int                   `json:"total_pages"`
	Statistik  RetensiStatistik      `json:"statistik"`
}

type RetensiStatistik struct {
	TotalDokumen int `json:"total_dokumen"`
	TotalSudah   int `json:"total_sudah"`
	TotalBelum   int `json:"total_belum"`
}

func (svc *retensiService) GetAll(ctx context.Context, page, perPage int) (*RetensiPagination, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	retensi, err := svc.repo.GetAllRetensi(ctx, perPage, offset)
	if err != nil {
		return nil, err
	}

	total, sudah, belum, err := svc.repo.GetStatistikRetensi(ctx)
	if err != nil {
		return nil, err
	}

	totalPages := total / perPage
	if total%perPage > 0 {
		totalPages++
	}

	return &RetensiPagination{
		Data:       retensi,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
		Statistik: RetensiStatistik{
			TotalDokumen: total,
			TotalSudah:   sudah,
			TotalBelum:   belum,
		},
	}, nil
}

// func (svc *retensiService) GetAll(ctx context.Context, page, perPage int) (*RetensiPagination, error) {
// 	if page < 1 {
// 		page = 1
// 	}
// 	if perPage < 1 || perPage > 100 {
// 		perPage = 10
// 	}

// 	offset := (page - 1) * perPage

// 	retensi, err := svc.repo.GetAllRetensi(ctx, perPage, offset)
// 	if err != nil {
// 		return nil, err
// 	}

// 	total, err := svc.repo.GetTotalRetensi(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	totalPages := total / perPage
// 	if total%perPage > 0 {
// 		totalPages++
// 	}

// 	return &RetensiPagination{
// 		Data:       retensi,
// 		Total:      total,
// 		Page:       page,
// 		PerPage:    perPage,
// 		TotalPages: totalPages,
// 	}, nil
// }

func (svc *retensiService) GetByID(ctx context.Context, id int) (*models.RetensiJoin, error) {
	if id <= 0 {
		return nil, errors.New("ID must can't be negative")
	}

	retensi, err := svc.repo.GetRetensiByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if retensi == nil {
		return nil, errors.New("Retensi not found")
	}

	return retensi, nil
}

func (svc *retensiService) Create(ctx context.Context, retensi models.Retensi) (*models.Retensi, error) {
	newRetensi, err := svc.repo.CreateRetensi(ctx, &retensi)
	if err != nil {
		return nil, err
	}

	return newRetensi, nil
}

func (svc *retensiService) Update(ctx context.Context, retensi models.Retensi) (*models.Retensi, error) {
	existing, err := svc.repo.GetRetensiByID(ctx, retensi.ID)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return nil, errors.New("Alih Media not found")
	}

	newRetensi, err := svc.repo.UpdateRetensi(ctx, retensi)

	if err != nil {
		return nil, err
	}

	return newRetensi, nil
}

func (svc *retensiService) Delete(ctx context.Context, id int) error {
	existing, err := svc.repo.GetRetensiByID(ctx, id)
	if err != nil {
		return err
	}

	if existing == nil {
		return errors.New("Alih Media not found")
	}

	return svc.repo.DeleteRetensi(ctx, id)
}

func (svc *retensiService) Export(ctx context.Context) ([]byte, error) {
	data, err := svc.repo.GetAllRetensiForExport(ctx)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		log.Println("[Export] Tidak ada data retensi")
	}

	f := excelize.NewFile()
	sheet := "Sheet1"
	f.SetSheetName(f.GetSheetName(0), sheet)

	// Header
	headers := []string{
		"ID", "Tanggal Laporan", "Status", "Jenis Kunjungan",
		"NoRM", "Nama Pasien", "Jenis Kelamin", "Tanggal Lahir",
		"Alamat", "Status Pasien", "Jenis Kasus",
		"Masa Aktif RI", "Masa Inaktif RI",
		"Masa Aktif RJ", "Masa Inaktif RJ", "Info Lain",
	}
	for i, h := range headers {
		col := string(rune('A' + i))
		f.SetCellValue(sheet, col+"1", h)
	}

	// Data
	for i, row := range data {
		r := i + 2
		f.SetCellValue(sheet, "A"+strconv.Itoa(r), row.ID)
		if row.TglLaporan != nil {
			f.SetCellValue(sheet, "B"+strconv.Itoa(r), row.TglLaporan.Format("2006-01-02"))
		} else {
			f.SetCellValue(sheet, "B"+strconv.Itoa(r), "-")
		}
		f.SetCellValue(sheet, "C"+strconv.Itoa(r), row.Status)
		f.SetCellValue(sheet, "D"+strconv.Itoa(r), row.JenisKunjungan)
		f.SetCellValue(sheet, "E"+strconv.Itoa(r), row.NoRM)
		f.SetCellValue(sheet, "F"+strconv.Itoa(r), row.NamaPasien)
		f.SetCellValue(sheet, "G"+strconv.Itoa(r), row.JenisKelamin)
		f.SetCellValue(sheet, "H"+strconv.Itoa(r), row.TglLahir.Format("2006-01-02"))
		f.SetCellValue(sheet, "I"+strconv.Itoa(r), row.Alamat)
		f.SetCellValue(sheet, "J"+strconv.Itoa(r), row.StatusPasien)
		f.SetCellValue(sheet, "K"+strconv.Itoa(r), row.JenisKasus)
		f.SetCellValue(sheet, "L"+strconv.Itoa(r), row.MasaAktifRi)
		f.SetCellValue(sheet, "M"+strconv.Itoa(r), row.MasaInaktifRi)
		f.SetCellValue(sheet, "N"+strconv.Itoa(r), row.MasaAktifRj)
		f.SetCellValue(sheet, "O"+strconv.Itoa(r), row.MasaInaktifRj)
		f.SetCellValue(sheet, "P"+strconv.Itoa(r), row.InfoLain)
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
