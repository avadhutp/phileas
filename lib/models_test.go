package lib

import (
	"errors"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func getCfg() *Cfg {
	cfg := &Cfg{}
	cfg.Mysql.Username = "user"
	cfg.Mysql.Password = "pass"
	cfg.Mysql.Host = "localhost"
	cfg.Mysql.Port = "3306"
	cfg.Mysql.Database = "test"

	return cfg
}

func TestGetDBConnString(t *testing.T) {
	expected := "user:pass@tcp(localhost:3306)/test?charset=utf8mb4&character_set_server=utf8mb4&parseTime=True&loc=Local"
	actual := getDBConnString(getCfg())

	assert.Equal(t, expected, actual)
}

func TestGetDB(t *testing.T) {
	oldGormOpen := gormOpen
	oldDBSingularTable := dbSingularTable

	defer func() {
		gormOpen = oldGormOpen
		dbSingularTable = oldDBSingularTable
	}()

	mockDB := gorm.DB{}
	gormOpen = func(string, ...interface{}) (gorm.DB, error) {
		return mockDB, nil
	}

	tableSingularity := false
	dbSingularTable = func(db *gorm.DB, flag bool) {
		tableSingularity = flag
	}

	actual := GetDB(getCfg())

	assert.Equal(t, &mockDB, actual)
	assert.True(t, tableSingularity)
}

func TestGetDBErrorHandling(t *testing.T) {
	oldGormOpen := gormOpen
	oldDBSingularTable := dbSingularTable

	defer func() {
		gormOpen = oldGormOpen
		dbSingularTable = oldDBSingularTable
	}()

	gormOpen = func(string, ...interface{}) (gorm.DB, error) {
		return gorm.DB{}, errors.New("Test func")
	}

	dbSingularTable = func(db *gorm.DB, flag bool) {}

	testFunc := assert.PanicTestFunc(func() {
		GetDB(getCfg())
	})

	assert.Panics(t, testFunc, "gorm.open returns an error, so our code should panic")
}
