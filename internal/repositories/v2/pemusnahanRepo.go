package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
)

type PemusnahanRepository interface {
	GetStatistikPemusnahan(ctx context.Context) (int, int, int, error)

	GetAllPemusnahan(ctx context.Context, limit, offset int) ([]*models.PemusnahanJoin, error)
	GetPemusnahanByID(ctx context.Context, id int) (*models.PemusnahanJoin, error)
	GetTotalPemusnahan(ctx context.Context) (int, error)
	CreatePemusnahan(ctx context.Context, pemusnahan *models.Pemusnahan) (*models.Pemusnahan, error)
	UpdatePemusnahan(ctx context.Context, pemusnahan models.Pemusnahan) (*models.Pemusnahan, error)
	DeletePemusnahan(ctx context.Context, id int) error
	GetAllPemusnahanForExport(ctx context.Context) ([]*models.PemusnahanJoin, error)
}

type pemusnahanRepository struct {
	db *sql.DB
}

func NewRepoPemusnahan(db *sql.DB) PemusnahanRepository {
	return &pemusnahanRepository{
		db: db,
	}
}

func (repo *pemusnahanRepository) GetStatistikPemusnahan(ctx context.Context) (int, int, int, error) {
	var total, sudah, belum int

	if err := repo.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM pemusnahan`).Scan(&total); err != nil {
		return 0, 0, 0, err
	}

	if err := repo.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM pemusnahan WHERE Status = 'sudah di musnahkan'`).Scan(&sudah); err != nil {
		return 0, 0, 0, err
	}

	if err := repo.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM pemusnahan WHERE Status = 'belum di musnahkan'`).Scan(&belum); err != nil {
		return 0, 0, 0, err
	}

	return total, sudah, belum, nil
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
		return nil, err
	}
	defer rows.Close()

	var pemusnahan []*models.PemusnahanJoin
	for rows.Next() {
		var pms models.PemusnahanJoin
		var tglLaporan sql.NullTime

		err := rows.Scan(
			&pms.ID,
			&tglLaporan,
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

		if tglLaporan.Valid {
			pms.TglLaporan = &tglLaporan.Time
		} else {
			pms.TglLaporan = nil
		}

		pemusnahan = append(pemusnahan, &pms)
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
	WHERE pemusnahan.Id =  ?
	LIMIT 1
	`

	var pemusnahan models.PemusnahanJoin
	var tglLaporan sql.NullTime

	row := repo.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&pemusnahan.ID,
		&tglLaporan,
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

	if tglLaporan.Valid {
		pemusnahan.TglLaporan = &tglLaporan.Time
	} else {
		pemusnahan.TglLaporan = nil
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

func (repo *pemusnahanRepository) GetAllPemusnahanForExport(ctx context.Context) ([]*models.PemusnahanJoin, error) {
	query := `
		SELECT
			pemusnahan.Id AS Id,
			pemusnahan.TglLaporan,
			pemusnahan.Status,
			kunjungan.JenisKunjungan,
			pasien.NoRM,
			pasien.NamaPasien,
			pasien.JenisKelamin,
			pasien.TglLahir,
			pasien.Alamat,
			pasien.Status AS StatusPasien,
			kasus.JenisKasus,
			kasus.MasaAktifRi,
			kasus.MasaInaktifRi,
			kasus.MasaAktifRj,
			kasus.MasaInaktifRj,
			kasus.InfoLain
		FROM pemusnahan
		INNER JOIN kunjungan ON kunjungan.Id = pemusnahan.Id
		INNER JOIN pasien ON pasien.Id = kunjungan.IdPasien
		INNER JOIN kasus ON kasus.Id = kunjungan.IdKasus
		ORDER BY pemusnahan.TglLaporan DESC
	`

	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*models.PemusnahanJoin
	for rows.Next() {
		var am models.PemusnahanJoin
		var tglLaporan sql.NullTime

		err := rows.Scan(
			&am.ID,
			&tglLaporan,
			&am.Status,
			&am.JenisKunjungan,
			&am.NoRM,
			&am.NamaPasien,
			&am.JenisKelamin,
			&am.TglLahir,
			&am.Alamat,
			&am.StatusPasien,
			&am.JenisKasus,
			&am.MasaAktifRi,
			&am.MasaInaktifRi,
			&am.MasaAktifRj,
			&am.MasaInaktifRj,
			&am.InfoLain,
		)
		if err != nil {
			return nil, err
		}

		if tglLaporan.Valid {
			am.TglLaporan = &tglLaporan.Time
		} else {
			am.TglLaporan = nil
		}

		result = append(result, &am)
	}

	return result, nil
}
