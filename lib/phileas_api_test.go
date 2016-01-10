package lib

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func getSUT() *gin.Engine {
	cfg := &Cfg{}
	cfg.Common.GoogleMapsKey = "test-key"

	oldGinLoadHTMLGlob := ginLoadHTMLGlob
	defer func() {
		ginLoadHTMLGlob = oldGinLoadHTMLGlob
	}()
	ginLoadHTMLGlob = func(*gin.Engine, string) {}

	r := NewService(cfg, &gorm.DB{}, &InstaAPI{})
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
