package lib

import (
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

	return NewService(cfg, &gorm.DB{}, &InstaAPI{})
}

func TestPing(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	sut := getSUT()
	sut.ServeHTTP(w, req)

	assert.Equal(t, "pong", w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestMapper(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/top", nil)
	sut := getSUT()
	sut.ServeHTTP(w, req)
}
