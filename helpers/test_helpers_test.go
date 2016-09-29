package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTestResult(t *testing.T) {
	actual := NewTestResult(1, 1)

	l, _ := actual.LastInsertId()
	a, _ := actual.RowsAffected()

	assert.Equal(t, 1, int(a))
	assert.Equal(t, 1, int(l))
}
