package results

// Test result for a run
type Test struct {
	Id   string `json:"id"`
	Pass bool   `json:"pass"`
}

// NewTestFailResult returns a failed test
func NewTestFailResult(id string) Test {
	return NewTestResult(id, false)
}

// NewTestResult return a new Test instance
func NewTestResult(id string, pass bool) Test {
	return Test{
		Id:   id,
		Pass: pass,
	}
}
