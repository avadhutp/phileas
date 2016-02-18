package lib

import (
	"database/sql/driver"
	"strings"
	"testing"
	"time"

	"github.com/avadhutp/phileas/vendor"

	"github.com/jasonwinn/geocoder"

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
	sql := []string{
		`SELECT  * FROM "locations"  WHERE (city = ? and country = ?) LIMIT 10`,
	}

	for _, q := range sql {
		testdb.StubQuery(q, testdb.RowsFromCSVString(locationCols, r))
	}
}

func resetSUT() {
	sut.waits = map[int]time.Duration{
		typeLoc: waitBetweenEnrichment,
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
	resetSUT()
	testdb.Reset()
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

func TestEnrichLocationNoGeo(t *testing.T) {
	tests := []struct {
		country          string
		shouldCallInsert bool
		msg              string
	}{
		{country: "UK", shouldCallInsert: true, msg: "All ok, so DB insert should be called"},
	}

	for _, test := range tests {
		resetSUT()
		testdb.Reset()

		insertCalled := false
		testdb.SetExecWithArgsFunc(func(q string, args []driver.Value) (result driver.Result, err error) {
			if strings.Contains(q, `INSERT INTO "locations"`) {
				insertCalled = true
			}

			return vendor.NewTestResult(1, 0), nil
		})

		r := `
		location-1, 1.0, 1.0, test address, UK, London, ,
		`
		stubQuery(r)

		db, _ := gorm.Open("testdb", "")
		sut.db = &db

		oldEsEnrichLocation := esEnrichLocation
		oldTimeSleep := timeSleep
		oldGeocoderReverseGeocode := geocoderReverseGeocode

		defer func() {
			timeSleep = oldTimeSleep
			esEnrichLocation = oldEsEnrichLocation
			geocoderReverseGeocode = oldGeocoderReverseGeocode
		}()

		timeSleep = func(s time.Duration) {}
		esEnrichLocation = func(*EnrichmentService) {}

		geocoderReverseGeocode = func(float64, float64) (*geocoder.Location, error) {
			g := &geocoder.Location{}
			g.CountryCode = test.country
			return g, nil
		}

		sut.EnrichLocation()

		assert.Equal(t, test.shouldCallInsert, insertCalled, test.msg)
	}
}
