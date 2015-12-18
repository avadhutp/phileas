package lib

import (
	"fmt"
	"time"

	"github.com/jasonwinn/geocoder"
	"github.com/jinzhu/gorm"
)

const (
	enrichmentLimit       = 10
	waitBetweenEnrichment = 30 * time.Second
)

// ReverseGeocoder Goes systematically and enriches existing Location records with city + country information
type ReverseGeocoder struct {
	db   *gorm.DB
	wait time.Duration
}

// NewReverseGeocoder Provider for ReverseGeocoder
func NewReverseGeocoder(cfg *Cfg, db *gorm.DB) *ReverseGeocoder {
	geocoder.SetAPIKey(cfg.Common.MapquestKey)

	rg := new(ReverseGeocoder)
	rg.db = db
	rg.wait = waitBetweenEnrichment

	return rg
}

// Enrich Periodically check the DB and enrich records for city + country info
func (rg *ReverseGeocoder) Enrich() {
	for {
		var locs []Location
		rg.db.Limit(enrichmentLimit).Where("city = ? and country = ?", "", "").Find(&locs)

		for _, loc := range locs {
			geo := reverseGeocode(loc.Lat, loc.Long)
			rg.updateLoc(geo, &loc)
		}

		rg.throttleWait(len(locs))
		time.Sleep(rg.wait)
	}
}

func (rg *ReverseGeocoder) updateLoc(geo *geocoder.Location, loc *Location) {
	if geo == nil || geo.CountryCode == "" {
		return
	}

	loc.Country = geo.CountryCode

	if geo.City != "" {
		loc.City = geo.City
	} else if geo.County != "" {
		loc.City = geo.County
	} else if geo.State != "" {
		loc.City = geo.State
	}

	rg.db.Save(loc)
}

func (rg *ReverseGeocoder) throttleWait(found int) {
	if found < enrichmentLimit {
		rg.wait += waitBetweenEnrichment
	} else {
		rg.wait = waitBetweenEnrichment
	}
}

func reverseGeocode(lat float64, long float64) *geocoder.Location {
	if geo, err := geocoder.ReverseGeocode(lat, long); err == nil {
		return geo
	} else {
		logger.Error(fmt.Sprintf("Reverse geocoding encountered and error: %s", err.Error()))
	}

	return nil
}
