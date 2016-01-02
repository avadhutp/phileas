package lib

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// NewService Provides the gonic gin service
func NewService(cfg *Cfg, db *gorm.DB, instaAPI *InstaAPI) *gin.Engine {
	api := NewPhileasAPI(cfg, db, instaAPI)

	s := gin.New()
	s.Use(gin.Logger())
	s.LoadHTMLGlob("templates/*")
	s.StaticFile("/favicon.ico", "static/favicon.ico")

	s.GET("/ping", api.ping)
	s.GET("/loc/:location-id", api.location)
	s.GET("/top", api.mapper)
	s.GET("/top.json", api.topJSON)

	return s
}
