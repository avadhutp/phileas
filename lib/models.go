package lib

import (
	"time"
)

// Entries Contains all of the instagrams likes indexed by required fields
type Entries struct {
	ID        int       `sql:"AUTO_INCREMENT"`
	Type      string    `sql:"NOT NULL"`
	VendorID  string    `sql:"NOT NULL"`
	Timestamp time.Time `sql:"NOT NULL"`
}
