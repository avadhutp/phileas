package lib

import (
	"github.com/jinzhu/gorm"

	"fmt"
)

// Entry Contains all of the instagrams likes indexed by required fields
type Entry struct {
	ID         int    `sql:"AUTO_INCREMENT"`
	Type       string `sql:"NOT NULL"`
	VendorID   string `sql:"NOT NULL"`
	Timestamp  int64  `sql:"NOT NULL"`
	LocationID int    `sql:"NOT NULL"`
}

type Location struct {
	ID       int     `sql:"AUTO_INCREMENT"`
	Name     string  `sql:"NOT NULL;type: varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci"`
	Lat      float64 `sql:"NOT NULL"`
	Long     float64 `sql:"NOT NULL"`
	Address  string  `sql:"type: text CHARACTER SET utf8 COLLATE utf8_general_ci"`
	Country  string  `sql:"type: varchar(2) CHARACTER SET utf8 COLLATE utf8_general_ci"`
	City     string  `sql:"type: text CHARACTER SET utf8 COLLATE utf8_general_ci"`
	YelpType string  `sql:"type: text CHARACTER SET utf8 COLLATE utf8_general_ci"`
	YelpURL  string  `sql:"type: text CHARACTER SET utf8 COLLATE utf8_general_ci"`
}

func getDBConnString(cfg *Cfg) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?character_set_server=utf8mb4&parseTime=True&loc=Local",
		cfg.Mysql.Username,
		cfg.Mysql.Password,
		cfg.Mysql.Host,
		cfg.Mysql.Port,
		cfg.Mysql.Database)
}

// GetDB Provides the DB connection with all the required options set
func GetDB(cfg *Cfg) *gorm.DB {
	db, err := gorm.Open("mysql", getDBConnString(cfg))

	if err != nil {
		panic(err)
	}

	db.SingularTable(true)

	return &db
}
