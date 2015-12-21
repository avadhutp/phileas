package lib

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NewService Provides the gonic gin service
func NewService(cfg *Cfg) *gin.Engine {
	s := gin.New()
	s.Use(gin.Logger())
	s.LoadHTMLGlob("templates/*")

	s.GET("/ping", ping)

	s.GET("/top", func(c *gin.Context) {
		c.HTML(http.StatusOK, "mapper.tmpl", gin.H{
			"title": "Top destinations",
			"key":   cfg.Common.GoogleMapsKey,
		})
	})

	return s
}

func ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
