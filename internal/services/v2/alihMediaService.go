package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/repositories/v2"
	"github.com/xuri/excelize/v2"
)

type AlihMediaService interface {
	GetAll(ctx context.Context, page, perPage int) (*AlihMediaPagination, error)
	GetByID(ctx context.Context, id int) (*models.AlihMediaJoin, error)
	Create(ctx context.Context, alihMedia models.AlihMedia) (*models.AlihMedia, error)
	Update(ctx context.Context, alihMedia models.AlihMedia) (*models.AlihMedia, error)
	Delete(ctx context.Context, id int) error
	CreateAndCheckAlihMedia(ctx context.Context, kunjunganID int) error
	CheckAllExpiredKunjungan(ctx context.Context) error
	Export(ctx context.Context) ([]byte, error)
}

type alihMediaService struct {
	repo          repositories.AlihMediaRepository
	kunjunganRepo repositories.KunjunganRepository
	kasusRepo     repositories.KasusRepository
}

func NewServiceAlihMedia(
	repo repositories.AlihMediaRepository,
	kunjunganRepo repositories.KunjunganRepository,
	kasusRepo repositories.KasusRepository,
) AlihMediaService {
	return &alihMediaService{
		repo:          repo,
		kunjunganRepo: kunjunganRepo,
		kasusRepo:     kasusRepo,
	}
}

type AlihMediaPagination struct {
	Data       []*models.AlihMediaJoin `json:"data"`
	Total      int                     `json:"total"`
	Page       int                     `json:"page"`
	PerPage    int                     `json:"per_page"`
	TotalPages int                     `json:"total_pages"`
	Statistik  AlihMediaStatistik      `json:"statistik"`
}

type AlihMediaStatistik struct {
	TotalDokumen int `json:"total_dokumen"`
	TotalSudah   int `json:"total_sudah"`
	TotalBelum   int `json:"total_belum"`
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

	total, sudah, belum, err := svc.repo.GetStatistikAlihMedia(ctx)
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
		Statistik: AlihMediaStatistik{
			TotalDokumen: total,
			TotalSudah:   sudah,
			TotalBelum:   belum,
		},
	}, nil
}

// func (svc *alihMediaService) GetAll(ctx context.Context, page, perPage int) (*AlihMediaPagination, error) {
// 	if page < 1 {
// 		page = 1
// 	}
// 	if perPage < 1 || perPage > 100 {
// 		perPage = 10
// 	}

// 	offset := (page - 1) * perPage

// 	alihMedia, err := svc.repo.GetAllAlihMedia(ctx, perPage, offset)
// 	if err != nil {
// 		return nil, err
// 	}

// 	total, err := svc.repo.GetTotalAlihMedia(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	totalPages := total / perPage
// 	if total%perPage > 0 {
// 		totalPages++
// 	}

// 	return &AlihMediaPagination{
// 		Data:       alihMedia,
// 		Total:      total,
// 		Page:       page,
// 		PerPage:    perPage,
// 		TotalPages: totalPages,
// 	}, nil
// }

