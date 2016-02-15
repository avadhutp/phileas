package lib

import (
	"testing"
	"time"

	testdb "github.com/erikstmartin/go-testdb"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

var (
	sut *EnrichmentService
)

func init() {
	db, _ := gorm.Open("testdb", "")
	cfg := &Cfg{}

	sut = NewEnrichmentService(cfg, &db)
}

func stubQuery(r string) {
	testdb.Reset()

	sql := []string{
		`SELECT  * FROM "locations"  WHERE (city = ? and country = ?) LIMIT 10`,
	}

	for _, q := range sql {
		testdb.StubQuery(q, testdb.RowsFromCSVString(locationCols, r))
	}
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

func TestEnrichLocationAllDone(t *testing.T) {
	stubQuery(``)

	oldEsEnrichLocation := esEnrichLocation
	oldTimeSleep := timeSleep

	defer func() {
		timeSleep = oldTimeSleep
		esEnrichLocation = oldEsEnrichLocation
	}()

	var sleepFor time.Duration
	timeSleep = func(s time.Duration) {
		sleepFor = s
	}

	esEnrichLocationCalled := false
	esEnrichLocation = func(*EnrichmentService) {
		esEnrichLocationCalled = true
	}

	sut.EnrichLocation()

	assert.Equal(t, waitBetweenEnrichment*2, sleepFor)
	assert.Equal(t, waitBetweenEnrichment*2, sut.waits[typeLoc])
	assert.True(t, esEnrichLocationCalled)
}
