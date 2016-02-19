package lib

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/net/context"

	"googlemaps.github.io/maps"

	"github.com/jasonwinn/geocoder"
	"github.com/jinzhu/gorm"
)

const (
	typeLoc = iota
	typeGooglePlaces
)

var (
	esEnrichLocation func(*EnrichmentService)

	geocoderSetAPIKey      = geocoder.SetAPIKey
	geocoderReverseGeocode = geocoder.ReverseGeocode
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
	db               *gorm.DB
	sanitizeRegex    *strings.Replacer
	waits            map[int]time.Duration
	googleMapsClient *maps.Client
}

// NewEnrichmentService Provider for EnrichmentService
func NewEnrichmentService(cfg *Cfg, db *gorm.DB) *EnrichmentService {
	geocoderSetAPIKey(cfg.Common.MapquestKey)

	es := new(EnrichmentService)
	es.db = db
	es.sanitizeRegex = strings.NewReplacer("-", "", ",", "")
	es.waits = map[int]time.Duration{
		typeLoc:          waitBetweenEnrichment,
		typeGooglePlaces: waitBetweenEnrichment,
	}

	es.googleMapsClient, _ = maps.NewClient(maps.WithAPIKey("AIzaSyBKLh4PJpmZF5YE4tQwul8yfld_Z-Qu_Gw"))

	return es
}

// EnrichLocation Periodically check the DB and enrich records for city + country info
func (es *EnrichmentService) EnrichLocation() {
	var locs []Location
	es.db.Limit(enrichmentLimit).Where("city = ? and country = ?", "", "").Find(&locs)

	logger.Infof("Enriching %d locations", len(locs))

	for _, loc := range locs {
		geo := reverseGeocode(&loc)
		es.updateLocGeo(geo, &loc)
	}

	es.throttleWait(len(locs), typeLoc)
	timeSleep(es.waits[typeLoc])
	esEnrichLocation(es)
}

func (es *EnrichmentService) EnrichGooglePlacesIDs() {
	var locs []Location
	es.db.Limit(enrichmentLimit).Where("google_places_id IS NULL").Find(&locs)

	logger.Infof("Enriching %d locations for google places IDs", len(locs))

	for _, loc := range locs {
		req := newRadarSearch(&loc)
		_, err := es.googleMapsClient.RadarSearch(context.Background(), req)

		fmt.Println(err.Error())
	}

	es.throttleWait(len(locs), typeGooglePlaces)
	timeSleep(es.waits[typeGooglePlaces])
	es.EnrichGooglePlacesIDs()
}

func (es *EnrichmentService) updateLocGeo(geo *geocoder.Location, loc *Location) {
	if geo == nil || geo.CountryCode == "" {
		return
	}

	loc.Country = geo.CountryCode
	loc.Address = makeAddress(geo)
	copyGeoToLoc(loc, geo)

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
	if geo, err := geocoderReverseGeocode(loc.Lat, loc.Long); err != nil {
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

func copyGeoToLoc(loc *Location, geo *geocoder.Location) {
	if geo.City != "" {
		loc.City = geo.City
	} else if geo.County != "" {
		loc.City = geo.County
	} else if geo.State != "" {
		loc.City = geo.State
	}
}

func newRadarSearch(l *Location) *maps.RadarSearchRequest {
	r := &maps.RadarSearchRequest{}

	r.Name = l.Name
	r.Radius = 500
	r.Location = &maps.LatLng{
		Lat: l.Lat,
		Lng: l.Long,
	}

	return r
}
