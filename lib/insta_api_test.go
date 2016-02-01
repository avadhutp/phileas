package lib

import (
	"fmt"
	"testing"

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
	sql := `SELECT  * FROM "locations"  WHERE ("long" = ?) AND ("name" = ?) AND ("lat" = ?) ORDER BY "locations"."id" ASC LIMIT 1`
	result := `
	1, location-1, 1.0, 1.0, test address, UK, London, NA, NA		
	`
	testdb.StubQuery(sql, testdb.RowsFromCSVString(locationCols, result))

	m := &instagram.Media{}
	m.Location = &instagram.MediaLocation{
		ID:        1,
		Name:      "test-location",
		Latitude:  1.0,
		Longitude: 1.0,
	}

	x := instaAPI.saveLocation(m)

	fmt.Println(x)
}
