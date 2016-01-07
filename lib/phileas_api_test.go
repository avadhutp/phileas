package lib

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
)

func TestPing(t *testing.T) {
	oldGinContextString := ginContextString
	defer func() { ginContextString = oldGinContextString }()

	c := &gin.Context{}
	sut := &PhileasAPI{}

	var status int
	var body string
	ginContextString = func(pe *gin.Context, s int, b string, args ...interface{}) {
		status = s
		body = b
	}

	sut.ping(c)

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "pong", body)
}
