package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
)

type RetensiRepository interface {
	GetAllRetensi(ctx context.Context, limit, offset int) ([]*models.RetensiJoin, error)
	GetRetensiByID(ctx context.Context, id int) (*models.RetensiJoin, error)
	GetTotalRetensi(ctx context.Context) (int, error)
	CreateRetensi(ctx context.Context, retensi *models.Retensi) (*models.Retensi, error)
	UpdateRetensi(ctx context.Context, retensi models.Retensi) (*models.Retensi, error)
	DeleteRetensi(ctx context.Context, id int) error
}

type retensiRepository struct {
	db *sql.DB
}

func NewRepoRetensi(db *sql.DB) RetensiRepository {
	return &retensiRepository{
		db: db,
	}
}

func (repo *retensiRepository) GetAllRetensi(ctx context.Context, limit, offset int) ([]*models.RetensiJoin, error) {
	query := `
	SELECT
		retensi.Id AS Id,
		TglLaporan,
		retensi.Status AS Status,
		JenisKunjungan,
		pasien.NoRM AS NoRM,
		NamaPasien,
		JenisKelamin,
		TglLahir,
		Alamat,
		pasien.Status AS StatusPasien,
		JenisKasus,
		MasaAktifRi,
		MasaInaktifRi,
		MasaAktifRj,
		MasaInaktifRj,
		InfoLain
	FROM
		retensi
	INNER JOIN
		kunjungan
	ON
		kunjungan.Id = retensi.Id
	INNER JOIN
		pasien
	ON
		pasien.Id = kunjungan.IdPasien
	INNER JOIN
		kasus
	ON
		kasus.Id = kunjungan.IdKasus
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

	var retensi []*models.RetensiJoin
	for rows.Next() {
		var rts models.RetensiJoin
		err := rows.Scan(
			&rts.ID,
			&rts.TglLaporan,
			&rts.Status,
			&rts.JenisKunjungan,
			&rts.NoRM,
			&rts.NamaPasien,
			&rts.JenisKelamin,
			&rts.TglLahir,
			&rts.Alamat,
			&rts.StatusPasien,
			&rts.JenisKasus,
			&rts.MasaAktifRi,
			&rts.MasaInaktifRi,
			&rts.MasaAktifRj,
			&rts.MasaInaktifRj,
			&rts.InfoLain,
		)
		if err != nil {
			return nil, err
		}
		retensi = append(retensi, &rts)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return retensi, nil
}

func (repo *retensiRepository) GetTotalRetensi(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM retensi`

	var count int
	err := repo.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *retensiRepository) GetRetensiByID(ctx context.Context, id int) (*models.RetensiJoin, error) {
	query := `
	SELECT 
		retensi.Id AS Id,
		TglLaporan,
		retensi.Status AS Status,
		JenisKunjungan,
		pasien.NoRM AS NoRM,
		NamaPasien,
		JenisKelamin,
		TglLahir,
		Alamat,
		pasien.Status AS StatusPasien,
		JenisKasus,
		MasaAktifRi,
		MasaInaktifRi,
		MasaAktifRj,
		MasaInaktifRj,
		InfoLain
	FROM
		retensi
	INNER JOIN
		kunjungan
	ON
		kunjungan.Id = retensi.Id
	INNER JOIN
		pasien
	ON
		pasien.Id = kunjungan.IdPasien
	INNER JOIN
		kasus
	ON
		kasus.Id = kunjungan.IdKasus
	WHERE retensi.Id =  ?
	LIMIT 1
	`

	var retensi models.RetensiJoin
	row := repo.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&retensi.ID,
		&retensi.TglLaporan,
		&retensi.Status,
		&retensi.JenisKunjungan,
		&retensi.NoRM,
		&retensi.NamaPasien,
		&retensi.JenisKelamin,
		&retensi.TglLahir,
		&retensi.Alamat,
		&retensi.StatusPasien,
		&retensi.JenisKasus,
		&retensi.MasaAktifRi,
		&retensi.MasaInaktifRi,
		&retensi.MasaAktifRj,
		&retensi.MasaInaktifRj,
		&retensi.InfoLain,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	return &retensi, nil
}

func (repo *retensiRepository) CreateRetensi(ctx context.Context, retensi *models.Retensi) (*models.Retensi, error) {
	query := `
	INSERT INTO retensi(Id, TglLaporan, Status)
	VALUES (?,?,?)
	`

	result, err := repo.db.ExecContext(
		ctx,
		query,
		&retensi.ID,
		&retensi.TglLaporan,
		&retensi.Status,
	)

	if err != nil {
		return nil, err
	}

	_, err = result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return retensi, nil
}

func (repo *retensiRepository) UpdateRetensi(ctx context.Context, retensi models.Retensi) (*models.Retensi, error) {
	query := `
	UPDATE retensi
	SET TglLaporan = ?, Status = ?
	WHERE Id = ?
	`

	_, err := repo.db.ExecContext(ctx, query, retensi.TglLaporan, retensi.Status, retensi.ID)
	if err != nil {
		return nil, err
	}

	return &retensi, nil
}

func (repo *retensiRepository) DeleteRetensi(ctx context.Context, id int) error {
	query := `
	DELETE FROM retensi
	WHERE Id = ?
	`

	_, err := repo.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
