package repositories

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
)

type AlihMediaRepository interface {
	GetStatistikAlihMedia(ctx context.Context) (int, int, int, error)
	GetAllAlihMedia(ctx context.Context, limit, offset int) ([]*models.AlihMediaJoin, error)
	FindAlihMedia(ctx context.Context, filter map[string]interface{}) ([]*models.AlihMediaJoin, error)
	GetAlihMediaByIDKunjungan(ctx context.Context, id int) (*models.AlihMediaJoin, error)
	GetAlihMediaByID(ctx context.Context, id int) (*models.AlihMedia, error)
	GetTotalAlihMedia(ctx context.Context) (int, error)
	CreateAlihMedia(ctx context.Context, alihMedia *models.AlihMedia) (*models.AlihMedia, error)
	UpdateAlihMedia(ctx context.Context, alihMedia models.AlihMedia) (*models.AlihMedia, error)
	DeleteAlihMedia(ctx context.Context, id int) error
	GetAllAlihMediaForExport(ctx context.Context) ([]*models.AlihMediaJoin, error)
}

type alihMediaRepository struct {
	db *sql.DB
}

func NewRepoAlihMedia(db *sql.DB) AlihMediaRepository {
	return &alihMediaRepository{
		db: db,
	}
}

func (repo *alihMediaRepository) GetStatistikAlihMedia(ctx context.Context) (int, int, int, error) {
	var total, sudah, belum int

	if err := repo.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM alih_media`).Scan(&total); err != nil {
		return 0, 0, 0, err
	}

	if err := repo.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM alih_media WHERE Status = 'Sudah dialih media'`).Scan(&sudah); err != nil {
		return 0, 0, 0, err
	}

	if err := repo.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM alih_media WHERE Status = 'Belum dialih media'`).Scan(&belum); err != nil {
		return 0, 0, 0, err
	}

	return total, sudah, belum, nil
}

func (repo *alihMediaRepository) GetAllAlihMedia(ctx context.Context, limit, offset int) ([]*models.AlihMediaJoin, error) {
	query := `
	SELECT
		alih_media.Id AS Id,
		TglLaporan,
		alih_media.Status AS Status,
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
	FROM alih_media
	INNER JOIN kunjungan ON kunjungan.Id = alih_media.Id
	INNER JOIN pasien ON pasien.Id = kunjungan.IdPasien
	INNER JOIN kasus ON kasus.Id = kunjungan.IdKasus
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

	var alihMedia []*models.AlihMediaJoin
	for rows.Next() {
		var am models.AlihMediaJoin
		err := rows.Scan(
			&am.ID,
			&am.TglLaporan,
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
		alihMedia = append(alihMedia, &am)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return alihMedia, nil
}

func (repo *alihMediaRepository) FindAlihMedia(ctx context.Context, filter map[string]interface{}) ([]*models.AlihMediaJoin, error) {
	query := `
	SELECT
		alih_media.Id AS Id,
		TglLaporan,
		alih_media.Status AS Status,
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
	FROM alih_media
	INNER JOIN kunjungan ON kunjungan.Id = alih_media.Id
	INNER JOIN pasien ON pasien.Id = kunjungan.IdPasien
	INNER JOIN kasus ON kasus.Id = kunjungan.IdKasus
	WHERE 1=1
	`

	var args []interface{}

	if noRM, ok := filter["NoRM"]; ok {
		query += " AND pasien.NoRM LIKE ?"
		args = append(args, "%"+noRM.(string)+"%")
	}
	if name, ok := filter["NamaPasien"]; ok {
		query += " AND pasien.NamaPasien LIKE ?"
		args = append(args, "%"+name.(string)+"%")
	}

	if limit, ok := filter["Limit"]; ok {
		query += " LIMIT ?"
		args = append(args, limit)
	} else {
		query += " LIMIT 100"
	}

	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*models.AlihMediaJoin
	for rows.Next() {
		var am models.AlihMediaJoin
		err := rows.Scan(
			&am.ID,
			&am.TglLaporan,
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
		results = append(results, &am)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *alihMediaRepository) GetTotalAlihMedia(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM alih_media`

	var count int
	err := repo.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *alihMediaRepository) GetAlihMediaByIDKunjungan(ctx context.Context, id int) (*models.AlihMediaJoin, error) {
	query := `
	SELECT 
		alih_media.Id AS Id,
		TglLaporan,
		alih_media.Status AS Status,
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
		alih_media
	INNER JOIN kunjungan ON kunjungan.Id = alih_media.Id
	INNER JOIN pasien ON pasien.Id = kunjungan.IdPasien
	INNER JOIN kasus ON kasus.Id = kunjungan.IdKasus
	WHERE alih_media.Id =  ?
	LIMIT 1
	`

	var alihMedia models.AlihMediaJoin
	var tglLaporan sql.NullTime

	row := repo.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&alihMedia.ID,
		&tglLaporan,
		&alihMedia.Status,
		&alihMedia.JenisKunjungan,
		&alihMedia.NoRM,
		&alihMedia.NamaPasien,
		&alihMedia.JenisKelamin,
		&alihMedia.TglLahir,
		&alihMedia.Alamat,
		&alihMedia.StatusPasien,
		&alihMedia.JenisKasus,
		&alihMedia.MasaAktifRi,
		&alihMedia.MasaInaktifRi,
		&alihMedia.MasaAktifRj,
		&alihMedia.MasaInaktifRj,
		&alihMedia.InfoLain,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if tglLaporan.Valid {
		alihMedia.TglLaporan = &tglLaporan.Time
	} else {
		alihMedia.TglLaporan = nil
	}

	return &alihMedia, nil
}

func (repo *alihMediaRepository) GetAlihMediaByID(ctx context.Context, id int) (*models.AlihMedia, error) {
	query := `
	SELECT 
		Id,
		TglLaporan,
		Status,
		CreatedAt,
		UpdatedAt
	FROM
		alih_media
	WHERE Id =  ?
	LIMIT 1
	`

	log.Printf("Executing query: %s with ID: %d", query, id)

	var alihMedia models.AlihMedia
	var tglLaporan, createdAt, updatedAt sql.NullTime

	row := repo.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&alihMedia.ID,
		&tglLaporan,
		&alihMedia.Status,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		log.Printf("GetAlihMediaByID error: %v, ID: %d", err, id)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if tglLaporan.Valid {
		alihMedia.TglLaporan = &tglLaporan.Time
	} else {
		alihMedia.TglLaporan = nil
	}

	if createdAt.Valid {
		alihMedia.CreatedAt = createdAt.Time
	} else {
		alihMedia.CreatedAt = time.Now()
	}

	if updatedAt.Valid {
		alihMedia.UpdatedAt = updatedAt.Time
	} else {
		alihMedia.UpdatedAt = time.Now()
	}

	log.Printf("Successfully found alih_media for ID: %d", id)

	return &alihMedia, nil
}

func (repo *alihMediaRepository) CreateAlihMedia(ctx context.Context, alihMedia *models.AlihMedia) (*models.AlihMedia, error) {
	query := `
	INSERT INTO alih_media(Id, TglLaporan, Status)
	VALUES (?,?,?)
	`

	result, err := repo.db.ExecContext(
		ctx,
		query,
		&alihMedia.ID,
		&alihMedia.TglLaporan,
		&alihMedia.Status,
	)

	if err != nil {
		return nil, err
	}

	_, err = result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return alihMedia, nil
}

func (repo *alihMediaRepository) UpdateAlihMedia(ctx context.Context, alihMedia models.AlihMedia) (*models.AlihMedia, error) {
	query := `
	UPDATE alih_media
	SET TglLaporan = ?, Status = ?
	WHERE Id = ?
	`

	_, err := repo.db.ExecContext(ctx, query, alihMedia.TglLaporan, alihMedia.Status, alihMedia.ID)
	if err != nil {
		return nil, err
	}

	return &alihMedia, nil
}

func (repo *alihMediaRepository) DeleteAlihMedia(ctx context.Context, id int) error {
	query := `
	DELETE FROM alih_media
	WHERE Id = ?
	`

	_, err := repo.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (repo *alihMediaRepository) GetAllAlihMediaForExport(ctx context.Context) ([]*models.AlihMediaJoin, error) {
	query := `
		SELECT
			alih_media.Id AS Id,
			alih_media.TglLaporan,
			alih_media.Status,
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
		FROM alih_media
		INNER JOIN kunjungan ON kunjungan.Id = alih_media.Id
		INNER JOIN pasien ON pasien.Id = kunjungan.IdPasien
		INNER JOIN kasus ON kasus.Id = kunjungan.IdKasus
		ORDER BY alih_media.TglLaporan DESC
	`

	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*models.AlihMediaJoin
	for rows.Next() {
		var am models.AlihMediaJoin
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
