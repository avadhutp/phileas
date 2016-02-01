package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"

	testdb "github.com/erikstmartin/go-testdb"
	"github.com/gedex/go-instagram/instagram"
	"github.com/jinzhu/gorm"
)

var (
	instaAPI *InstaAPI
)

func init() {
	cfg := &Cfg{}
	cfg.Instagram.ClientID = "test-client-id"
	cfg.Instagram.Secret = "test-secret"
	cfg.Instagram.Token = "test-token"

	db, _ := gorm.Open("testdb", "")

	instaAPI = NewInstaAPI(cfg, &db)
}
func TestSaveLocation(t *testing.T) {
	sql := []string{
		`SELECT  * FROM "locations"  WHERE ("long" = ?) AND ("name" = ?) AND ("lat" = ?) ORDER BY "locations"."id" ASC LIMIT 1`,
		`SELECT  * FROM "locations"  WHERE ("long" = ?) AND ("lat" = ?) AND ("name" = ?) ORDER BY "locations"."id" ASC LIMIT 1`,
		`SELECT  * FROM "locations"  WHERE ("name" = ?) AND ("lat" = ?) AND ("long" = ?) ORDER BY "locations"."id" ASC LIMIT 1`,
		`SELECT  * FROM "locations"  WHERE ("name" = ?) AND ("long" = ?) AND ("lat" = ?) ORDER BY "locations"."id" ASC LIMIT 1`,
		`SELECT  * FROM "locations"  WHERE ("lat" = ?) AND ("name" = ?) AND ("long" = ?) ORDER BY "locations"."id" ASC LIMIT 1`,
		`SELECT  * FROM "locations"  WHERE ("lat" = ?) AND ("long" = ?) AND ("name" = ?) ORDER BY "locations"."id" ASC LIMIT 1`,
	}
	result := `
	1, location-1, 1.0, 1.0, test address, UK, London, ,
	`

	for _, q := range sql {
		testdb.StubQuery(q, testdb.RowsFromCSVString(locationCols, result))
	}

	expected := &Location{
		ID:       1,
		Name:     "location-1",
		Lat:      1.0,
		Long:     1.0,
		Address:  "test address",
		Country:  "UK",
		City:     "London",
		YelpType: "",
		YelpURL:  "",
	}

	m := &instagram.Media{}
	m.Location = &instagram.MediaLocation{
		ID:        1,
		Name:      "location-1",
		Latitude:  1.0,
		Longitude: 1.0,
	}

	actual := instaAPI.saveLocation(m)

	assert.Equal(t, expected, actual)
}
