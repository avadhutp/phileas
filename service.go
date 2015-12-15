package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getService() *gin.Engine {
	s := gin.New()
	s.Use(gin.Logger())

	s.GET("/ping", ping)

	return s
}

func ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
