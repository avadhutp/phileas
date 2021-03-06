package lib

import (
	"github.com/jinzhu/gorm"

	"fmt"
)

var (
	gormOpen        = gorm.Open
	dbSingularTable = (*gorm.DB).SingularTable

	locationCols = []string{"id", "name", "lat", "long", "address", "country", "city", "google_places_id"}
	entryCols    = []string{"id", "type", "vendorid", "thumbnail", "url", "caption", "timestamp", "loctionid"}
)

// Entry Contains all of the instagrams likes indexed by required fields
type Entry struct {
	ID         int    `sql:"AUTO_INCREMENT"`
	Type       string `sql:"NOT NULL"`
	VendorID   string `sql:"NOT NULL"`
	Thumbnail  string `sql:"type: varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci"`
	URL        string `sql:"type: varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci"`
	Caption    string `sql:"type: text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci"`
	Timestamp  int64  `sql:"NOT NULL"`
	LocationID int    `sql:"NOT NULL"`
}

// Location Struct to hold all the location information
type Location struct {
	ID             int     `sql:"AUTO_INCREMENT"`
	Name           string  `sql:"NOT NULL;type: varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci"`
	Lat            float64 `sql:"NOT NULL"`
	Long           float64 `sql:"NOT NULL"`
	Address        string  `sql:"type: text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci"`
	Country        string  `sql:"type: varchar(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci"`
	City           string  `sql:"type: text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci"`
	GooglePlacesID string  `sql:"type: varchar(255)"`
}

func getDBConnString(cfg *Cfg) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&character_set_server=utf8mb4&parseTime=True&loc=Local",
		cfg.Mysql.Username,
		cfg.Mysql.Password,
		cfg.Mysql.Host,
		cfg.Mysql.Port,
		cfg.Mysql.Database)
}

// GetDB Provides the DB connection with all the required options set
func GetDB(cfg *Cfg) *gorm.DB {
	db, err := gormOpen("mysql", getDBConnString(cfg))

	if err != nil {
		panic(err)
	}

	dbSingularTable(db, true)

	return db
}
