package lib

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jasonwinn/geocoder"
	"github.com/jinzhu/gorm"
	"github.com/kpawlik/geojson"
)

const (
	countryCacheFilePath = "static/country_latlon.csv"
	imgTmpl              = `
		<div class="row">
			<div class="left"><a href="%s" target="_blank"><img src="%s" /></a></div>
			<div class="right">%s</div>
		</div>`
)

// CountryInfo struct to store country obj
type CountryInfo struct {
	name string
	geo  geocoder.LatLng
}

// LocEntry Results struct for MySQL join queries
type LocEntry struct {
	VendorID   string
	LocationID int
	Name       string
	Lat        float64
	Long       float64
}

// PhileasAPI Provides the data for phileas's API
type PhileasAPI struct {
	googleBrowserKey string
	instaAPI         *InstaAPI
	db               *gorm.DB
	countryCache     map[string]CountryInfo
}

// NewPhileasAPI Go-style constructor to provide an instance of Phileas's API
func NewPhileasAPI(cfg *Cfg, db *gorm.DB, instaAPI *InstaAPI) *PhileasAPI {
	api := &PhileasAPI{}
	api.googleBrowserKey = cfg.Google.BrowserKey

	api.db = db
	api.instaAPI = instaAPI

	api.countryCache = cacheCountryLatLong(countryCacheFilePath)

	return api
}

// ping — /ping
func (pe *PhileasAPI) ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

// location - /loc/:location-id
func (pe *PhileasAPI) location(c *gin.Context) {
	l, _ := strconv.Atoi(c.Param("location-id"))
	var entries []*Entry
	pe.db.Where(&Entry{LocationID: l}).Find(&entries)

	var out []string
	for _, e := range entries {
		out = append(out, fmt.Sprintf(imgTmpl, e.URL, e.Thumbnail, e.Caption))
	}

	c.String(http.StatusOK, strings.Join(out, "<br />"))
}

// mapper -/top
func (pe *PhileasAPI) mapper(c *gin.Context) {
	c.HTML(http.StatusOK, "mapper.tmpl", gin.H{
		"title": "Top destinations",
		"key":   pe.googleBrowserKey,
	})
}

// countriesJSON - /countries.json
func (pe *PhileasAPI) countriesJSON(c *gin.Context) {
	var total int
	pe.db.Model(&Location{}).Where("country != ?", "").Count(&total)

	rows, err := pe.db.Table("location").Select("`id`, `country`, count(*)").Group("country").Having("`country` != ''").Rows()

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		col := makeGroupedGeoJSON(rows, pe.countryCache, total)
		c.JSON(http.StatusOK, col)
	}
}

// topJSON — /top.json
func (pe *PhileasAPI) topJSON(c *gin.Context) {
	var locs []*Location
	pe.db.Find(&locs)
	col := makeGeoJSON(locs)

	c.JSON(http.StatusOK, col)
}

// statsJSON - /stats.json
func (pe *PhileasAPI) statsJSON(c *gin.Context) {
	var totalLocations int
	var gpEnriched int
	var gpUnenriched int
	var gpWithoutID int

	pe.db.Model(&Location{}).Count(&totalLocations)

	pe.db.Model(&Location{}).Where("google_places_id IS NOT NULL and google_places_id != ?", "").Count(&gpEnriched)
	pe.db.Model(&Location{}).Where("google_places_id IS NULL").Count(&gpUnenriched)
	pe.db.Model(&Location{}).Where("google_places_id = ?", "").Count(&gpWithoutID)

	fmt.Printf("All nice %d", totalLocations)

	stats := map[string]interface{}{
		"total_locations": totalLocations,
		"enrichment": map[string]interface{}{
			"google_places": map[string]interface{}{
				"enriched":         gpEnriched,
				"un_enriched":      gpUnenriched,
				"without_place_id": gpWithoutID,
			},
		},
	}

	c.JSON(http.StatusOK, stats)
}

func makeGroupedGeoJSON(rows *sql.Rows, cache map[string]CountryInfo, total int) *geojson.FeatureCollection {
	var all []*geojson.Feature

	for rows.Next() {
		var ID, count int
		var country string

		rows.Scan(&ID, &country, &count)
		p := geojson.NewPoint(geojson.Coordinate{geojson.CoordType(cache[country].geo.Lng), geojson.CoordType(cache[country].geo.Lat)})

		props := map[string]interface{}{
			"id":      ID,
			"size":    (float64(count*100) / float64(total)),
			"country": cache[country].name,
			"total":   count,
		}

		f := geojson.NewFeature(p, props, nil)
		all = append(all, f)
	}

	col := geojson.NewFeatureCollection(all)
	return col
}

func makeGeoJSON(locs []*Location) *geojson.FeatureCollection {
	var all []*geojson.Feature

	for _, loc := range locs {
		p := geojson.NewPoint(geojson.Coordinate{geojson.CoordType(loc.Long), geojson.CoordType(loc.Lat)})

		props := map[string]interface{}{
			"id":              loc.ID,
			"name":            loc.Name,
			"google_place_id": loc.GooglePlacesID,
		}

		f := geojson.NewFeature(p, props, nil)
		all = append(all, f)
	}

	col := geojson.NewFeatureCollection(all)
	return col
}

func cacheCountryLatLong(path string) map[string]CountryInfo {
	file, _ := os.Open(path)
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	rows, _ := reader.ReadAll()
	cache := make(map[string]CountryInfo, len(rows))

	for i := range rows {
		lat, _ := strconv.ParseFloat(rows[i][1], 64)
		lng, _ := strconv.ParseFloat(rows[i][2], 64)
		cache[rows[i][0]] = CountryInfo{
			name: rows[i][3],
			geo: geocoder.LatLng{
				Lat: lat,
				Lng: lng,
			},
		}
	}

	return cache
}
