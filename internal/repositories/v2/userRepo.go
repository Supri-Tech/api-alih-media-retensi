package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
)

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	Create(ctx context.Context, user models.User) (*models.User, error)
	UpdateStatus(ctx context.Context, user models.User) (*models.User, error)
}

type userrepository struct {
	db *sql.DB
}

func NewRepoUser(db *sql.DB) UserRepository {
	return &userrepository{
		db: db,
	}
}

func (repo *userrepository) GetAllUsers(ctx context.Context, limit, offset int) ([]*models.User, error) {
	query := `
	SELECT
		Id, Name, Email, Role, Status, CreatedAt
	FROM
		users
	LIMIT ?
	OFFSET ?
	`

	rows, err := repo.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Role,
			&user.Status,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *userrepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
	SELECT id, name, email, password, role, status
	FROM users
	WHERE email = ?
	LIMIT 1
	`
	row := repo.db.QueryRowContext(ctx, query, username)

	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (repo *userrepository) Create(ctx context.Context, user models.User) (*models.User, error) {
	query := `
	INSERT INTO users(name, email, password, role, status)
	VALUES (?, ?, ?, ?, ?)
	`

	user.Status = "tidak aktif"
	result, err := repo.db.ExecContext(ctx, query, user.Name, user.Email, user.Password, user.Role, user.Status)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	user.ID = int(id)
	return &user, nil
}

func (repo *userrepository) UpdateStatus(ctx context.Context, user models.User) (*models.User, error) {
	query := `
	UPDATE users 
	SET status = ? 
	WHERE email = ?
	`

	user.Status = "aktif"
	_, err := repo.db.ExecContext(ctx, query, user.Status, user.Email)
	if err != nil {
		return nil, err
	}

	return repo.GetByUsername(ctx, user.Email)
}
