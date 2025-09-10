package repositories

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
)

type KunjunganRepository interface {
	GetAllKunjungan(ctx context.Context, limit, offset int) ([]*models.KunjunganJoin, error)
	GetTotalKunjungan(ctx context.Context) (int, error)
	GetKunjunganByID(ctx context.Context, id int) (*models.KunjunganJoin, error)
	CreateKunjungan(ctx context.Context, kunjungan *models.Kunjungan) (*models.Kunjungan, error)
	UpdateKunjungan(ctx context.Context, kunjungan models.Kunjungan) (*models.Kunjungan, error)
	DeleteKunjungan(ctx context.Context, id int) error
	GetPotentiallyExpiredKunjungan(ctx context.Context, monthsThreshold int, limit, offset int) ([]*models.Kunjungan, error)
	GetKunjunganBasicByID(ctx context.Context, id int) (*models.Kunjungan, error)
	UpdateKunjunganStatus(ctx context.Context, id int, status string) error
	GetActiveKunjungan(ctx context.Context) ([]*models.Kunjungan, error)
	GetTotalActiveKunjungan(ctx context.Context) (int, error)
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
	SELECT
		kunjungan.Id,
		pasien.Id AS IDPasien,
		pasien.NamaPasien AS NamaPasien,
		pasien.NoRM AS NoRM,
		pasien.NIK AS NIK,
		pasien.JenisKelamin AS JenisKelamin,
		pasien.TglLahir AS TglLahir,
		pasien.Alamat AS Alamat,
		pasien.Status AS Status,
		kunjungan.TglMasuk AS TglMasuk,
		kunjungan.JenisKunjungan AS JenisKunjungan,
		kasus.Id AS IDKasus,
		kasus.JenisKasus AS JenisKasus,
		kasus.MasaAktifRi AS MasaAktifRi,
		kasus.MasaInaktifRi AS MasaInaktifRi,
		kasus.MasaAktifRj AS MasaAktifRj,
		kasus.MasaInaktifRj AS MasaInaktifRj,
		kasus.InfoLain AS InfoLain,
		dokumen.Path AS path
	FROM kunjungan
	INNER JOIN pasien ON pasien.Id = kunjungan.IdPasien
	INNER JOIN kasus ON kasus.Id = kunjungan.IdKasus
	LEFT JOIN dokumen ON dokumen.IdKunjungan = kunjungan.Id
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
		var path sql.NullString

		err := rows.Scan(
			&k.ID,
			&k.IDPasien,
			&k.NamaPasien,
			&k.NoRM,
			&k.NIK,
			&k.JenisKelamin,
			&k.TglLahir,
			&k.Alamat,
			&k.Status,
			&k.TglMasuk,
			&k.JenisKunjungan,
			&k.IDKasus,
			&k.JenisKasus,
			&k.MasaAktifRi,
			&k.MasaInaktifRi,
			&k.MasaAktifRj,
			&k.MasaInaktifRj,
			&k.InfoLain,
			&path,
		)
		if err != nil {
			return nil, err
		}

		if path.Valid {
			k.Dokumen = path.String
		} else {
			k.Dokumen = ""
		}

		kunjungan = append(kunjungan, &k)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return kunjungan, nil
}

