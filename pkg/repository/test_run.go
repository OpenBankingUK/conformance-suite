package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/lib/pq"
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

type TestCaseResult struct {
	ID         string
	TestCaseID string
	Pass       bool
	Fail       []string
	Detail     string
	TestRunID  string
	RefURI     string
	CreatedAt  time.Time
}

func NewTestRunRepository(db *sql.DB) TestRunRepository {
	return TestRunRepository{db: db}
}

func (r TestRunRepository) GetByID(ctx context.Context, testRunID string) (TestRun, error) {
	var testRun TestRun
	query := `SELECT id, discovery_model, configuration, user_id, created_at FROM test_runs WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, testRunID)
	err := row.Scan(&testRun.ID, &testRun.DiscoveryModel, &testRun.Configuration, &testRun.UserID, &testRun.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return TestRun{}, nil
		}
		return TestRun{}, err
	}
	return testRun, nil
}

func (r TestRunRepository) CreateTestRun(ctx context.Context, testRun TestRun) error {
	query := `INSERT INTO test_runs (id, discovery_model, configuration, user_id, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, testRun.ID, testRun.DiscoveryModel, testRun.Configuration, testRun.UserID, testRun.CreatedAt)
	return err
}

func (r TestRunRepository) CreateTestCaseResult(ctx context.Context, testCaseResult TestCaseResult) error {
	query := `INSERT INTO test_test_case_results (id, test_case_id, test_run_id, pass, fail, detail, ref_uri, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query, testCaseResult.ID, testCaseResult.TestCaseID, testCaseResult.TestRunID, testCaseResult.Pass, pq.Array(testCaseResult.Fail), testCaseResult.Detail, testCaseResult.RefURI, testCaseResult.CreatedAt)
	return err
}
