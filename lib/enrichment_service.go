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
	typeYelp = iota
	typeLoc
)

const (
	enrichmentLimit       = 10
	waitBetweenEnrichment = 30 * time.Second
	exemptErrorID         = "UNAVAILABLE_FOR_LOCATION"
)

// EnrichmentService Goes systematically and enriches existing Location records with city + country information
type EnrichmentService struct {
	db            *gorm.DB
	yelpClient    *yelp.Client
	sanitizeRegex *strings.Replacer
	waits         map[int]time.Duration
}

// NewEnrichmentService Provider for EnrichmentService
func NewEnrichmentService(cfg *Cfg, db *gorm.DB) *EnrichmentService {
	geocoder.SetAPIKey(cfg.Common.MapquestKey)

	es := new(EnrichmentService)
	es.db = db
	es.sanitizeRegex = strings.NewReplacer("-", "", ",", "")
	es.waits = map[int]time.Duration{
		typeYelp: waitBetweenEnrichment,
		typeLoc:  waitBetweenEnrichment,
	}

	auth := &yelp.AuthOptions{
		ConsumerKey:       cfg.Yelp.ConsumerKey,
		ConsumerSecret:    cfg.Yelp.ConsumerSecret,
		AccessToken:       cfg.Yelp.AccessToken,
		AccessTokenSecret: cfg.Yelp.AccessTokenSecret,
	}
	es.yelpClient = yelp.New(auth, nil)

	return es
}

//EnrichYelp Periodically check the DB and enrich records for yelp category and URL
func (es *EnrichmentService) EnrichYelp() {
	for {
		var locs []Location
		es.db.Limit(enrichmentLimit).Where("yelp_type = ? and yelp_url = ?", "", "").Find(&locs)

		for _, loc := range locs {
			info := es.getYelpInfo(&loc)
			es.updateLocYelp(info, &loc)
		}

		es.throttleWait(len(locs), typeYelp)
		time.Sleep(es.waits[typeYelp])
	}
}

// EnrichLocation Periodically check the DB and enrich records for city + country info
func (es *EnrichmentService) EnrichLocation() {
	for {
		var locs []Location
		es.db.Limit(enrichmentLimit).Where("city = ? and country = ?", "", "").Find(&locs)

		for _, loc := range locs {
			geo := reverseGeocode(&loc)
			es.updateLocGeo(geo, &loc)
		}

		es.throttleWait(len(locs), typeLoc)
		time.Sleep(es.waits[typeLoc])
	}
}

func (es *EnrichmentService) updateLocYelp(info *yelp.Business, loc *Location) {
	if info == nil {
		loc.YelpType = "NA"
		loc.YelpURL = "NA"
	} else {
		if len(info.Categories) > 0 {
			loc.YelpType = info.Categories[0][0]
		}
		loc.YelpURL = info.URL
	}

	es.db.Save(loc)
}

func (es *EnrichmentService) updateLocGeo(geo *geocoder.Location, loc *Location) {
	if geo == nil || geo.CountryCode == "" {
		return
	}

	loc.Country = geo.CountryCode
	loc.Address = makeAddress(geo)

	if geo.City != "" {
		loc.City = geo.City
	} else if geo.County != "" {
		loc.City = geo.County
	} else if geo.State != "" {
		loc.City = geo.State
	}

	es.db.Save(loc)
}

func (es *EnrichmentService) throttleWait(found int, w int) {
	if found < enrichmentLimit {
		es.waits[w] += waitBetweenEnrichment
	} else {
		es.waits[w] = waitBetweenEnrichment
	}
}

func (es *EnrichmentService) getYelpInfo(loc *Location) *yelp.Business {
	opts := yelp.SearchOptions{
		GeneralOptions: &yelp.GeneralOptions{
			Term: es.sanitize(loc.Name),
		},
		CoordinateOptions: &yelp.CoordinateOptions{
			Latitude:  null.FloatFrom(loc.Lat),
			Longitude: null.FloatFrom(loc.Long),
		},
	}

	if rs, err := es.yelpClient.DoSearch(opts); err != nil {
		if !exemptYelpError(err) {
			logger.Error(fmt.Sprintf("Error fetching yelp info: %s", err.Error()))
		}
	} else {
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
	if geo, err := geocoder.ReverseGeocode(loc.Lat, loc.Long); err != nil {
		logger.Error(fmt.Sprintf("Reverse geocoding encountered and error: %s", err.Error()))
	} else {
		return geo
	}

	return nil
}

func exemptYelpError(err error) bool {
	return strings.Contains(err.Error(), exemptErrorID)
}

func notEmpty(s string) string {
	if s == "" {
		return "NA"
	}

	return s
}

func makeAddress(geo *geocoder.Location) string {
	els := []string{
		geo.Street,
		geo.City,
		geo.State,
		geo.County,
		geo.PostalCode,
	}

	var addr []string

	for _, el := range els {
		if el != "" {
			addr = append(addr, el)
		}
	}

	return strings.Join(addr, ", ")
}
