package lib

import (
	"github.com/jasonwinn/geocoder"
	"github.com/jinzhu/gorm"
)

// ReverseGeocoder Goes systematically and enriches existing Location records with city + country information
type ReverseGeocoder struct {
	db *gorm.DB
}

// NewReverseGeocoder Provider for ReverseGeocoder
func NewReverseGeocoder(cfg *Cfg, db *gorm.DB) *ReverseGeocoder {
	rg := new(ReverseGeocoder)
	rg.db = db

	return rg
}

func reverseGeocode(lat float64, long float64) *geocoder.Location {
	if loc, err := geocoder.ReverseGeocode(lat, long); err != nil {
		return loc
	}

	return nil
}
