package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
)

type PasienRepository interface {
	GetAllPasien(ctx context.Context, limit, offset int) ([]*models.Pasien, error)
	GetTotalPasien(ctx context.Context) (int, error)
	GetPasienByID(ctx context.Context, id int) (*models.Pasien, error)
	FindPasien(ctx context.Context, filter map[string]string) ([]*models.Pasien, error)
	GetPasienByNoRM(ctx context.Context, noRM string) (*models.Pasien, error)
	GetPasienByName(ctx context.Context, name string) (*models.Pasien, error)
	GetPasienByNIK(ctx context.Context, NIK string) (*models.Pasien, error)
	CreatePasien(ctx context.Context, pasien models.Pasien) (*models.Pasien, error)
	UpdatePasien(ctx context.Context, pasien models.Pasien) (*models.Pasien, error)
	DeletePasien(ctx context.Context, id int) error
}

type pasienRepository struct {
	db *sql.DB
}

func NewRepoPasien(db *sql.DB) PasienRepository {
	return &pasienRepository{
		db: db,
	}
}

func (repo *pasienRepository) GetAllPasien(ctx context.Context, limit, offset int) ([]*models.Pasien, error) {
	query := `
	SELECT Id, NoRM, NamaPasien, JenisKelamin, TglLahir, NIK, Alamat, Status, CreatedAt
	FROM pasien
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

	var pasien []*models.Pasien
	for rows.Next() {
		var p models.Pasien
		err := rows.Scan(
			&p.ID,
			&p.NoRM,
			&p.NamaPasien,
			&p.JenisKelamin,
			&p.TanggalLahir,
			&p.NIK,
			&p.Alamat,
			&p.Status,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		pasien = append(pasien, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pasien, nil
}

func (repo *pasienRepository) GetTotalPasien(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM pasien`

	var count int
	err := repo.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *pasienRepository) GetPasienByID(ctx context.Context, id int) (*models.Pasien, error) {
	query := `
	SELECT Id, NoRM, NamaPasien, JenisKelamin, TglLahir, NIK, Alamat, Status, CreatedAt
	FROM pasien
	WHERE Id = ?
	LIMIT 1
	`

	var pasien models.Pasien
	row := repo.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&pasien.ID,
		&pasien.NoRM,
		&pasien.NamaPasien,
		&pasien.JenisKelamin,
		&pasien.TanggalLahir,
		&pasien.NIK,
		&pasien.Alamat,
		&pasien.Status,
		&pasien.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	return &pasien, nil
}

func (repo *pasienRepository) FindPasien(ctx context.Context, filter map[string]string) ([]*models.Pasien, error) {
	query := `
	SELECT Id, NoRM, NamaPasien, JenisKelamin, TglLahir, NIK, Alamat, Status, CreatedAt
	FROM pasien
	WHERE 1=1
	`

	var args []interface{}

	if noRM, ok := filter["NoRM"]; ok {
		query += " AND NoRM LIKE ?"
		args = append(args, "%"+noRM+"")
	}
	if name, ok := filter["NamaPasien"]; ok {
		query += " AND NamaPasien LIKE ?"
		args = append(args, "%"+name+"%")
	}
	if nik, ok := filter["NIK"]; ok {
		query += " AND NIK LIKE ?"
		args = append(args, "%"+nik+"%")
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

	var pasien []*models.Pasien
	for rows.Next() {
		var p models.Pasien
		err := rows.Scan(
			&p.ID,
			&p.NoRM,
			&p.NamaPasien,
			&p.JenisKelamin,
			&p.TanggalLahir,
			&p.NIK,
			&p.Alamat,
			&p.Status,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		pasien = append(pasien, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pasien, nil
}

func (repo *pasienRepository) GetPasienByNoRM(ctx context.Context, noRM string) (*models.Pasien, error) {
	query := `
	SELECT Id, NoRM, NamaPasien, JenisKelamin, TglLahir, NIK, Alamat, Status, CreatedAt
	FROM pasien
	WHERE NoRM = ?
	LIMIT 1
	`
	var pasien models.Pasien
	row := repo.db.QueryRowContext(ctx, query, noRM)
	err := row.Scan(
		&pasien.ID,
		&pasien.NoRM,
		&pasien.NamaPasien,
		&pasien.JenisKelamin,
		&pasien.TanggalLahir,
		&pasien.NIK,
		&pasien.Alamat,
		&pasien.Status,
		&pasien.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	return &pasien, nil
}

func (repo *pasienRepository) GetPasienByName(ctx context.Context, email string) (*models.Pasien, error) {
	query := `
	SELECT Id, NoRM, NamaPasien, JenisKelamin, TglLahir, NIK, Alamat, Status, CreatedAt
	FROM pasien
	WHERE NamaPasien = ?
	LIMIT 1
	`
	var pasien models.Pasien
	row := repo.db.QueryRowContext(ctx, query, email)
	err := row.Scan(
		&pasien.ID,
		&pasien.NoRM,
		&pasien.NamaPasien,
		&pasien.JenisKelamin,
		&pasien.TanggalLahir,
		&pasien.NIK,
		&pasien.Alamat,
		&pasien.Status,
		&pasien.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}
	return &pasien, nil
}

func (repo *pasienRepository) GetPasienByNIK(ctx context.Context, NIK string) (*models.Pasien, error) {
	query := `
	SELECT Id, NoRM, NamaPasien, JenisKelamin, TglLahir, NIK, Alamat, Status, CreatedAt
	FROM pasien
	WHERE NIK = ?
	LIMIT 1
	`
	var pasien models.Pasien
	row := repo.db.QueryRowContext(ctx, query, NIK)
	err := row.Scan(
		&pasien.ID,
		&pasien.NoRM,
		&pasien.NamaPasien,
		&pasien.JenisKelamin,
		&pasien.TanggalLahir,
		&pasien.NIK,
		&pasien.Alamat,
		&pasien.Status,
		&pasien.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}
	return &pasien, nil
}

func (repo *pasienRepository) CreatePasien(ctx context.Context, pasien models.Pasien) (*models.Pasien, error) {
	query := `
	INSERT INTO pasien(NoRM, NamaPasien, JenisKelamin, TglLahir, NIK, Alamat, Status, CreatedAt) 
	VALUES (?,?,?,?,?,?,?,?)
	`

	result, err := repo.db.ExecContext(
		ctx,
		query,
		&pasien.NoRM,
		&pasien.NamaPasien,
		&pasien.JenisKelamin,
		&pasien.TanggalLahir,
		&pasien.NIK,
		&pasien.Alamat,
		&pasien.Status,
		&pasien.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	pasien.ID = int(id)

	return &pasien, nil
}

func (repo *pasienRepository) UpdatePasien(ctx context.Context, pasien models.Pasien) (*models.Pasien, error) {
	query := `
	UPDATE pasien
	SET 
		NoRM = ?, 
		NamaPasien = ?, 
		JenisKelamin = ?, 
		TglLahir = ?, 
		NIK = ?, 
		Alamat = ?, 
		Status = ?
	WHERE Id = ?
	`

	_, err := repo.db.ExecContext(ctx, query, pasien.NoRM, pasien.NamaPasien, pasien.JenisKelamin, pasien.TanggalLahir, pasien.NIK, pasien.Alamat, pasien.Status, pasien.ID)
	if err != nil {
		return nil, err
	}

	return &pasien, nil
}

func (repo *pasienRepository) DeletePasien(ctx context.Context, id int) error {
	query := `
	DELETE FROM pasien
	WHERE Id = ?
	`

	_, err := repo.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
