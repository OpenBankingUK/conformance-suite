package repository

import (
	"context"
	"database/sql"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

type User struct {
	ID          string
	FirstName   string
	LastName    string
	Email       string
	CompanyName string
	CreatedAt   time.Time
}

func NewUserRepository(db *sql.DB) UserRepository {
	return UserRepository{db: db}
}

func (r UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, first_name, last_name, company_name, created_at
		FROM users
		WHERE email = $1
	`
	row := r.db.QueryRowContext(ctx, query, email)

	var user User
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.CompanyName, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r UserRepository) GetByID(ctx context.Context, userID string) (User, error) {
	query := `
		SELECT id, first_name, last_name, company_name, created_at
		FROM users
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, userID)

	var user User
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.CompanyName, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, nil // or return a custom error indicating user not found
		}
		return User{}, err
	}

	return user, nil
}

func (r UserRepository) Create(ctx context.Context, user User) error {
	query := `
		INSERT INTO users (id, first_name, last_name, company_name, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.FirstName, user.LastName, user.CompanyName, user.CreatedAt)
	return err
}
