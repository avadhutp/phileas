package lib

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/kpawlik/geojson"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

const (
	imgTmpl = `
		<div class="row">
			<div class="left"><a href="%s" target="_blank"><img src="%s" /></a></div>
			<div class="right">%s</div>
		</div>`
)

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
	googleKey string
	instaAPI  *InstaAPI
	db        *gorm.DB
}

// NewPhileasAPI Go-style constructor to provide an instance of Phileas's API
func NewPhileasAPI(cfg *Cfg, db *gorm.DB, instaAPI *InstaAPI) *PhileasAPI {
	api := &PhileasAPI{}
	api.googleKey = cfg.Common.GoogleMapsKey

	api.db = db
	api.instaAPI = instaAPI

	return api
}

// ping — /ping
func (pe *PhileasAPI) ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

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
		"key":   pe.googleKey,
	})
}

// countriesJSON - /countries.json
func (pe *PhileasAPI) countriesJSON(c *gin.Context) {
	rows, err := pe.db.Table("location").Select("`id`, `country`, `lat`, `long`, count(*)").Group("country").Having("`country` != ''").Rows()

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		col := makeGroupedGeoJSON(rows)
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

func makeGroupedGeoJSON(rows *sql.Rows) *geojson.FeatureCollection {
	var all []*geojson.Feature

	for rows.Next() {
		var ID, count int
		var country string
		var lat, long float64

		rows.Scan(&ID, &country, &lat, &long, &count)

		p := geojson.NewPoint(geojson.Coordinate{geojson.CoordType(long), geojson.CoordType(lat)})

		props := map[string]interface{}{
			"id":      ID,
			"count":   count,
			"country": country,
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
			"id":   loc.ID,
			"name": loc.Name,
		}

		f := geojson.NewFeature(p, props, nil)
		all = append(all, f)
	}

	col := geojson.NewFeatureCollection(all)
	return col
}
