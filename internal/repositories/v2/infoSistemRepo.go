package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
)

type InfoSistemRepository interface {
	GetAllInfoSistem(ctx context.Context, infoSistem models.InfoSistem) (*models.InfoSistem, error)
	GetInfoSistem(ctx context.Context, id int) (*models.InfoSistem, error)
	CreateInfoSistem(ctx context.Context, infoSistem models.InfoSistem) (*models.InfoSistem, error)
	UpdateInfoSistem(ctx context.Context, infoSistem models.InfoSistem) (*models.InfoSistem, error)
}

type infoSistemRepository struct {
	db *sql.DB
}

func NewRepoInfoSistem(db *sql.DB) InfoSistemRepository {
	return &infoSistemRepository{
		db: db,
	}
}

func (repo *infoSistemRepository) GetAllInfoSistem(ctx context.Context, infoSistem models.InfoSistem) (*models.InfoSistem, error) {
	query := `
	SELECT NamaAplikasi, Logo, CreatedAt
	FROM info_sistem
	LIMIT 1
	`

	row := repo.db.QueryRowContext(ctx, query)
	err := row.Scan(
		&infoSistem.NamaAplikasi,
		&infoSistem.Logo,
		&infoSistem.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}
	return &infoSistem, nil
}

func (repo *infoSistemRepository) GetInfoSistem(ctx context.Context, id int) (*models.InfoSistem, error) {
	query := `
	SELECT NamaAplikasi, Logo, CreatedAt
	FROM info_sistem
	WHERE NamaAplikasi = ?
	LIMIT 1
	`

	var infoSistem models.InfoSistem
	row := repo.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&infoSistem.NamaAplikasi,
		&infoSistem.Logo,
		&infoSistem.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}
	return &infoSistem, nil
}

func (repo *infoSistemRepository) CreateInfoSistem(ctx context.Context, infoSistem models.InfoSistem) (*models.InfoSistem, error) {
	query := `
	INSERT INTO info_sistem(NamaAplikasi, Logo, CreatedAt, UpdatedAt)
	VALUES (?,?,?,?)
	`

	infoSistem.CreatedAt = time.Now()
	infoSistem.UpdatedAt = time.Now()
	result, err := repo.db.ExecContext(ctx, query, infoSistem.NamaAplikasi, infoSistem.Logo, infoSistem.CreatedAt, infoSistem.UpdatedAt)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	infoSistem.ID = int(id)
	return &infoSistem, nil
}

func (repo *infoSistemRepository) UpdateInfoSistem(ctx context.Context, infoSistem models.InfoSistem) (*models.InfoSistem, error) {
	query := `
	UPDATE info_sistem
	SET NamaAplikasi = ?, Logo = ?, UpdatedAt = ?
	WHERE Id = ?
	`

	_, err := repo.db.ExecContext(ctx, query, infoSistem.NamaAplikasi, infoSistem.Logo, infoSistem.UpdatedAt, infoSistem.ID)
	if err != nil {
		return nil, err
	}

	return &infoSistem, nil
}
