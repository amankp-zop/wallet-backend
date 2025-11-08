package repository

import (
	"context"
	"database/sql"

	"github.com/amankp-zop/wallet/internal/domain"
)

type mysqlUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &mysqlUserRepository{
		db: db,
	}
}

func (r *mysqlUserRepository) Create(ctx context.Context, user *domain.User) error {
	query := "INSERT INTO users (name, email, password) VALUES (?, ?, ?)"
	result, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = id

	return nil
}

func (r *mysqlUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := "SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = ?"
	row := r.db.QueryRowContext(ctx, query, email)

	var user domain.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *mysqlUserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := "SELECT id, name, email, created_at, updated_at FROM users WHERE id = ?"
	row := r.db.QueryRowContext(ctx, query, id)

	var user domain.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}