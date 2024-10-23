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
	CompanyName string
	CreatedAt   time.Time
}

func NewUserRepository(db *sql.DB) UserRepository {
	return UserRepository{db: db}
}

func (r UserRepository) GetByID(ctx context.Context, userID string) (User, error) {
	return User{}, nil
}
