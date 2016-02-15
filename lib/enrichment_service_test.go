package lib

import (
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestNewEnrichmentService(t *testing.T) {
	cfg := &Cfg{}
	common := &common{}
	common.MapquestKey = "test-key"
	cfg.Common = *common

	db := &gorm.DB{}

	actual := NewEnrichmentService(cfg, db)

	assert.Equal(t, db, actual.db)
}
