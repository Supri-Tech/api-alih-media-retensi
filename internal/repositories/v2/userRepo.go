package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
)

type UserRepository interface {
	GetAllUsers(ctx context.Context, limit, offset int) ([]*models.User, error)
	GetTotalUsers(ctx context.Context) (int, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByID(ctx context.Context, id int) (*models.User, error) // 🔹 baru
	Create(ctx context.Context, user models.User) (*models.User, error)
	UpdateData(ctx context.Context, user models.User) (*models.User, error)
	UpdateStatus(ctx context.Context, user models.User) (*models.User, error)
	UpdateProfile(ctx context.Context, user models.User) (*models.User, error) // 🔹 baru
	UpdatePassword(ctx context.Context, id int, hashedPassword string) error   // 🔹 baru
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

func (repo *userrepository) GetTotalUsers(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM users`

	var count int
	err := repo.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
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

func (repo *userrepository) UpdateData(ctx context.Context, user models.User) (*models.User, error) {
	query := `
		UPDATE users 
		SET Name = ?, Email = ?, Role = ?, Status = ?
		WHERE id = ?
		LIMIT 1
	`

	_, err := repo.db.ExecContext(ctx, query,
		user.Name,
		user.Email,
		user.Role,
		user.Status,
		user.ID,
	)
	if err != nil {
		return nil, err
	}

	return repo.GetByID(ctx, user.ID)
}

func (repo *userrepository) UpdateStatus(ctx context.Context, user models.User) (*models.User, error) {
	query := `
		UPDATE users 
		SET Status = ?
		WHERE id = ?
		LIMIT 1
	`

	_, err := repo.db.ExecContext(ctx, query,
		user.Status,
		user.ID,
	)
	if err != nil {
		return nil, err
	}

	return repo.GetByID(ctx, user.ID)
}

// func (repo *userrepository) UpdateData(ctx context.Context, user models.User) (*models.User, error) {
// 	query := `
// 	UPDATE users
// 	SET Name = ?, Email = ?, Role = ?, Status = ?
// 	WHERE email = ?
// 	LIMIT 1
// 	`

// 	_, err := repo.db.ExecContext(ctx, query, user.Status, user.Email)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return repo.GetByUsername(ctx, user.Email)
// }

// func (repo *userrepository) UpdateStatus(ctx context.Context, user models.User) (*models.User, error) {
// 	query := `
// 	UPDATE users
// 	SET status = ?
// 	WHERE email = ?
// 	`

// 	user.Status = "aktif"
// 	_, err := repo.db.ExecContext(ctx, query, user.Status, user.Email)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return repo.GetByUsername(ctx, user.Email)
// }

func (repo *userrepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	query := `SELECT id, name, email, password, role, status FROM users WHERE id = ? LIMIT 1`
	row := repo.db.QueryRowContext(ctx, query, id)

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

func (repo *userrepository) UpdateProfile(ctx context.Context, user models.User) (*models.User, error) {
	query := `
	UPDATE users 
	SET name = ?, email = ?
	WHERE id = ?
	`
	_, err := repo.db.ExecContext(ctx, query, user.Name, user.Email, user.ID)
	if err != nil {
		return nil, err
	}
	return repo.GetByID(ctx, user.ID)
}

func (repo *userrepository) UpdatePassword(ctx context.Context, id int, hashedPassword string) error {
	query := `
	UPDATE users 
	SET password = ?
	WHERE id = ?
	`
	_, err := repo.db.ExecContext(ctx, query, hashedPassword, id)
	return err
}
