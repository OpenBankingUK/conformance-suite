
ALTER TABLE test_test_case_results ADD COLUMN method VARCHAR(255) DEFAULT NULL;
ALTER TABLE test_test_case_results ADD COLUMN expected_status_code VARCHAR(255) DEFAULT NULL;
ALTER TABLE test_test_case_results ADD COLUMN response_time VARCHAR(255) DEFAULT NULL;
ALTER TABLE test_test_case_results ADD COLUMN response_size_bytes INTEGER DEFAULT NULL;
ALTER TABLE test_test_case_results ADD COLUMN pass BOOLEAN DEFAULT NULL;