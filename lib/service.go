package lib

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NewService Provides the gonic gin service
func NewService(cfg *Cfg) *gin.Engine {
	eng := &PhileasEngine{cfg}

	s := gin.New()
	s.Use(gin.Logger())
	s.LoadHTMLGlob("templates/*")

	s.GET("/ping", eng.ping)
	s.GET("/top", eng.mapper)

	return s
}

type PhileasEngine struct {
	cfg *Cfg
}

func (pe *PhileasEngine) ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func (pe *PhileasEngine) mapper(c *gin.Context) {
	c.HTML(http.StatusOK, "mapper.tmpl", gin.H{
		"title": "Top destinations",
		"key":   pe.cfg.Common.GoogleMapsKey,
	})
}
