package vendor

// TestResult Test struct to help with stubbing insert queries for testdb; this has to be put in the vendor folder seperately because it fails golint (lastId instead of lastID is used). And this cannot be changed as the interface demanding this struct requires the variable/func names to be so.
type TestResult struct {
	lastId       int64
	affectedRows int64
}

func (r TestResult) LastInsertId() (int64, error) {
	return r.lastId, nil
}

func (r TestResult) RowsAffected() (int64, error) {
	return r.affectedRows, nil
}

// NewTestResult Provides a new test result
func NewTestResult(lastID int64, affectedRows int64) TestResult {
	return TestResult{lastID, affectedRows}
}