func (svc *alihMediaService) GetByID(ctx context.Context, id int) (*models.AlihMediaJoin, error) {
	if id <= 0 {
		return nil, errors.New("ID must can't be negative")
	}

	alihMedia, err := svc.repo.GetAlihMediaByIDKunjungan(ctx, id)
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

func (svc *alihMediaService) CreateAndCheckAlihMedia(ctx context.Context, kunjunganID int) error {
	kunjungan, err := svc.kunjunganRepo.GetKunjunganByID(ctx, kunjunganID)

	log.Println(kunjungan)

	if err != nil {
		log.Printf("Error getting kunjungan: %v", err)
		return err
	}
	if kunjungan == nil {
		return nil
	}

	kasus, err := svc.kasusRepo.GetKasusByID(ctx, kunjungan.IDKasus)

	log.Println(kasus)

	if err != nil {
		log.Printf("Error getting kasus: %v", err)
		return err
	}
	if kasus == nil {
		return errors.New("Kasus not found")
	}

	var masaInaktif int
	var expirationDate time.Time

	switch kunjungan.JenisKunjungan {
	case "RI":
		masaInaktif = kasus.MasaInaktifRI
		expirationDate = kunjungan.TglMasuk.AddDate(masaInaktif, 0, 0)
	case "RJ":
		masaInaktif = kasus.MasaInaktifRJ
		expirationDate = kunjungan.TglMasuk.AddDate(masaInaktif, 0, 0)
	default:
		return errors.New("Jenis kunjungan invalid")
	}

	now := time.Now()
	if now.After(expirationDate) {
		// Use GetAlihMediaByID to check if alih_media exists
		existing, err := svc.repo.GetAlihMediaByID(ctx, kunjunganID)
		if err != nil {
			log.Printf("Error checking existing alih media: %v", err)
			return err
		}

		if existing == nil {
			alihMedia := models.AlihMedia{
				ID:         kunjunganID,
				TglLaporan: nil, // Use current time instead of nil
				Status:     "belum di alih media",
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}

			_, err = svc.repo.CreateAlihMedia(ctx, &alihMedia)
			if err != nil {
				log.Printf("Error creating alih media: %v", err)
				return err
			}
			log.Printf("Created alih media for kunjungan ID: %d", kunjunganID)
		} else {
			log.Printf("Alih media already exists for kunjungan ID: %d", kunjunganID)
		}
	} else {
		log.Printf("Kunjungan ID %d not expired yet (expires: %s)", kunjunganID, expirationDate.Format("2006-01-02"))
	}

	return nil
}

func (svc *alihMediaService) CheckAllExpiredKunjungan(ctx context.Context) error {
	monthsThreshold := 3

	total, err := svc.kunjunganRepo.GetTotalKunjungan(ctx)
	if err != nil {
		return fmt.Errorf("Failed to get total count: %v", err)
	}

	if total == 0 {
		log.Println("No potentially expired kunjungan found")
		return nil
	}

	workerCount := 10
	pageSize := 100
	totalPages := (total + pageSize - 1) / pageSize

	log.Printf("Processing %d potentially expired kunjungan across %d pages", total, totalPages)

	jobs := make(chan []*models.Kunjungan, workerCount)
	errorChan := make(chan error, totalPages)
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go svc.worker(ctx, jobs, errorChan, &wg)
	}

	for page := 0; page < totalPages; page++ {
		offset := page * pageSize
		kunjunganList, err := svc.kunjunganRepo.GetPotentiallyExpiredKunjungan(ctx, monthsThreshold, pageSize, offset)
		if err != nil {
			errorChan <- fmt.Errorf("error getting page %d: %v", page, err)
			continue
		}
		jobs <- kunjunganList
	}

	close(jobs)
	wg.Wait()
	close(errorChan)

	var errorList []string
	for err := range errorChan {
		if err != nil {
			errorList = append(errorList, err.Error())
		}
	}

	if len(errorList) > 0 {
		return fmt.Errorf("errors during processing: %s", strings.Join(errorList, "; "))
	}

	return nil
}

func (svc *alihMediaService) worker(ctx context.Context, jobs <-chan []*models.Kunjungan, errors chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	for kunjunganList := range jobs {
		for _, kunjungan := range kunjunganList {
			select {
			case <-ctx.Done():
				errors <- ctx.Err()
				return
			default:
				err := svc.CreateAndCheckAlihMedia(ctx, kunjungan.ID)
				if err != nil {
					errors <- fmt.Errorf("Kunjungan ID %d: %v", kunjungan.ID, err.Error())
				}
			}
		}
	}
}

func (svc *alihMediaService) Export(ctx context.Context) ([]byte, error) {
	data, err := svc.repo.GetAllAlihMediaForExport(ctx)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		log.Println("[Export] Tidak ada data alih_media")
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
