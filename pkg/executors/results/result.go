package results

// TestCase result for a run
type TestCase struct {
	Id   string `json:"id"`
	Pass bool   `json:"pass"`
}

// NewTestCaseFail returns a failed test
func NewTestCaseFail(id string) TestCase {
	return NewTestCaseResult(id, false)
}

// NewTestCaseResult return a new TestCase instance
func NewTestCaseResult(id string, pass bool) TestCase {
	return TestCase{
		Id:   id,
		Pass: pass,
	}
}
