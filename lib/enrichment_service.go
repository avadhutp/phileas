package lib

import (
	"fmt"
	"strings"
	"time"

	"github.com/jasonwinn/geocoder"
	"github.com/jinzhu/gorm"
)

const (
	typeLoc = iota
)

var (
	esEnrichLocation  func(*EnrichmentService)
	geocoderSetAPIKey = geocoder.SetAPIKey
)

func init() {
	esEnrichLocation = (*EnrichmentService).EnrichLocation
}

const (
	enrichmentLimit       = 10
	waitBetweenEnrichment = 30 * time.Second
	exemptErrorID         = "UNAVAILABLE_FOR_LOCATION"
)

// EnrichmentService Goes systematically and enriches existing Location records with city + country information
type EnrichmentService struct {
	db            *gorm.DB
	sanitizeRegex *strings.Replacer
	waits         map[int]time.Duration
}

// NewEnrichmentService Provider for EnrichmentService
func NewEnrichmentService(cfg *Cfg, db *gorm.DB) *EnrichmentService {
	geocoderSetAPIKey(cfg.Common.MapquestKey)

	es := new(EnrichmentService)
	es.db = db
	es.sanitizeRegex = strings.NewReplacer("-", "", ",", "")
	es.waits = map[int]time.Duration{
		typeLoc: waitBetweenEnrichment,
	}

	return es
}

// EnrichLocation Periodically check the DB and enrich records for city + country info
func (es *EnrichmentService) EnrichLocation() {
	var locs []Location
	es.db.Limit(enrichmentLimit).Where("city = ? and country = ?", "", "").Find(&locs)

	for _, loc := range locs {
		geo := reverseGeocode(&loc)
		es.updateLocGeo(geo, &loc)
	}

	es.throttleWait(len(locs), typeLoc)
	timeSleep(es.waits[typeLoc])
	esEnrichLocation(es)
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

func reverseGeocode(loc *Location) *geocoder.Location {
	if geo, err := geocoder.ReverseGeocode(loc.Lat, loc.Long); err != nil {
		logger.Error(fmt.Sprintf("Reverse geocoding encountered and error: %s", err.Error()))
	} else {
		return geo
	}

	return nil
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
