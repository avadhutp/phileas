package lib

import (
	"fmt"
	"strings"
	"time"

	"github.com/JustinBeckwith/go-yelp/yelp"
	"github.com/guregu/null"
	"github.com/jasonwinn/geocoder"
	"github.com/jinzhu/gorm"
)

const (
	enrichmentLimit       = 10
	waitBetweenEnrichment = 30 * time.Second
)

// EnrichmentService Goes systematically and enriches existing Location records with city + country information
type EnrichmentService struct {
	db            *gorm.DB
	yelpClient    *yelp.Client
	locWait       time.Duration
	yelpWait      time.Duration
	sanitizeRegex strings.Replacer
}

// NewEnrichmentService Provider for EnrichmentService
func NewEnrichmentService(cfg *Cfg, db *gorm.DB) *EnrichmentService {
	geocoder.SetAPIKey(cfg.Common.MapquestKey)

	es := new(EnrichmentService)
	es.db = db
	es.locWait = waitBetweenEnrichment
	es.yelpWait = waitBetweenEnrichment
	es.sanitizeRegex = strings.NewReplacer("-", "", ",", "")

	auth := &yelp.AuthOptions{
		ConsumerKey:       cfg.Yelp.ConsumerKey,
		ConsumerSecret:    cfg.Yelp.ConsumerSecret,
		AccessToken:       cfg.Yelp.AccessToken,
		AccessTokenSecret: cfg.Yelp.AccessTokenSecret,
	}
	es.yelpClient = yelp.New(auth, nil)

	return es
}

func (es *EnrichmentService) EnrichYelp() {
	for {
		var locs []Location
		es.db.Limit(enrichmentLimit).Where("yelp_type = ? and yelp_url = ?", "", "").Find(&locs)

		for _, loc := range locs {
			info := es.getYelpInfo(loc)
		}

		es.throttleWait(len(locs))
		time.Sleep(es.wait)
	}
}

// Enrich Periodically check the DB and enrich records for city + country info
func (es *EnrichmentService) EnrichLocation() {
	for {
		var locs []Location
		es.db.Limit(enrichmentLimit).Where("city = ? and country = ?", "", "").Find(&locs)

		for _, loc := range locs {
			geo := reverseGeocode(&loc)
			es.updateLocGeo(geo, &loc)
		}

		es.throttleWait(len(locs))
		time.Sleep(es.wait)
	}
}

func (es *EnrichmentService) updateLocGeo(geo *geocoder.Location, loc *Location) {
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

	es.db.Save(loc)
}

func (es *EnrichmentService) throttleWait(found int) {
	if found < enrichmentLimit {
		es.wait += waitBetweenEnrichment
	} else {
		es.wait = waitBetweenEnrichment
	}
}

func (es *EnrichmentService) getYelpInfo(loc *Location) *yelp.Business {
	opts := yelp.SearchOptions{
		GeneralOptions: &yelp.GeneralOptions{
			Term: es.sanitize(loc.Name),
		},
		LocationOptions: &yelp.LocationOptions{
			CoordinateOptions: &yelp.CoordinateOptions{
				Latitude:  null.FloatFrom(loc.Lat),
				Longitude: null.FloatFrom(loc.Long),
			},
		},
	}

	if rs, err := es.yelpClient.DoSearch(opts); err != nil {
		return es.filterYelpResults(loc, &rs)
	}

	return nil
}

func (es *EnrichmentService) filterYelpResults(loc *Location, rs *yelp.SearchResult) *yelp.Business {
	for _, r := range rs.Businesses {
		if es.sanitize(loc.Name) == es.sanitize(r.Name) {
			return &r
		}
	}

	return nil
}
func (es *EnrichmentService) sanitize(s string) string {
	return es.sanitizeRegex.Replace(strings.ToLower(s))
}

func reverseGeocode(loc *Location) *geocoder.Location {
	if geo, err := geocoder.ReverseGeocode(loc.Lat, loc.Long); err == nil {
		return geo
	} else {
		logger.Error(fmt.Sprintf("Reverse geocoding encountered and error: %s", err.Error()))
	}

	return nil
}
