package repositories

import (
	"context"
	"database/sql"
)

type GeneralRepository interface {
	GetTotalDokumenByYear(ctx context.Context, year int) (int, error)
	GetTotalPasienByYear(ctx context.Context, year int) (int, error)
	GetTotalKasusByYear(ctx context.Context, year int) (int, error)
	GetTotalAlihMediaByYear(ctx context.Context, year int) (int, error)
	GetTotalRetensiByYear(ctx context.Context, year int) (int, error)
	GetTotalPemusnahanByYear(ctx context.Context, year int) (int, error)
	GetMostCommonKasus(ctx context.Context) (string, int, error)
	GetKasusList(ctx context.Context, year int) (map[string]int, error)
	GetTotalKunjunganByYear(ctx context.Context, year int) (int, error)
}

type generalRepository struct {
	db *sql.DB
}

func NewRepoGeneral(db *sql.DB) GeneralRepository {
	return &generalRepository{db: db}
}

func (repo *generalRepository) GetTotalDokumenByYear(ctx context.Context, year int) (int, error) {
	var total int
	err := repo.db.QueryRowContext(ctx, `
		SELECT COUNT(*) 
		FROM kunjungan`).Scan(&total)
	return total, err
}

func (repo *generalRepository) GetTotalPasienByYear(ctx context.Context, year int) (int, error) {
	var total int
	err := repo.db.QueryRowContext(ctx, `
		SELECT COUNT(*) 
		FROM pasien 
		WHERE YEAR(CreatedAt) = ?`, year).Scan(&total)
	return total, err
}

func (repo *generalRepository) GetTotalKasusByYear(ctx context.Context, year int) (int, error) {
	var total int
	err := repo.db.QueryRowContext(ctx, `
		SELECT COUNT(*) 
		FROM kasus`).Scan(&total)
	return total, err
}

func (repo *generalRepository) GetTotalAlihMediaByYear(ctx context.Context, year int) (int, error) {
	var total int
	err := repo.db.QueryRowContext(ctx, `
		SELECT COUNT(*) 
		FROM alih_media 
		WHERE YEAR(CreatedAt) = ?`, year).Scan(&total)
	return total, err
}

func (repo *generalRepository) GetTotalRetensiByYear(ctx context.Context, year int) (int, error) {
	var total int
	err := repo.db.QueryRowContext(ctx, `
		SELECT COUNT(*) 
		FROM retensi 
		WHERE YEAR(CreatedAt) = ?`, year).Scan(&total)
	return total, err
}

func (repo *generalRepository) GetTotalPemusnahanByYear(ctx context.Context, year int) (int, error) {
	var total int
	err := repo.db.QueryRowContext(ctx, `
		SELECT COUNT(*) 
		FROM pemusnahan 
		WHERE YEAR(CreatedAt) = ?`, year).Scan(&total)
	return total, err
}

func (repo *generalRepository) GetMostCommonKasus(ctx context.Context) (string, int, error) {
	var jenisKasus string
	var total int

	query := `
		SELECT k.JenisKasus, COUNT(*) as total
		FROM kunjungan kj
		JOIN kasus k ON kj.IdKasus = k.Id
		GROUP BY k.Id, k.JenisKasus
		ORDER BY total DESC
		LIMIT 1;
	`

	err := repo.db.QueryRowContext(ctx, query).Scan(&jenisKasus, &total)
	return jenisKasus, total, err
}

func (repo *generalRepository) GetKasusList(ctx context.Context, year int) (map[string]int, error) {
	rows, err := repo.db.QueryContext(ctx, `
		SELECT k.JenisKasus, COUNT(kj.Id) as total
		FROM kasus k
		LEFT JOIN kunjungan kj ON kj.IdKasus = k.Id
		GROUP BY k.Id, k.JenisKasus`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var nama string
		var total int
		if err := rows.Scan(&nama, &total); err != nil {
			return nil, err
		}
		result[nama] = total
	}

	return result, nil
}

func (repo *generalRepository) GetTotalKunjunganByYear(ctx context.Context, year int) (int, error) {
	var total int
	err := repo.db.QueryRowContext(ctx, `
		SELECT COUNT(*) 
		FROM kunjungan`).Scan(&total)
	return total, err
}
