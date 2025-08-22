package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/repositories/v2"
)

type DokumenService interface {
	UploadDokumen(ctx context.Context, idKunjungan int, file multipart.File, header *multipart.FileHeader) (*models.Dokumen, error)
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

	if _, err := dst.ReadFrom(file); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	dokumen := models.Dokumen{
		IDKunjungan: idKunjungan,
		Nama:        header.Filename,
		Path:        filePath,
		CreatedAt:   time.Now(),
	}

	return svc.repo.CreateDokumen(ctx, dokumen)
}
