package lib

import (
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	testdb "github.com/erikstmartin/go-testdb"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

var (
	testDB       = &gorm.DB{}
	locationCols = []string{"id", "name", "lat", "long", "address", "country", "city", "yelptype", "yelpurl"}
)

func init() {
	tables := []interface{}{&Entry{}, &Location{}}

	db, _ := gorm.Open("testdb", "")
	db.AutoMigrate(tables...)

	testDB = &db
}

func getSUT() *gin.Engine {
	cfg := &Cfg{}
	cfg.Common.GoogleMapsKey = "test-key"

	ginLoadHTMLGlob = func(*gin.Engine, string) {}

	r := NewService(cfg, testDB, &InstaAPI{})

	tmpl, _ := template.New("mapper.tmpl").Parse(`{{ .title }} | {{ .key }}`)
	r.SetHTMLTemplate(tmpl)

	return r
}

func peformRequest(method string, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, nil)
	sut := getSUT()
	sut.ServeHTTP(w, req)

	return w
}

func TestPing(t *testing.T) {
	w := peformRequest("GET", "/ping")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestMapper(t *testing.T) {
	w := peformRequest("GET", "/top")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Top destinations | test-key", w.Body.String())
}

func TestTopJSON(t *testing.T) {
	sql := `SELECT * from "locations"`
	result := `
	1, location-1, 1.0, 1.0, test address, UK, London, NA, NA
	2, location-2, 1.0, 1.0, test address, IN, Mumbai, NA, NA
	`
	testdb.StubQuery(sql, testdb.RowsFromCSVString(locationCols, result))
	expected := `{"type":"FeatureCollection","features":[{"type":"Feature","geometry":{"type":"Point","coordinates":[1,1]},"properties":{"id":1,"name":"location-1"}},{"type":"Feature","geometry":{"type":"Point","coordinates":[1,1]},"properties":{"id":2,"name":"location-2"}}]}`

	w := peformRequest("GET", "/top.json")

	assert.Equal(t, fmt.Sprintf("%s\n", expected), w.Body.String())
}
