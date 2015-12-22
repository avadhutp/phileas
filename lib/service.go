package lib

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// NewService Provides the gonic gin service
func NewService(cfg *Cfg, db *gorm.DB) *gin.Engine {
	api := NewPhileasAPI(cfg, db)

	s := gin.New()
	s.Use(gin.Logger())
	s.LoadHTMLGlob("templates/*")

	s.GET("/ping", api.ping)
	s.GET("/top", api.mapper)
	s.GET("/top.json", api.topJSON)

	return s
}
