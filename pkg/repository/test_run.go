package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

type TestRunRepository struct {
	db *sql.DB
}

type TestRun struct {
	ID             string
	DiscoveryModel json.RawMessage
	Configuration  json.RawMessage
	UserID         string
	CreatedAt      time.Time
}

func NewTestRunRepository(db *sql.DB) TestRunRepository {
	return TestRunRepository{db: db}
}

func (r TestRunRepository) GetByID(ctx context.Context, testRunID string) (TestRun, error) {
	return TestRun{}, nil
}

func (r TestRunRepository) Create(ctx context.Context, testRun TestRun) error {
	return nil
}
