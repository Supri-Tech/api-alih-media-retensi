package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
)

type DokumenRepository interface {
	CreateDokumen(ctx context.Context, dokumen models.Dokumen) (*models.Dokumen, error)
}

type dokumenRepository struct {
	db *sql.DB
}

func NewRepoDokumen(db *sql.DB) DokumenRepository {
	return &dokumenRepository{
		db: db,
	}
}

func (repo *dokumenRepository) CreateDokumen(ctx context.Context, dokumen models.Dokumen) (*models.Dokumen, error) {
	query := `
	INSERTT INTO dokumen(IdKunjungan, Nama, Path)
	VALUES (?,?,?)
	`

	dokumen.CreatedAt = time.Now()
	result, err := repo.db.ExecContext(ctx, query, dokumen.IDKunjungan, dokumen.Nama, dokumen.Path)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	dokumen.ID = int(id)
	return &dokumen, nil
}
