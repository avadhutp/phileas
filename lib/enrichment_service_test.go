package lib

import (
	"testing"
	"time"

	_ "github.com/erikstmartin/go-testdb"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

var (
	es *EnrichmentService
)

func init() {
	db, _ := gorm.Open("testdb", "")
	cfg := &Cfg{}

	es = NewEnrichmentService(cfg, &db)
}

func TestNewEnrichmentService(t *testing.T) {
	cfg := &Cfg{}
	common := &common{}
	common.MapquestKey = "test-key"
	cfg.Common = *common

	db := &gorm.DB{}

	oldGeocoderSetAPIKey := geocoderSetAPIKey
	defer func() { geocoderSetAPIKey = oldGeocoderSetAPIKey }()

	geocoderSetAPIKeyCalled := false
	geocoderSetAPIKey = func(key string) {
		if key == "test-key" {
			geocoderSetAPIKeyCalled = true
		}
	}

	actual := NewEnrichmentService(cfg, db)

	assert.Equal(t, db, actual.db)
	assert.True(t, geocoderSetAPIKeyCalled)
	assert.Equal(t, map[int]time.Duration{typeLoc: waitBetweenEnrichment}, actual.waits)
}
