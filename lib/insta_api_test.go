package lib

import (
	_ "github.com/erikstmartin/go-testdb"
	"github.com/gedex/go-instagram/instagram"
	"github.com/jinzhu/gorm"
	"testing"
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
	m := &instagram.Media{}
	m.Location = &instagram.MediaLocation{
		ID:        1,
		Name:      "test-location",
		Latitude:  1.0,
		Longitude: 1.0,
	}
}
