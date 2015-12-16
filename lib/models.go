package lib

import (
	"github.com/jinzhu/gorm"

	"fmt"
)

// Entry Contains all of the instagrams likes indexed by required fields
type Entry struct {
	ID        int    `sql:"AUTO_INCREMENT"`
	Type      string `sql:"NOT NULL"`
	VendorID  string `sql:"NOT NULL"`
	Location  string `sql:"NOT NULL"`
	Timestamp int64  `sql:"NOT NULL"`
}

func getDBConnString(cfg *Cfg) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
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