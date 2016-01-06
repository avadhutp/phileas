package command

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/avadhutp/phileas/lib"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
)

func TestStartPhileas(t *testing.T) {
	oldLibNewCfg := libNewCfg
	oldLibGetDB := libGetDB
	oldLibNewInstaAPI := libNewInstaAPI
	oldLibNewService := libNewService
	oldServiceRun := serviceRun

	defer func() {
		libNewCfg = oldLibNewCfg
		libGetDB = oldLibGetDB
		libNewInstaAPI = oldLibNewInstaAPI
		libNewService = oldLibNewService
		serviceRun = oldServiceRun
	}()

	libNewCfg = func(string) *lib.Cfg {
		return &lib.Cfg{}
	}

	libGetDB = func(*lib.Cfg) *gorm.DB {
		return nil
	}

	serviceRunCalled := false
	serviceRun = func(*gin.Engine, ...string) error {
		serviceRunCalled = true
		return nil
	}

	libNewService = func(*lib.Cfg, *gorm.DB, *lib.InstaAPI) *gin.Engine {
		return nil
	}

	cmd := &cobra.Command{}
	args := []string{}
	startPhileas(cmd, args)

	assert.True(t, serviceRunCalled)
}
