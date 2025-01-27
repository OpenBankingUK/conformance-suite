package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/OpenBankingUK/conformance-suite/pkg/discovery"
	"github.com/OpenBankingUK/conformance-suite/pkg/executors/results"
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
	ID                 string
	TestCaseID         string
	Pass               bool
	Fail               []string
	Detail             string
	TestRunID          string
	RefURI             string
	Endpoint           string
	API                string
	APIVersion         string
	HTTPStatus         string
	CreatedAt          time.Time
	ExpectedStatusCode int
	Method             string
	ResponseTime       string
	ResponseSizeBytes  int
}

type PublicTestCaseResult struct {
	APISpecification discovery.ModelAPISpecification `json:"apiSpecification"`
	TestCases        []PublicTestCase                `json:"testCases"`
}

type PublicTestCase struct {
	Id    string `json:"@id"`
	Name  string `json:"name"`
	Input struct {
		Method   string `json:"method"`
		Endpoint string `json:"endpoint"`
	} `json:"input"`
	Expect struct {
		StatusCode int `json:"status-code"`
	} `json:"expect"`
	Meta struct {
		Status  string `json:"status"`
		Metrics struct {
			ResponseTime string `json:"responseTime"`
			ResponseSize string `json:"responseSize"`
		} `json:"metrics"`
	} `json:"meta"`
	Pass       bool            `json:"pass"`
	Metrics    results.Metrics `json:"metrics,omitempty"`
	Fail       []string        `json:"fail,omitempty"`
	Detail     string          `json:"detail"`
	RefURI     string          `json:"refURI"`
	API        string          `json:"-"`
	APIVersion string          `json:"-"`
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

func (r TestRunRepository) GetAllByUserID(ctx context.Context, userID string) ([]TestRun, error) {
	var ret []TestRun
	rows, err := r.db.QueryContext(ctx, `SELECT
	test_runs.id,
	test_runs.discovery_model,
	test_runs.configuration,
	test_runs.user_id,
	test_runs.created_at
	FROM
	test_runs
	JOIN users ON test_runs.user_id = users.id
	WHERE users.id = $1`,
		userID,
	)
	if err != nil {
		return ret, err
	}
	for rows.Next() {
		var testRun TestRun
		if err := rows.Scan(&testRun.ID, &testRun.DiscoveryModel, &testRun.Configuration, &testRun.UserID, &testRun.CreatedAt); err != nil {
			return ret, nil
		}
		ret = append(ret, testRun)
	}

	return ret, nil
}

func (r TestRunRepository) RetrieveRunResults(ctx context.Context, testRunID string) (PublicTestCaseResult, error) {
	var ret PublicTestCaseResult

	// Get all test case results
	rows, err := r.db.QueryContext(ctx, `
        SELECT 
            test_case_id,
            detail,
            pass,
            fail,
            detail,
            ref_uri,
            endpoint,
            api,
            api_version,
            http_status,
			method,
			response_time,
			response_size_bytes
        FROM test_test_case_results 
        WHERE test_run_id = $1`, testRunID)
	if err != nil {
		return ret, fmt.Errorf("fetching test results: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tc PublicTestCase
		var httpStatus string
		if err := rows.Scan(
			&tc.Id,
			&tc.Name,
			&tc.Pass,
			pq.Array(&tc.Fail),
			&tc.Detail,
			&tc.RefURI,
			&tc.Input.Endpoint,
			&ret.APISpecification.Name,
			&ret.APISpecification.Version,
			&httpStatus,
			&tc.Input.Method,
			&tc.Meta.Metrics.ResponseTime,
			&tc.Meta.Metrics.ResponseSize,
		); err != nil {
			return ret, fmt.Errorf("scanning test case: %w", err)
		}

		// Set meta status from http_status
		tc.Expect.StatusCode = 200
		tc.Meta.Status = httpStatus
		ret.TestCases = append(ret.TestCases, tc)
	}

	return ret, nil
}

func (r TestRunRepository) CreateTestRun(ctx context.Context, testRun TestRun) error {
	query := `INSERT INTO test_runs (id, discovery_model, configuration, user_id, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, testRun.ID, testRun.DiscoveryModel, testRun.Configuration, testRun.UserID, testRun.CreatedAt)
	return err
}

func (r TestRunRepository) CreateTestCaseResult(ctx context.Context, testCaseResult TestCaseResult) error {
	query := `
		INSERT INTO test_test_case_results (
			id, 
			test_case_id, 
			test_run_id, 
			pass, 
			fail, 
			detail, 
			ref_uri, 
			endpoint, 
			api, 
			api_version, 
			http_status, 
			method, 
			expected_status_code, 
			response_time, 
			response_size_bytes, 
			created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		)`
	_, err := r.db.ExecContext(
		ctx,
		query,
		testCaseResult.ID,
		testCaseResult.TestCaseID,
		testCaseResult.TestRunID,
		testCaseResult.Pass,
		pq.Array(testCaseResult.Fail),
		testCaseResult.Detail,
		testCaseResult.RefURI,
		testCaseResult.Endpoint,
		testCaseResult.API,
		testCaseResult.APIVersion,
		testCaseResult.HTTPStatus,
		testCaseResult.Method,
		testCaseResult.ExpectedStatusCode,
		testCaseResult.ResponseTime,
		testCaseResult.ResponseSizeBytes,
		testCaseResult.CreatedAt,
	)
	return err
}
