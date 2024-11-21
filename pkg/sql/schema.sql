CREATE TABLE users (
    id CHAR(36) PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    company_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);


CREATE TABLE test_runs (
    id CHAR(36) PRIMARY KEY,
    user_id CHAR(36) NOT NULL REFERENCES users(id),
    configuration JSONB NOT NULL,
    discovery_model JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS test_test_case_results (
    id CHAR(36) PRIMARY KEY,
    test_case_id CHAR(36),
    test_run_id CHAR(36) NOT NULL REFERENCES test_runs (id),
    pass BOOLEAN NOT NULL,
    fail TEXT[] NOT NULL,
    detail TEXT NOT NULL,
    ref_uri VARCHAR(255) NOT NULL,
    endpoint TEXT NOT NULL,
    api TEXT NOT NULL,
    api_version VARCHAR(255) NOT NULL,
    http_status TEXT NOT NULL,
    test_run_tracking_id TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);
