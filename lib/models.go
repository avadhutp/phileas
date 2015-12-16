package lib

import (
	"fmt"
	"time"
)

// Entry Contains all of the instagrams likes indexed by required fields
type Entry struct {
	ID        int       `sql:"AUTO_INCREMENT"`
	Type      string    `sql:"NOT NULL"`
	VendorID  string    `sql:"NOT NULL"`
	Timestamp time.Time `sql:"NOT NULL"`
}

// GetDBConnString Provides the MySQL DSN
func GetDBConnString(cfg *Cfg) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.Mysql.Username,
		cfg.Mysql.Password,
		cfg.Mysql.Host,
		cfg.Mysql.Port,
		cfg.Mysql.Database)
}
