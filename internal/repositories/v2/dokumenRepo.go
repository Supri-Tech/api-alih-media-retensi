package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
)

type DokumenRepository interface {
	CreateDokumen(ctx context.Context, dokumen models.Dokumen) (*models.Dokumen, error)
	GetDokumenByID(ctx context.Context, id int) (*models.Dokumen, error)
	UpdateDokumen(ctx context.Context, dokumen models.Dokumen) (*models.Dokumen, error)
	DeleteDokumen(ctx context.Context, id int) error
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
	INSERT INTO dokumen(IdKunjungan, Nama, Path)
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

func (repo *dokumenRepository) GetDokumenByID(ctx context.Context, id int) (*models.Dokumen, error) {
	query := `
	SELECT 
		Id, 
		IdKunjungan, 
		Nama, 
		Path, 
		CreatedAt 
	FROM
		dokumen
	WHERE
		Id = ?
	LIMIT 1
	`

	var dokumen models.Dokumen
	row := repo.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&dokumen.ID,
		&dokumen.IDKunjungan,
		&dokumen.Nama,
		&dokumen.Path,
		&dokumen.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &dokumen, nil
}

func (repo *dokumenRepository) UpdateDokumen(ctx context.Context, dokumen models.Dokumen) (*models.Dokumen, error) {
	query := `
	UPDATE dokumen
	SET Nama = ?, Path = ?
	Where Id = ?
	`

	_, err := repo.db.ExecContext(ctx, query, dokumen.Nama, dokumen.Path, dokumen.ID)
	if err != nil {
		return nil, err
	}

	return &dokumen, nil
}

func (repo *dokumenRepository) DeleteDokumen(ctx context.Context, id int) error {
	query := `DELETE FROM dokumen WHERE Id = ?`
	_, err := repo.db.ExecContext(ctx, query, id)

	return err
}
