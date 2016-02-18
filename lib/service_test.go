package lib

import (
	"reflect"
	"runtime"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	cfg := &Cfg{}
	cfg.Common.GoogleMapsKey = "test-key"
	api := NewPhileasAPI(cfg, &gorm.DB{}, &InstaAPI{})

	oldGinLoadHTMLGlob := ginLoadHTMLGlob
	defer func() {
		ginLoadHTMLGlob = oldGinLoadHTMLGlob
	}()

	ginLoadHTMLGlob = func(*gin.Engine, string) {}

	actual := NewService(cfg, &gorm.DB{}, &InstaAPI{})

	assert.Contains(t, actual.Routes(), gin.RouteInfo{
		Method:  "GET",
		Path:    "/ping",
		Handler: nameOfFunction(api.ping),
	})
	assert.Contains(t, actual.Routes(), gin.RouteInfo{
		Method:  "GET",
		Path:    "/loc/:location-id",
		Handler: nameOfFunction(api.location),
	})
	assert.Contains(t, actual.Routes(), gin.RouteInfo{
		Method:  "GET",
		Path:    "/top",
		Handler: nameOfFunction(api.mapper),
	})
	assert.Contains(t, actual.Routes(), gin.RouteInfo{
		Method:  "GET",
		Path:    "/top.json",
		Handler: nameOfFunction(api.topJSON),
	})
	assert.Contains(t, actual.Routes(), gin.RouteInfo{
		Method:  "GET",
		Path:    "/countries.json",
		Handler: nameOfFunction(api.countriesJSON),
	})
}

func nameOfFunction(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
