package lib

import (
	"fmt"
	"time"

	"github.com/JustinBeckwith/go-yelp/yelp"

	"github.com/jasonwinn/geocoder"
	"github.com/jinzhu/gorm"
)

const (
	enrichmentLimit       = 10
	waitBetweenEnrichment = 30 * time.Second
)

// EnrichmentService Goes systematically and enriches existing Location records with city + country information
type EnrichmentService struct {
	db         *gorm.DB
	yelpClient *yelp.Client
	wait       time.Duration
}

// NewEnrichmentService Provider for EnrichmentService
func NewEnrichmentService(cfg *Cfg, db *gorm.DB) *EnrichmentService {
	geocoder.SetAPIKey(cfg.Common.MapquestKey)

	rg := new(EnrichmentService)
	rg.db = db
	rg.wait = waitBetweenEnrichment

	auth := yelp.AuthOptions{
		ConsumerKey:       cfg.Yelp.ConsumerKey,
		ConsumerSecret:    cfg.Yelp.ConsumerSecret,
		AccessToken:       cfg.Yelp.AccessToken,
		AccessTokenSecret: cfg.Yelp.AccessTokenSecret,
	}
	rg.yelpClient = yelp.New(auth, nil)

	return rg
}

// Enrich Periodically check the DB and enrich records for city + country info
func (rg *EnrichmentService) EnrichLocation() {
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

func (rg *EnrichmentService) updateLoc(geo *geocoder.Location, loc *Location) {
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

func (rg *EnrichmentService) throttleWait(found int) {
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
