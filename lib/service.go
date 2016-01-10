package lib

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var (
	ginLoadHTMLGlob = (*gin.Engine).LoadHTMLGlob
)

// NewService Provides the gonic gin service
func NewService(cfg *Cfg, db *gorm.DB, instaAPI *InstaAPI) *gin.Engine {
	api := NewPhileasAPI(cfg, db, instaAPI)

	s := gin.New()
	s.Use(gin.Logger())
	ginLoadHTMLGlob(s, "templates/*")
	s.Static("/static", "static")
	s.StaticFile("/favicon.ico", "static/favicon.ico")

	s.GET("/ping", api.ping)
	s.GET("/loc/:location-id", api.location)
	s.GET("/top", api.mapper)
	s.GET("/top.json", api.topJSON)

	return s
}
