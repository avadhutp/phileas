package lib

import (
	"github.com/jinzhu/gorm"
)

type APIProvider struct {
	db *gorm.DB
}

func (api *APIProvider) top() []Location {
	var locs []Location
	api.db.Find(&locs)

	return locs
}
