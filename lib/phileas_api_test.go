package lib

import (
	"database/sql/driver"
	"errors"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/jasonwinn/geocoder"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/json"

	testdb "github.com/erikstmartin/go-testdb"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

var (
	service *gin.Engine
	minHTML = minify.New()
	minJSON = minify.New()
)

func init() {
	minHTML.AddFunc("text/html", html.Minify)
	minJSON.AddFunc("text/json", json.Minify)
}

func init() {
	db, _ := gorm.Open("testdb", "")

	cfg := &Cfg{}
	cfg.Google.BrowserKey = "test-key"

	ginLoadHTMLGlob = func(*gin.Engine, string) {}
	service = NewService(cfg, db, &InstaAPI{})
}

func performRequest(method string, path string) *httptest.ResponseRecorder {
	return performRequestWithService(service, method, path)
}

func performRequestWithService(s *gin.Engine, method string, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, nil)
	s.ServeHTTP(w, req)

	return w
}

func TestPing(t *testing.T) {
	w := performRequest("GET", "/ping")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestMapper(t *testing.T) {
	tmpl, _ := template.New("mapper.tmpl").Parse(`{{ .title }} | {{ .key }}`)
	service.SetHTMLTemplate(tmpl)

	w := performRequest("GET", "/top")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Top destinations | test-key", w.Body.String())
}

func TestStatsJSON(t *testing.T) {
	countCols := []string{"count"}
	sqlStubs := []struct {
		query  string
		output int
	}{
		{`SELECT  count(*) FROM "locations"`, 10},
		{`SELECT  count(*) FROM "locations"  WHERE (google_places_id IS NOT NULL and google_places_id != ?)`, 20},
		{`SELECT  count(*) FROM "locations"  WHERE (google_places_id IS NULL)`, 30},
		{`SELECT  count(*) FROM "locations"  WHERE (google_places_id = ?)`, 40},
	}

	for _, stub := range sqlStubs {
		result := `
		` + strconv.Itoa(stub.output) + `
		`
		testdb.StubQuery(stub.query, testdb.RowsFromCSVString(countCols, result))
	}

	expected := `{"enrichment":{"google_places":{"enriched":20,"un_enriched":30,"without_place_id":40}},"total_locations":10}`
	w := performRequest("GET", "/stats.json")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, minifyJSON(expected), minifyJSON(w.Body.String()))
}

func TestCountriesJSON(t *testing.T) {
	selectSQL := "SELECT  `id`, `country`, count(*) FROM \"location\"   GROUP BY country HAVING (`country` != '')"
	selectCols := []string{"id", "country", "count"}
	selectResult := `
	1, UK, 5
	`
	testdb.StubQuery(selectSQL, testdb.RowsFromCSVString(selectCols, selectResult))

	countSQL := `SELECT  count(*) FROM "locations"  WHERE (country != ?)`
	countCols := []string{"count"}
	countResult := `
	5
	`
	testdb.StubQuery(countSQL, testdb.RowsFromCSVString(countCols, countResult))

	expected := `{"type":"FeatureCollection","features":[{"type":"Feature","geometry":{"type":"Point","coordinates":[0,0]},"properties":{"country":"","id":1,"size":100,"total":5}}]}`

	w := performRequest("GET", "/countries.json")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, minifyJSON(expected), minifyJSON(w.Body.String()))
}

func TestCountriesJSONErrorHandling(t *testing.T) {
	testdb.Reset()
	testdb.SetQueryFunc(func(q string) (result driver.Rows, err error) {
		return nil, errors.New("Test exception")
	})

	db, _ := gorm.Open("testdb", "")

	cfg := &Cfg{}
	cfg.Google.BrowserKey = "test-key"
	service := NewService(cfg, db, &InstaAPI{})

	w := performRequestWithService(service, "GET", "/countries.json")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Test exception", w.Body.String())
}

func TestTopJSON(t *testing.T) {
	sql := `SELECT * from "locations"`
	result := `
	1, location-1, 1.0, 1.0, test address, UK, London
	2, location-2, 1.0, 1.0, test address, IN, Mumbai
	`
	testdb.StubQuery(sql, testdb.RowsFromCSVString(locationCols, result))
	expected := `{"type":"FeatureCollection","features":[{"type":"Feature","geometry":{"type":"Point","coordinates":[1,1]},"properties":{"google_place_id":"","id":1,"name":"location-1"}},{"type":"Feature","geometry":{"type":"Point","coordinates":[1,1]},"properties":{"google_place_id":"","id":2,"name":"location-2"}}]}`

	w := performRequest("GET", "/top.json")

	assert.Equal(t, minifyJSON(expected), minifyJSON(w.Body.String()))
}

func TestLocation(t *testing.T) {
	sql := `SELECT * FROM "entries"  WHERE ("location_id" = ?)`
	result := `
	1, instagram, 12345, http://thumbnail/url, http://full/url, test caption, 12345678, 1
	`
	expected := `
	<div class="row">
		<div class="left"><a href="http://full/url" target="_blank"><img src="http://thumbnail/url" /></a></div>
        <div class="right">test caption</div>
    </div>
	`
	testdb.StubQuery(sql, testdb.RowsFromCSVString(entryCols, result))
	w := performRequest("GET", "/loc/1")

	assert.Equal(t, minifyHTML(expected), minifyHTML(w.Body.String()))
}

func TestCacheCountryLatLong(t *testing.T) {
	cache := cacheCountryLatLong("../static/country_latlon.csv")

	expected := CountryInfo{
		name: "United Kingdom",
		geo: geocoder.LatLng{
			Lat: 55.378051,
			Lng: -3.435973,
		},
	}
	assert.Equal(t, expected, cache["GB"])
}

func minifyHTML(raw string) string {
	out, _ := minHTML.String("text/html", raw)

	return out
}

func minifyJSON(raw string) string {
	out, _ := minJSON.String("text/json", raw)

	return out
}
