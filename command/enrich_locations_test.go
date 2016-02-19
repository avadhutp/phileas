package command

import (
	"github.com/avadhutp/phileas/lib"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestEnrichPhileas(t *testing.T) {
	oldLibNewCfg := libNewCfg
	oldLibGetDB := libGetDB
	oldLibNewEnrichmentService := libNewEnrichmentService
	oldESEnrichLocations := esEnrichLocations

	defer func() {
		libNewCfg = oldLibNewCfg
		libGetDB = oldLibGetDB
		libNewEnrichmentService = oldLibNewEnrichmentService
		esEnrichLocations = oldESEnrichLocations
	}()

	libNewCfg = func(string) *lib.Cfg {
		return &lib.Cfg{}
	}

	libGetDB = func(*lib.Cfg) *gorm.DB {
		return nil
	}

	enrichCalled := false
	esEnrichLocations = func(*lib.EnrichmentService) {
		enrichCalled = true
	}

	cmd := &cobra.Command{}
	args := []string{}
	enrichLocations(cmd, args)

	assert.True(t, enrichCalled)
}
