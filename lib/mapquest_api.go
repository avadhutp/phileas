package lib

import (
	"fmt"

	"github.com/jasonwinn/geocoder"
	"github.com/jinzhu/gorm"
)

const (
	enrichmentLimit = 10
)

// ReverseGeocoder Goes systematically and enriches existing Location records with city + country information
type ReverseGeocoder struct {
	db *gorm.DB
}

// NewReverseGeocoder Provider for ReverseGeocoder
func NewReverseGeocoder(cfg *Cfg, db *gorm.DB) *ReverseGeocoder {
	geocoder.SetAPIKey(cfg.Common.MapquestKey)
	rg := new(ReverseGeocoder)
	rg.db = db

	return rg
}

// Enrich Periodically check the DB and enrich records for city + country info
func (rg *ReverseGeocoder) Enrich() {
	var locs []Location
	rg.db.Limit(enrichmentLimit).Where("city = ? and country = ?", "", "").Find(&locs)

	for _, loc := range locs {
		logger.Info(fmt.Sprintf("Lat: %f; Long: %f", loc.Lat, loc.Long))
		geo := reverseGeocode(loc.Lat, loc.Long)
		rg.updateLoc(geo, &loc)
	}
}

func (rg *ReverseGeocoder) updateLoc(geo *geocoder.Location, loc *Location) {
	if geo == nil || (geo.CountryCode == "" && geo.City == "") {
		return
	}

	loc.City = geo.City
	loc.Country = geo.CountryCode

	rg.db.Debug().Save(loc)
}

func reverseGeocode(lat float64, long float64) *geocoder.Location {
	if geo, err := geocoder.ReverseGeocode(lat, long); err != nil {
		return geo
	}

	return nil
}
