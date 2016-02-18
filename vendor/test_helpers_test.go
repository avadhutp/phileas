package vendor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTestResult(t *testing.T) {
	actual := NewTestResult(1, 1)

	assert.Equal(t, TestResult{1, 1}, actual)
}
