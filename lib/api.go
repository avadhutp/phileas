package lib

import (
	"net/http"

	"github.com/kpawlik/geojson"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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

// mapper -/top
func (pe *PhileasAPI) mapper(c *gin.Context) {
	c.HTML(http.StatusOK, "mapper.tmpl", gin.H{
		"title": "Top destinations",
		"key":   pe.googleKey,
	})
}

// topJSON — /top.json
func (pe *PhileasAPI) topJSON(c *gin.Context) {
	var results []*LocEntry
	pe.db.Table("location").Select("entry.vendor_id, entry.location_id, location.name, location.lat, location.long").Joins("join entry on location.id = entry.location_id").Scan(&results)
	// col := makeGeoJSON(locs)

	c.JSON(http.StatusOK, results)
}

func makeGeoJSON(locs []*Location) *geojson.FeatureCollection {
	var all []*geojson.Feature

	for _, loc := range locs {
		p := geojson.NewPoint(geojson.Coordinate{geojson.CoordType(loc.Long), geojson.CoordType(loc.Lat)})
		props := map[string]interface{}{
			"content": loc.Name,
		}

		f := geojson.NewFeature(p, props, nil)
		all = append(all, f)
	}

	col := geojson.NewFeatureCollection(all)
	return col
}
