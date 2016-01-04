package lib

import (
	"errors"

	"github.com/stretchr/testify/assert"

	"testing"
)

func TestNewCfg(t *testing.T) {
	oldIniMapTo := iniMapTo
	defer func() { iniMapTo = oldIniMapTo }()

	iniMapTo = func(interface{}, interface{}, ...interface{}) error {
		return nil
	}

	actual := NewCfg("file.ini")

	assert.Equal(t, &Cfg{}, actual)
}

func TestNewCfgHandlError(t *testing.T) {
	oldIniMapTo := iniMapTo
	oldLogErr := logErr
	defer func() {
		iniMapTo = oldIniMapTo
		logErr = oldLogErr
	}()

	expected := errors.New("Test error")
	iniMapTo = func(interface{}, interface{}, ...interface{}) error {
		return expected
	}

	logErrCalled := false
	logErr = func(...interface{}) {
		logErrCalled = true
	}

	actual := NewCfg("file.ini")

	assert.Nil(t, actual)
	assert.True(t, logErrCalled)
}
