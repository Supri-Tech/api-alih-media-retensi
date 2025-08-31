package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
)

type PemusnahanRepository interface {
	GetAllPemusnahan(ctx context.Context, limit, offset int) ([]*models.PemusnahanJoin, error)
	GetPemusnahanByID(ctx context.Context, id int) (*models.PemusnahanJoin, error)
	GetTotalPemusnahan(ctx context.Context) (int, error)
	CreatePemusnahan(ctx context.Context, pemusnahan *models.Pemusnahan) (*models.Pemusnahan, error)
	UpdatePemusnahan(ctx context.Context, pemusnahan models.Pemusnahan) (*models.Pemusnahan, error)
	DeletePemusnahan(ctx context.Context, id int) error
}

type pemusnahanRepository struct {
	db *sql.DB
}

func NewRepoPemusnahan(db *sql.DB) PemusnahanRepository {
	return &pemusnahanRepository{
		db: db,
	}
}

func (repo *pemusnahanRepository) GetAllPemusnahan(ctx context.Context, limit, offset int) ([]*models.PemusnahanJoin, error) {
	query := `
	SELECT
		pemusnahan.Id AS Id,
		TglLaporan,
		pemusnahan.Status AS Status,
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
		pemusnahan
	INNER JOIN
		kunjungan
	ON
		kunjungan.Id = pemusnahan.Id
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

	var pemusnahan []*models.PemusnahanJoin
	for rows.Next() {
		var pms models.PemusnahanJoin
		err := rows.Scan(
			&pms.ID,
			&pms.TglLaporan,
			&pms.Status,
			&pms.JenisKunjungan,
			&pms.NoRM,
			&pms.NamaPasien,
			&pms.JenisKelamin,
			&pms.TglLahir,
			&pms.Alamat,
			&pms.StatusPasien,
			&pms.JenisKasus,
			&pms.MasaAktifRi,
			&pms.MasaInaktifRi,
			&pms.MasaAktifRj,
			&pms.MasaInaktifRj,
			&pms.InfoLain,
		)
		if err != nil {
			return nil, err
		}
		pemusnahan = append(pemusnahan, &pms)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pemusnahan, nil
}

func (repo *pemusnahanRepository) GetTotalPemusnahan(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM pemusnahan`

	var count int
	err := repo.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *pemusnahanRepository) GetPemusnahanByID(ctx context.Context, id int) (*models.PemusnahanJoin, error) {
	query := `
	SELECT 
		pemusnahan.Id AS Id,
		TglLaporan,
		pemusnahan.Status AS Status,
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
		pemusnahan
	INNER JOIN
		kunjungan
	ON
		kunjungan.Id = pemusnahan.Id
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

	var pemusnahan models.PemusnahanJoin
	row := repo.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&pemusnahan.ID,
		&pemusnahan.TglLaporan,
		&pemusnahan.Status,
		&pemusnahan.JenisKunjungan,
		&pemusnahan.NoRM,
		&pemusnahan.NamaPasien,
		&pemusnahan.JenisKelamin,
		&pemusnahan.TglLahir,
		&pemusnahan.Alamat,
		&pemusnahan.StatusPasien,
		&pemusnahan.JenisKasus,
		&pemusnahan.MasaAktifRi,
		&pemusnahan.MasaInaktifRi,
		&pemusnahan.MasaAktifRj,
		&pemusnahan.MasaInaktifRj,
		&pemusnahan.InfoLain,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	return &pemusnahan, nil
}

func (repo *pemusnahanRepository) CreatePemusnahan(ctx context.Context, pemusnahan *models.Pemusnahan) (*models.Pemusnahan, error) {
	query := `
	INSERT INTO pemusnahan(Id, TglLaporan, Status)
	VALUES (?,?,?)
	`

	result, err := repo.db.ExecContext(
		ctx,
		query,
		&pemusnahan.ID,
		&pemusnahan.TglLaporan,
		&pemusnahan.Status,
	)

	if err != nil {
		return nil, err
	}

	_, err = result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return pemusnahan, nil
}

func (repo *pemusnahanRepository) UpdatePemusnahan(ctx context.Context, pemusnahan models.Pemusnahan) (*models.Pemusnahan, error) {
	query := `
	UPDATE pemusnahan
	SET TglLaporan = ?, Status = ?
	WHERE Id = ?
	`

	_, err := repo.db.ExecContext(ctx, query, pemusnahan.TglLaporan, pemusnahan.Status, pemusnahan.ID)
	if err != nil {
		return nil, err
	}

	return &pemusnahan, nil
}

func (repo *pemusnahanRepository) DeletePemusnahan(ctx context.Context, id int) error {
	query := `
	DELETE FROM pemusnahan
	WHERE Id = ?
	`

	_, err := repo.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
