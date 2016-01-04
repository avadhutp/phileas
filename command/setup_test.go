package command

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/spf13/cobra"

	"github.com/jinzhu/gorm"

	"github.com/avadhutp/phileas/lib"
)

func TestSetup(t *testing.T) {
	oldDbSet := dbSet
	oldDBAutoMigrate := dbAutoMigrate
	oldLibNewCfg := libNewCfg
	oldLibGetDB := libGetDB
	defer func() {
		dbSet = oldDbSet
		dbAutoMigrate = oldDBAutoMigrate
		libNewCfg = oldLibNewCfg
		libGetDB = oldLibGetDB
	}()

	libNewCfg = func(string) *lib.Cfg {
		return &lib.Cfg{}
	}

	mockDB := &gorm.DB{}

	libGetDB = func(*lib.Cfg) *gorm.DB {
		return mockDB
	}

	dbSetCalled := false
	dbSet = func(*gorm.DB, string, interface{}) *gorm.DB {
		dbSetCalled = true
		return mockDB
	}

	dbAutoMigrateCalled := false
	dbAutoMigrate = func(*gorm.DB, ...interface{}) *gorm.DB {
		dbAutoMigrateCalled = true
		return mockDB
	}

	cmd := &cobra.Command{}
	args := []string{}
	setup(cmd, args)

	assert.True(t, dbSetCalled)
	assert.True(t, dbAutoMigrateCalled)
}
