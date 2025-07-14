package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
)

type KunjunganRepository interface {
	GetAllKunjungan(ctx context.Context, limit, offset int) ([]*models.KunjunganJoin, error)
	GetTotalKunjungan(ctx context.Context) (int, error)
	GetKunjunganByID(ctx context.Context, id int) (*models.KunjunganJoin, error)
	CreateKunjungna(ctx context.Context, kunjungan *models.Kunjungan) (*models.Kunjungan, error)
}

type kunjunganRepository struct {
	db *sql.DB
}

func NewRepoKunjungan(db *sql.DB) KunjunganRepository {
	return &kunjunganRepository{
		db: db,
	}
}

func (repo *kunjunganRepository) GetAllKunjungan(ctx context.Context, limit, offset int) ([]*models.KunjunganJoin, error) {
	query := `
	SELECT Id, pasien.NamaPasien AS NamaPasien, pasien.NoRM AS NoRM, pasien.TglLahir AS TglLahir, pasien.Alamat AS Alamat, kasus.JenisKasus AS JenisKasus
	FROM kunjungan
	LIMIT ? OFFSET ?
	`

	rows, err := repo.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}
	defer rows.Close()

	var kunjungan []*models.KunjunganJoin
	for rows.Next() {
		var k models.KunjunganJoin
		err := rows.Scan(
			&k.ID,
			&k.NamaPasien,
			&k.NoRM,
			&k.TglLahir,
			&k.Alamat,
			&k.JenisKasus,
		)
		if err != nil {
			return nil, err
		}
		kunjungan = append(kunjungan, &k)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return kunjungan, nil
}

func (repo *kunjunganRepository) GetTotalKunjungan(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM kunjungan`

	var count int
	err := repo.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *kunjunganRepository) GetKunjunganByID(ctx context.Context, id int) (*models.KunjunganJoin, error) {
	query := `
	SELECT Id, pasien.NamaPasien AS NamaPasien, pasien.NoRM AS NoRM, pasien.TglLahir AS TglLahir, pasien.Alamat AS Alamat, kasus.JenisKasus AS JenisKasus
	FROM kunjungan
	WHERE Id = ?
	LIMIT 1
	`
	var kunjungan models.KunjunganJoin
	row := repo.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&kunjungan.ID,
		&kunjungan.NamaPasien,
		&kunjungan.NoRM,
		&kunjungan.TglLahir,
		&kunjungan.Alamat,
		&kunjungan.JenisKasus,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	return &kunjungan, nil
}

func (repo *kunjunganRepository) CreateKunjungna(ctx context.Context, kunjungan *models.Kunjungan) (*models.Kunjungan, error) {
	query := `
	INSERT INTO kunjungan(Id, IdPasien, IdKasus, TanggalMasuk, JenisKunjungan)
	VALUES (?,?,?,?,?)
	`

	result, err := repo.db.ExecContext(
		ctx,
		query,
		&kunjungan.ID,
		&kunjungan.IDPasien,
		&kunjungan.IDKasus,
		&kunjungan.TanggalMasuk,
		&kunjungan.JenisKunjungan,
	)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	kunjungan.ID = int(id)

	return kunjungan, nil
}
