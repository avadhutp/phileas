package lib

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
	"strings"
	"testing"
	"time"

	"github.com/avadhutp/phileas/helpers"
	testdb "github.com/erikstmartin/go-testdb"
	"github.com/jasonwinn/geocoder"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

var (
	sut *EnrichmentService
)

func init() {
	db, _ := gorm.Open("testdb", "")
	cfg := &Cfg{}

	sut = NewEnrichmentService(cfg, db)
}

func stubQuery(r string) {
	sql := []string{
		`SELECT  * FROM "locations"  WHERE (city = ? and country = ?) LIMIT 10`,
		`SELECT  * FROM "locations"  WHERE (google_places_id IS NULL) LIMIT 10`,
	}

	for _, q := range sql {
		testdb.StubQuery(q, testdb.RowsFromCSVString(locationCols, r))
	}
}

func resetSUT() {
	sut.waits = map[int]time.Duration{
		typeLoc:          waitBetweenEnrichment,
		typeGooglePlaces: waitBetweenEnrichment,
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
	assert.Equal(t, map[int]time.Duration{typeLoc: waitBetweenEnrichment, typeGooglePlaces: waitBetweenEnrichment}, actual.waits)
}

func TestEnrichGooglePlacesAllDone(t *testing.T) {
	resetSUT()
	testdb.Reset()
	stubQuery(``)

	oldEsEnrichGooglePlacesIDs := esEnrichGooglePlacesIDs
	oldTimeSleep := timeSleep

	defer func() {
		timeSleep = oldTimeSleep
		esEnrichGooglePlacesIDs = oldEsEnrichGooglePlacesIDs
	}()

	var sleepFor time.Duration
	timeSleep = func(s time.Duration) {
		sleepFor = s
	}

	esEnrichGooglePlacesIDsCalled := false
	esEnrichGooglePlacesIDs = func(*EnrichmentService) {
		esEnrichGooglePlacesIDsCalled = true
	}

	sut.EnrichGooglePlacesIDs()

	assert.Equal(t, waitBetweenEnrichment*2, sleepFor)
	assert.Equal(t, waitBetweenEnrichment*2, sut.waits[typeGooglePlaces])
	assert.True(t, esEnrichGooglePlacesIDsCalled)
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

func TestEnrichGooglePlacesIDs(t *testing.T) {
	tests := []struct {
		placeID string
		msg     string
	}{
		{"", "No result was found; therefore, location should be saved with an empty google places ID"},
		{"test-place-id", "Place result was found; therefore, location should be saved with the resultant google places ID"},
	}

	for _, test := range tests {

		resetSUT()
		testdb.Reset()

		insertCalled := false
		testdb.SetExecWithArgsFunc(func(q string, args []driver.Value) (result driver.Result, err error) {
			if strings.Contains(q, `UPDATE "locations"`) && argInSlice(int64(123), args) {
				insertCalled = true
			}
			if len(test.placeID) > 0 && !argInSlice(string(test.placeID), args) {
				insertCalled = false
			}

			return helpers.NewTestResult(1, 0), nil
		})

		r := `
		123, location-1, 1.0, 1.0, test address, UK, London,
		`
		stubQuery(r)
		db, _ := gorm.Open("testdb", "")
		sut.db = db

		oldEsEnrichGooglePlacesIDs := esEnrichGooglePlacesIDs
		oldTimeSleep := timeSleep
		oldRadarSearch := radarSearch

		defer func() {
			esEnrichGooglePlacesIDs = oldEsEnrichGooglePlacesIDs
			timeSleep = oldTimeSleep
			radarSearch = oldRadarSearch
		}()

		timeSleep = func(s time.Duration) {}
		esEnrichGooglePlacesIDs = func(*EnrichmentService) {}

		radarSearch = func(client *maps.Client, c context.Context, req *maps.RadarSearchRequest) (maps.PlacesSearchResponse, error) {
			results := []maps.PlacesSearchResult{}
			if len(test.placeID) > 0 {
				place := maps.PlacesSearchResult{}
				place.PlaceID = test.placeID
				results = append(results, place)
			}

			resp := maps.PlacesSearchResponse{}
			resp.Results = results

			return resp, nil
		}

		sut.EnrichGooglePlacesIDs()

		assert.True(t, insertCalled, test.msg)
	}
}

func TestEnrichLocation(t *testing.T) {
	tests := []struct {
		reverseGeocodeError error
		country             string
		shouldCallInsert    bool
		msg                 string
	}{
		{nil, "UK", true, "All ok, so DB insert should be called"},
		{nil, "", false, "Reverse geocoding did not return a country, so DB insert should not be called"},
		{errors.New("Test error"), "", false, "Reverse geocoding failed with an error, so DB insert should not be called"},
	}

	for _, test := range tests {
		resetSUT()
		testdb.Reset()

		insertCalled := false
		testdb.SetExecWithArgsFunc(func(q string, args []driver.Value) (result driver.Result, err error) {
			if strings.Contains(q, `INSERT INTO "locations"`) {
				insertCalled = true
			}

			return helpers.NewTestResult(1, 0), nil
		})

		r := `
		location-1, 1.0, 1.0, test address, UK, London
		`
		stubQuery(r)

		db, _ := gorm.Open("testdb", "")
		sut.db = db

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
			g.City = "London"
			g.County = "London"
			return g, test.reverseGeocodeError
		}

		sut.EnrichLocation()

		assert.Equal(t, test.shouldCallInsert, insertCalled, test.msg)
	}
}

func TestThrottleWait(t *testing.T) {
	tests := []struct {
		found                int
		waitType             int
		postThrottleInterval time.Duration
	}{
		{10, typeLoc, waitBetweenEnrichment},
		{9, typeLoc, waitBetweenEnrichment * 2},
	}

	for _, test := range tests {
		resetSUT()
		sut.throttleWait(test.found, test.waitType)

		assert.Equal(t, test.postThrottleInterval, sut.waits[test.waitType])
	}
}

func TestCopyGeoToLoc(t *testing.T) {
	tests := []struct {
		expectedCity string
		input        geocoder.Location
		msg          string
	}{
		{"", geocoder.Location{"", "", "", "", "", "", geocoder.LatLng{}, "", false}, "Empty geocoder.Location"},
		{"London", geocoder.Location{"", "London", "", "", "", "", geocoder.LatLng{}, "", false}, "City in geocoder.Location.City field"},
		{"London", geocoder.Location{"", "", "", "", "London", "", geocoder.LatLng{}, "", false}, "City in geocoder.Location.County field"},
		{"London", geocoder.Location{"", "", "London", "", "", "", geocoder.LatLng{}, "", false}, "City in geocoder.Location.State field"},
	}

	for _, test := range tests {
		loc := &Location{}
		loc.City = ""

		copyGeoToLoc(loc, &test.input)

		assert.Equal(t, test.expectedCity, loc.City, test.msg)
	}
}

func argInSlice(a driver.Value, list []driver.Value) bool {
	for _, b := range list {
		fmt.Println(a, b)
		if b == a {
			return true
		}
	}
	return false
}
