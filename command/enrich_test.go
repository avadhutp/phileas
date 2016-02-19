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
	oldESEnrichGooglePlacesIDs := esEnrichGooglePlaceIDs

	defer func() {
		libNewCfg = oldLibNewCfg
		libGetDB = oldLibGetDB
		libNewEnrichmentService = oldLibNewEnrichmentService
		esEnrichLocations = oldESEnrichLocations
		esEnrichGooglePlaceIDs = oldESEnrichGooglePlacesIDs
	}()

	libNewCfg = func(string) *lib.Cfg {
		return &lib.Cfg{}
	}

	libGetDB = func(*lib.Cfg) *gorm.DB {
		return nil
	}

	tests := []struct {
		arg                string
		enrichLocations    bool
		enrichGooglePlaces bool
		msg                string
	}{
		{"locations", true, false, "enrich is called with the argument locations; thus, only locations should be enriched"},
		{"google-places", false, true, "enrich is called with the argument google-places; thus, only Google Place IDs should be enriched"},
		{"unsupported-argument", false, false, "enrich is called with an unknown argument; therefore, nothing should be enriched"},
	}

	for _, test := range tests {
		enrichLocationsCalled := false
		esEnrichLocations = func(*lib.EnrichmentService) {
			enrichLocationsCalled = true
		}

		enrichGooglePlacesCalled := false
		esEnrichGooglePlaceIDs = func(*lib.EnrichmentService) {
			enrichGooglePlacesCalled = true
		}

		cmd := &cobra.Command{}
		args := []string{test.arg}
		enrich(cmd, args)

		assert.Equal(t, test.enrichLocations, enrichLocationsCalled, test.msg)
		assert.Equal(t, test.enrichGooglePlaces, enrichGooglePlacesCalled, test.msg)
	}
}
