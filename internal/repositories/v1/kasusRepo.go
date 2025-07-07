package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v1"
)

type KasusRepository interface {
	GetAllKasus(ctx context.Context, limit, offset int) ([]*models.Kasus, error)
	GetTotalKasus(ctx context.Context) (int, error)
	GetKasusByID(ctx context.Context, id int) (*models.Kasus, error)
	CreateKasus(ctx context.Context, kasus models.Kasus) (*models.Kasus, error)
	UpdateKasus(ctx context.Context, kasus models.Kasus) (*models.Kasus, error)
	DeleteKasus(ctx context.Context, id int) error
}

type kasusrepository struct {
	db *sql.DB
}

func NewRepoKasus(db *sql.DB) KasusRepository {
	return &kasusrepository{
		db: db,
	}
}

func (repo *kasusrepository) GetAllKasus(ctx context.Context, limit, offset int) ([]*models.Kasus, error) {
	query := `
	SELECT Id, JenisKasus, MasaAktifRi, MasaInaktifRi, MasaAktifRj, MasaInaktifRj
	FROM kasus
	LIMIT ? OFFSET ?
	`

	rows, err := repo.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	var kasus []*models.Kasus
	for rows.Next() {
		var k models.Kasus
		err := rows.Scan(
			&k.ID,
			&k.JenisKasus,
			&k.MasaAktifRI,
			&k.MasaInaktifRI,
			&k.MasaAktifRJ,
			&k.MasaInaktifRJ,
		)
		if err != nil {
			return nil, err
		}
		kasus = append(kasus, &k)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return kasus, nil
}

func (repo *kasusrepository) GetTotalKasus(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM kasus`

	var count int
	err := repo.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *kasusrepository) GetKasusByID(ctx context.Context, id int) (*models.Kasus, error) {
	query := `
	SELECT Id, JenisKasus, MasaAktifRi, MasaInaktifRi, MasaAktifRj, MasaInaktifRj, InfoLain
	FROM kasus
	WHERE id = ?
	LIMIT 1
	`

	var kasus models.Kasus
	row := repo.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&kasus.ID, &kasus.JenisKasus, &kasus.MasaAktifRI, &kasus.MasaInaktifRI, &kasus.MasaAktifRJ, &kasus.MasaInaktifRJ, &kasus.InfoLain)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &kasus, nil
}

func (repo *kasusrepository) CreateKasus(ctx context.Context, kasus models.Kasus) (*models.Kasus, error) {
	query := `
	INSERT INTO kasus(JenisKasus, MasaAktifRi, MasaInaktifRi, MasaAktifRj, MasaInaktifRj, InfoLain)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := repo.db.ExecContext(ctx, query, kasus.JenisKasus, kasus.MasaAktifRI, kasus.MasaInaktifRI, kasus.MasaAktifRJ, kasus.MasaInaktifRJ, kasus.InfoLain)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	kasus.ID = int(id)
	return &kasus, nil
}

func (repo *kasusrepository) UpdateKasus(ctx context.Context, kasus models.Kasus) (*models.Kasus, error) {
	query := `
	UPDATE kasus
	SET JenisKasus = ?, MasaAktifRI = ?, MasaInaktifRI = ?, MasaAktifRJ = ?, MasaInaktifRJ = ?, InfoLain = ?
	WHERE id = ?
	`

	_, err := repo.db.ExecContext(ctx, query, kasus.JenisKasus, kasus.MasaAktifRI, kasus.MasaInaktifRI, kasus.MasaAktifRJ, kasus.MasaInaktifRJ, kasus.InfoLain, kasus.ID)
	if err != nil {
		return nil, err
	}

	return &kasus, nil
}

func (repo *kasusrepository) DeleteKasus(ctx context.Context, id int) error {
	query := `
	DELETE FROM kasus
	WHERE id = ?
	`

	_, err := repo.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
