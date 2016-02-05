package lib

import (
	"database/sql/driver"
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

func stubInsertLocation() {
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
}

func TestBackfill(t *testing.T) {
	oldGetLikedMedia := getLikedMedia
	oldInstaAPISaveMedia := instaAPISaveMedia

	defer func() {
		getLikedMedia = oldGetLikedMedia
		instaAPISaveMedia = oldInstaAPISaveMedia
	}()

	saveMediaCallCnt := 0
	instaAPISaveMedia = func(*InstaAPI, *instagram.Media) {
		saveMediaCallCnt += 1
	}

	getLikedMedia = func(*instagram.UsersService, *instagram.Parameters) ([]instagram.Media, *instagram.ResponsePagination, error) {
		r := &instagram.ResponsePagination{
			NextURL:   "http://some.random.url",
			NextMaxID: "",
		}

		m := []instagram.Media{
			*(&instagram.Media{}),
			*(&instagram.Media{}),
		}

		return m, r, nil
	}

	instaAPI.Backfill("")

	assert.Equal(t, 2, saveMediaCallCnt)
}

func TestSaveMedia(t *testing.T) {
	testdb.Reset()

	caption := &instagram.MediaCaption{}
	caption.Text = "Test caption"

	thumbnail := &instagram.MediaImage{}
	thumbnail.URL = "http://test/url"

	images := &instagram.MediaImages{}
	images.Thumbnail = thumbnail

	m := &instagram.Media{}
	m.Location = &instagram.MediaLocation{}
	m.ID = "test-id"
	m.Images = images
	m.Link = "http://full.image/url"
	m.Caption = caption
	m.CreatedTime = 12345678

	oldInstaAPISaveLocation := instaAPISaveLocation
	defer func() { instaAPISaveLocation = oldInstaAPISaveLocation }()

	isSaveLocationCalled := false
	instaAPISaveLocation = func(*InstaAPI, *instagram.Media) *Location {
		isSaveLocationCalled = true
		loc := &Location{}
		loc.ID = 123

		return loc
	}

	isDBQueryCalled := false
	testdb.SetQueryWithArgsFunc(func(query string, args []driver.Value) (driver.Rows, error) {
		isDBQueryCalled = true
		return testdb.RowsFromCSVString(entryCols, ``), nil
	})

	db, _ := gorm.Open("testdb", "")
	instaAPI.db = &db
	instaAPI.saveMedia(m)

	assert.True(t, isSaveLocationCalled)
	assert.True(t, isDBQueryCalled)
}

func TestSaveMediaWithoutLocation(t *testing.T) {
	m := &instagram.Media{}
	m.Location = nil

	oldInstaAPISaveLocation := instaAPISaveLocation
	defer func() { instaAPISaveLocation = oldInstaAPISaveLocation }()

	isSaveLocationCalled := false
	instaAPISaveLocation = func(*InstaAPI, *instagram.Media) *Location {
		isSaveLocationCalled = true
		return nil
	}

	instaAPI.saveMedia(m)

	assert.False(t, isSaveLocationCalled, "Media's location is nil; this implies that the instagram post was not geotagged. Therefore, it should not be saved.")
}

func TestSaveLocation(t *testing.T) {
	testdb.Reset()
	stubInsertLocation()

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

	db, _ := gorm.Open("testdb", "")
	instaAPI.db = &db
	actual := instaAPI.saveLocation(m)

	assert.Equal(t, expected, actual)
}
