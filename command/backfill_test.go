package command

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

	"github.com/avadhutp/phileas/lib"
	"github.com/spf13/cobra"
)

func TestBackfillPhileas(t *testing.T) {
	oldLibNewCfg := libNewCfg
	oldLibGetDB := libGetDB
	oldLibNewInstaAPI := libNewInstaAPI
	oldInstaAPIBackfill := instaAPIBackfill

	defer func() {
		libNewCfg = oldLibNewCfg
		libGetDB = oldLibGetDB
		libNewInstaAPI = oldLibNewInstaAPI
		instaAPIBackfill = oldInstaAPIBackfill
	}()

	libNewCfg = func(string) *lib.Cfg {
		return &lib.Cfg{}
	}

	libGetDB = func(*lib.Cfg) *gorm.DB {
		return nil
	}

	backfillCalled := false
	instaAPIBackfill = func(*lib.InstaAPI, string) {
		backfillCalled = true
	}

	cmd := &cobra.Command{}
	args := []string{}
	backfillPhileas(cmd, args)

	assert.True(t, backfillCalled)
}
