package services

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/repositories/v2"
)

type DokumenService interface {
	UploadDokumen(ctx context.Context, idKunjungan int, file multipart.File, header *multipart.FileHeader) (*models.Dokumen, error)
	UpdateDokumen(ctx context.Context, id int, file multipart.File, header *multipart.FileHeader) (*models.Dokumen, error)
	DeleteDokumen(ctx context.Context, id int) error
}

type dokumenService struct {
	repo repositories.DokumenRepository
}

func NewServiceDokumen(repo repositories.DokumenRepository) DokumenService {
	return &dokumenService{repo: repo}
}

func (svc *dokumenService) UploadDokumen(ctx context.Context, idKunjungan int, file multipart.File, header *multipart.FileHeader) (*models.Dokumen, error) {
	today := time.Now().Format("2006-01-02")
	uploadDir := filepath.Join("uploads", today)

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("Failed to create upload dir: %w", err)
	}

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), header.Filename)
	filePath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err := file.Seek(0, 0); err != nil {
		return nil, err
	}

	// if _, err := dst.ReadFrom(file); err != nil {
	// 	return nil, fmt.Errorf("failed to save file: %w", err)
	// }

	if _, err := io.Copy(dst, file); err != nil {
		return nil, fmt.Errorf("Failed to save file: %w", err)
	}

	dokumen := models.Dokumen{
		IDKunjungan: idKunjungan,
		Nama:        header.Filename,
		Path:        filePath,
		CreatedAt:   time.Now(),
	}

	return svc.repo.CreateDokumen(ctx, dokumen)
}

func (svc *dokumenService) UpdateDokumen(ctx context.Context, id int, file multipart.File, header *multipart.FileHeader) (*models.Dokumen, error) {
	existing, err := svc.repo.GetDokumenByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := os.Remove(existing.Path); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("Failed to delete old file: %w", err)
	}

	today := time.Now().Format("2006-01-02")
	uploadDir := filepath.Join("uploads", today)
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create upload dir: %w", err)
	}

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), header.Filename)
	filePath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err := file.Seek(0, 0); err != nil {
		return nil, err
	}
	// if _, err := dst.ReadFrom(file); err != nil {
	// 	return nil, fmt.Errorf("failed to save file: %w", err)
	// }

	if _, err := io.Copy(dst, file); err != nil {
		return nil, fmt.Errorf("Failed to save file: %w", err)
	}

	existing.Nama = header.Filename
	existing.Path = filePath

	return svc.repo.UpdateDokumen(ctx, *existing)
}

func (svc *dokumenService) DeleteDokumen(ctx context.Context, id int) error {
	existing, err := svc.repo.GetDokumenByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	if err := os.Remove(existing.Path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return svc.repo.DeleteDokumen(ctx, id)
}
