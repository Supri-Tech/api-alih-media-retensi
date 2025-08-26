package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
)

type AlihMediaStatRepository interface {
	AlihMediaStat(ctx context.Context) (*models.AlihMediaStats, error)
}

type alihMediaStatRepository struct {
	db *sql.DB
}

func NewRepoAlihMedia(db *sql.DB) AlihMediaStatRepository {
	return &alihMediaStatRepository{
		db: db,
	}
}

func (repo *alihMediaStatRepository) AlihMediaStat(ctx context.Context) (*models.AlihMediaStats, error) {
	query := `
		SELECT 
			COUNT(*) AS Total,
			COUNT(CASE WHEN Status = 'sudah di alih media' THEN 1 END) AS Aktif,
			COUNT(CASE WHEN Status = 'belum di alih media' THEN 1 END) AS Inaktif
		FROM alih_media
	`

	var stats models.AlihMediaStats
	row := repo.db.QueryRowContext(ctx, query)
	err := row.Scan(&stats.Total, &stats.Aktif, &stats.Inaktif)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	return &stats, nil
}