func (repo *kunjunganRepository) GetActiveKunjungan(ctx context.Context) ([]*models.Kunjungan, error) {
	query := `
	SELECT
		Id,
		IdPasien,
		IdKasus,
		TglMasuk,
		JenisKunjungan,
		Status
	WHERE Status = 'aktif'
	ORDER BY TglMasuk ASC
	`

	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*models.Kunjungan{}, nil
		}
		return nil, err
	}
	defer rows.Close()

	var kunjungan []*models.Kunjungan
	for rows.Next() {
		var k models.Kunjungan

		err := rows.Scan(
			&k.ID,
			&k.IDPasien,
			&k.IDKasus,
			&k.TanggalMasuk,
			&k.JenisKunjungan,
			&k.Status,
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

func (repo *kunjunganRepository) GetTotalActiveKunjungan(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM kunjungan WHERE status = 'aktif'`

	var count int
	err := repo.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
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
	SELECT
		kunjungan.Id,
		pasien.Id AS IDPasien,
		pasien.NamaPasien AS NamaPasien,
		pasien.NoRM AS NoRM,
		pasien.NIK AS NIK,
		pasien.JenisKelamin AS JenisKelamin,
		pasien.TglLahir AS TglLahir,
		pasien.Alamat AS Alamat,
		pasien.Status AS Status,
		kunjungan.TglMasuk AS TglMasuk,
		kunjungan.JenisKunjungan AS JenisKunjungan,
		kasus.Id AS IDKasus,
		kasus.JenisKasus AS JenisKasus,
		kasus.MasaAktifRi AS MasaAktifRi,
		kasus.MasaInaktifRi AS MasaInaktifRi,
		kasus.MasaAktifRj AS MasaAktifRj,
		kasus.MasaInaktifRj AS MasaInaktifRj,
		kasus.InfoLain AS InfoLain,
		dokumen.Path AS path
	FROM kunjungan
	INNER JOIN pasien ON pasien.Id = kunjungan.IdPasien
	INNER JOIN kasus ON kasus.Id = kunjungan.IdKasus
	LEFT JOIN dokumen ON dokumen.IdKunjungan = kunjungan.Id
	WHERE kunjungan.Id = ?
	LIMIT 1
	`
	var k models.KunjunganJoin
	var path sql.NullString

	row := repo.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&k.ID,
		&k.IDPasien,
		&k.NamaPasien,
		&k.NoRM,
		&k.NIK,
		&k.JenisKelamin,
		&k.TglLahir,
		&k.Alamat,
		&k.Status,
		&k.TglMasuk,
		&k.JenisKunjungan,
		&k.IDKasus,
		&k.JenisKasus,
		&k.MasaAktifRi,
		&k.MasaInaktifRi,
		&k.MasaAktifRj,
		&k.MasaInaktifRj,
		&k.InfoLain,
		&path,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if path.Valid {
		k.Dokumen = path.String
	} else {
		k.Dokumen = ""
	}

	log.Println(k)

	return &k, nil
}

func (repo *kunjunganRepository) GetPotentiallyExpiredKunjungan(ctx context.Context, monthsThreshold int, limit, offset int) ([]*models.Kunjungan, error) {
	query := `
	SELECT k.Id, k.IdPasien, k.IdKasus, k.TglMasuk, k.JenisKunjungan, k.CreatedAt, k.UpdatedAt
	FROM kunjungan k
	INNER JOIN kasus ks ON k.IdKasus = ks.Id
	WHERE 
		(k.JenisKunjungan = 'RI' AND DATE_ADD(k.TglMasuk, INTERVAL ks.MasaInaktifRI YEAR) <= DATE_ADD(CURDATE(), INTERVAL ? MONTH))
		OR 
		(k.JenisKunjungan = 'RJ' AND DATE_ADD(k.TglMasuk, INTERVAL ks.MasaInaktifRJ YEAR) <= DATE_ADD(CURDATE(), INTERVAL ? MONTH))
	ORDER BY k.TglMasuk ASC
	LIMIT ? OFFSET ?
	`

	rows, err := repo.db.QueryContext(ctx, query, monthsThreshold, monthsThreshold, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var kunjunganList []*models.Kunjungan
	for rows.Next() {
		var k models.Kunjungan
		err := rows.Scan(
			&k.ID,
			&k.IDPasien,
			&k.IDKasus,
			&k.TanggalMasuk,
			&k.JenisKunjungan,
		)
		if err != nil {
			return nil, err
		}
		kunjunganList = append(kunjunganList, &k)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return kunjunganList, nil
}

func (repo *kunjunganRepository) CreateKunjungan(ctx context.Context, kunjungan *models.Kunjungan) (*models.Kunjungan, error) {
	query := `
	INSERT INTO kunjungan(IdPasien, IdKasus, TglMasuk, JenisKunjungan)
	VALUES (?,?,?,?)
	`

	result, err := repo.db.ExecContext(
		ctx,
		query,
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

func (repo *kunjunganRepository) UpdateKunjungan(ctx context.Context, kunjungan models.Kunjungan) (*models.Kunjungan, error) {
	query := `
	UPDATE kunjungan
	SET IdPasien = ?, IdKasus = ?, TglMasuk = ?, JenisKunjungan = ?
	WHERE Id = ?
	`

	_, err := repo.db.ExecContext(ctx, query, kunjungan.IDPasien, kunjungan.IDKasus, kunjungan.TanggalMasuk, kunjungan.JenisKunjungan, kunjungan.ID)
	if err != nil {
		return nil, err
	}

	return &kunjungan, nil
}

func (repo *kunjunganRepository) DeleteKunjungan(ctx context.Context, id int) error {
	query := `
	DELETE FROM kunjungan
	WHERE Id = ?
	`

	_, err := repo.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (repo *kunjunganRepository) UpdateKunjunganStatus(ctx context.Context, id int, status string) error {
	query := `UPDATE kunjungan SET status = ? WHERE id = ?`
	_, err := repo.db.ExecContext(ctx, query, status, id)
	return err
}

func (repo *kunjunganRepository) GetKunjunganBasicByID(ctx context.Context, id int) (*models.Kunjungan, error) {
	query := `SELECT Id, IdPasien, IdKasus, TglMasuk, JenisKunjungan, Status FROM kunjungan WHERE id = ?`

	var k models.Kunjungan
	err := repo.db.QueryRowContext(ctx, query, id).Scan(
		&k.ID, &k.IDPasien, &k.IDKasus, &k.TanggalMasuk, &k.JenisKunjungan, &k.Status,
	)
	if err != nil {
		return nil, err
	}
	return &k, nil
}
